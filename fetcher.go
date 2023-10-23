package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Fetcher struct {
	client  *http.Client
	cache   *Cache
	timeout time.Duration
}

func NewFetcher(client *http.Client, cache *Cache, timeout time.Duration) *Fetcher {
	return &Fetcher{client: client, cache: cache, timeout: timeout}
}

func (f *Fetcher) FetchNumbersFromURL(url string) ([]int, error) {
	// If data exists in cache, return it
	if data, found := f.cache.Get(url); found {
		return data, nil
	}

	resultCh := make(chan []int, 1)
	errorCh := make(chan error, 1)

	go func() {
		var response struct {
			Numbers []int `json:"numbers"`
		}

		resp, err := f.client.Get(url)
		if err != nil {
			errorCh <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			errorCh <- err
			return
		}

		// Only store in cache if Cache-Control provides max-age
		if cacheControl := resp.Header.Get("Cache-Control"); cacheControl != "" {
			if match := regexp.MustCompile(`max-age=(\d+)`).FindStringSubmatch(cacheControl); len(match) == 2 {
				if seconds, err := strconv.Atoi(match[1]); err == nil {
					duration := time.Duration(seconds) * time.Second
					// Store the result in cache
					f.cache.Set(url, response.Numbers, duration)
				}
			}
		}

		resultCh <- response.Numbers
	}()

	// Wait for the fetch to complete or timeout
	select {
	case nums := <-resultCh:
		return nums, nil
	case err := <-errorCh:
		return nil, err
	case <-time.After(f.timeout):
		// Return nil synchronously but the asynchronous fetch will continue and store in cache if possible
		return nil, nil
	}
}

func (f *Fetcher) FetchNumbersInParallel(urls []string) <-chan []int {
	ch := make(chan []int)

	for _, url := range urls {
		go func(u string) {
			nums, _ := f.FetchNumbersFromURL(u) // We're ignoring errors here
			if nums != nil {
				ch <- nums
			}
		}(url)
	}

	return ch
}
