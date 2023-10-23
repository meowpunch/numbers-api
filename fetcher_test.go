package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetcher(t *testing.T) {
	cache := NewCache()
	client := &http.Client{}
	fetcher := NewFetcher(client, cache, 100*time.Millisecond)

	tests := []struct {
		name       string
		header     string
		hitCount   int
		expectHits int
	}{
		{
			name:       "fetch number without cache",
			hitCount:   2,
			expectHits: 2,
		},
		{
			name:       "fetch number with cache",
			header:     "max-age=2",
			hitCount:   1,
			expectHits: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hitCounter := 0
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hitCounter++
				if tc.header != "" {
					w.Header().Set("Cache-Control", tc.header)
				}
				w.Write([]byte(`{"numbers": [1,2,3,4,5]}`))
			}))
			defer mockServer.Close()

			// First fetch
			numbers, err := fetcher.FetchNumbersFromURL(mockServer.URL)
			require.NoError(t, err, "Expected successful fetch")
			assert.Len(t, numbers, 5, "Expected 5 numbers")

			// Fetch again
			_, err = fetcher.FetchNumbersFromURL(mockServer.URL)
			require.NoError(t, err, "Expected successful fetch")

			assert.Equal(t, tc.expectHits, hitCounter, "Unexpected number of hits to server")
		})
	}
}

func TestParallelFetching(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2")
		w.Write([]byte("{\"numbers\": [1,2,3,4,5]}")) // Corrected JSON structure here
	}))
	defer mockServer.Close()

	cache := NewCache()
	client := &http.Client{}
	fetcher := NewFetcher(client, cache, 100*time.Millisecond)

	urls := []string{mockServer.URL, mockServer.URL, mockServer.URL}
	ch := fetcher.FetchNumbersInParallel(urls)

	for range urls {
		select {
		case numbers := <-ch:
			if len(numbers) != 5 {
				t.Fatalf("Expected 5 numbers, got %d", len(numbers))
			}
		case <-time.After(2 * time.Second): // Reasonable timeout for test
			t.Fatal("Timed out waiting for numbers")
		}
	}
}
