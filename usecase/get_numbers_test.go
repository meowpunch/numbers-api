package usecase

import (
	"net/http"
	"net/http/httptest"
	"numbers-api/cache"
	"testing"
	"time"
)

func TestGetNumbersFunc(t *testing.T) {
	timeout := 10 * time.Millisecond

	tests := []struct {
		name     string
		handlers []http.HandlerFunc
		expected []int
	}{
		{
			name: "All URLs that were successfully retrieved within the given timeframe must influence the result",
			handlers: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"numbers": [3,4,5]}`))
				},
				func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(100 * time.Millisecond) // Delay to simulate timeout
					w.Write([]byte(`{"numbers": [5,6,7]}`))
				},
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name: "return an empty list when all URLs returned errors or took too long to respond and no previous response is stored in the cache",
			handlers: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(timeout)
					w.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				},
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var urls []string
			for _, handler := range tt.handlers {
				server := httptest.NewServer(handler)
				defer server.Close()
				urls = append(urls, server.URL)
			}

			cacheImpl := cache.NewInMemoryCache(1 * time.Minute) // Use the actual in-memory cache
			client := http.Client{}
			getNumbers := NewGetNumbersFunc(client, timeout, cacheImpl)
			result := getNumbers(urls)

			if !equalSlices(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
