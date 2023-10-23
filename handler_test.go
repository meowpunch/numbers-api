package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_NumbersEndpoint(t *testing.T) {
	// Mock Fetcher & Merger
	cache := NewCache()
	timeout := 10 * time.Millisecond
	fetcher := NewFetcher(http.DefaultClient, cache, timeout)
	merger := NewMerger()

	handler := NewHandler(fetcher, merger, timeout)

	overTimeout := 2 * timeout
	tests := []struct {
		name             string
		mockServerFuncs  []http.HandlerFunc
		expectedHTTPCode int
		expectedNumbers  []int
	}{
		{
			name: "All URLs return within timeout",
			mockServerFuncs: []http.HandlerFunc{
				func(rw http.ResponseWriter, req *http.Request) {
					rw.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(rw http.ResponseWriter, req *http.Request) {
					rw.Write([]byte(`{"numbers": [3,4,5]}`))
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedNumbers:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "First URL has a timeout",
			mockServerFuncs: []http.HandlerFunc{
				func(rw http.ResponseWriter, req *http.Request) {
					rw.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(rw http.ResponseWriter, req *http.Request) {
					time.Sleep(overTimeout)
					rw.Write([]byte(`{"numbers": [3,4,5]}`))
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedNumbers:  []int{1, 2, 3},
		},
		{
			name: "One URL returns a 500 internal server error",
			mockServerFuncs: []http.HandlerFunc{
				func(rw http.ResponseWriter, req *http.Request) {
					rw.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(rw http.ResponseWriter, req *http.Request) {
					rw.WriteHeader(http.StatusInternalServerError)
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedNumbers:  []int{1, 2, 3},
		},
		{
			name: "All URLs have a timeout",
			mockServerFuncs: []http.HandlerFunc{
				func(rw http.ResponseWriter, req *http.Request) {
					time.Sleep(overTimeout)
					rw.Write([]byte(`{"numbers": [1,2,3]}`))
				},
				func(rw http.ResponseWriter, req *http.Request) {
					time.Sleep(overTimeout)
					rw.Write([]byte(`{"numbers": [4,5,6]}`))
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedNumbers:  []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock servers
			var urls []string
			for _, fn := range tc.mockServerFuncs {
				server := httptest.NewServer(fn)
				defer server.Close()
				urls = append(urls, server.URL)
			}

			req, err := http.NewRequest("GET", "/numbers", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			for _, u := range urls {
				q.Add("u", u)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.NumbersEndpoint)
			handlerFunc.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedHTTPCode, rr.Code)

			var response map[string][]int
			err = json.NewDecoder(rr.Body).Decode(&response)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedNumbers, response["numbers"])
		})
	}
}
