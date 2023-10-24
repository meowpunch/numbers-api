package fetcher

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"numbers-api/cache"
	"testing"
	"time"
)

func TestFetchNumbers(t *testing.T) {
	inMemory := cache.NewInMemoryCache(1 * time.Second)
	client := &http.Client{}
	timeout := 10 * time.Millisecond

	t.Run("no pre-cached", func(t *testing.T) {
		t.Run("valid url", func(t *testing.T) {
			tests := []struct {
				name                  string
				setupMock             func(w http.ResponseWriter)
				expectedNumbers       []int
				expectedCachedNumbers []int
			}{
				{
					name: "fetch numbers without cache control",
					setupMock: func(w http.ResponseWriter) {
						w.Write([]byte(`{"numbers": [1,2,3]}`))
					},
					expectedNumbers:       []int{1, 2, 3},
					expectedCachedNumbers: nil,
				},
				{
					name: "fetch and cache numbers",
					setupMock: func(w http.ResponseWriter) {
						w.Header().Set("Cache-Control", "max-age=2")
						w.Write([]byte(`{"numbers": [1,2,3]}`))
					},
					expectedNumbers:       []int{1, 2, 3},
					expectedCachedNumbers: []int{1, 2, 3},
				},
				{
					name: "timeout but cache numbers in background",
					setupMock: func(w http.ResponseWriter) {
						time.Sleep(timeout)
						w.Header().Set("Cache-Control", "max-age=2")
						w.Write([]byte(`{"numbers": [1,2,3]}`))
					},
					expectedNumbers:       nil,
					expectedCachedNumbers: []int{1, 2, 3},
				},
				{
					name: "simply ignore incorrect format",
					setupMock: func(w http.ResponseWriter) {
						w.Write([]byte(`{"invalid": "format"}`))
					},
					expectedNumbers:       nil,
					expectedCachedNumbers: nil,
				},
				{
					name: "simply ignore any error from server ",
					setupMock: func(w http.ResponseWriter) {
						w.WriteHeader(http.StatusInternalServerError)
					},
					expectedNumbers:       nil,
					expectedCachedNumbers: nil,
				},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					hitCounter := 0
					mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						hitCounter++
						tc.setupMock(w)
					}))
					defer mockServer.Close()

					// First fetch
					numbers := FetchNumbers(*client, mockServer.URL, timeout, inMemory)
					assert.Equal(t, numbers, tc.expectedNumbers, "Expected fetched numbers")

					time.Sleep(10 * time.Millisecond)
					// Fetch again
					numbers, _ = inMemory.Get(mockServer.URL)
					assert.Equal(t, numbers, tc.expectedCachedNumbers, "Expected cached numbers")
				})
			}
		})

		t.Run("invalid url", func(t *testing.T) {
			tests := []struct {
				name                  string
				cacheData             []int
				expectedNumbers       []int
				expectedCachedNumbers []int
			}{
				{
					name:                  "simply ignore invalid url",
					cacheData:             []int{4, 5, 6},
					expectedNumbers:       nil,
					expectedCachedNumbers: nil,
				},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {

					// First fetch
					numbers := FetchNumbers(*client, "invalid url", timeout, inMemory)
					assert.Equal(t, numbers, tc.expectedNumbers, "Expected fetched numbers")

					time.Sleep(10 * time.Millisecond)
					// Fetch again
					numbers, _ = inMemory.Get("invalid url")
					assert.Equal(t, numbers, tc.expectedCachedNumbers, "Expected cached numbers")
				})
			}
		})
	})

	t.Run("pre-cached", func(t *testing.T) {
		tests := []struct {
			name            string
			cacheData       []int
			setupMock       func(w http.ResponseWriter)
			expectedNumbers []int
		}{
			{
				name:      "fetch from cache without hitting the server",
				cacheData: []int{4, 5, 6},
				setupMock: func(w http.ResponseWriter) {
					t.Fatal("Should not hit the server if data is cached")
				},
				expectedNumbers: []int{4, 5, 6},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				hitCounter := 0
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					hitCounter++
					tc.setupMock(w)
				}))
				defer mockServer.Close()

				// Preload the cache
				inMemory.Set(mockServer.URL, tc.cacheData, 5*time.Second)

				// Fetch
				numbers := FetchNumbers(*client, mockServer.URL, timeout, inMemory)
				assert.Equal(t, tc.expectedNumbers, numbers, "Expected fetched numbers from cache")
			})
		}
	})
}
