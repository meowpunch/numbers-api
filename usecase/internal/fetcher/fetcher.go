package fetcher

import (
	"encoding/json"
	"net/http"
	"numbers-api/cache"
	"regexp"
	"strconv"
	"time"
)

func FetchNumbers(client http.Client, url string, timeout time.Duration, cache cache.Cache) []int {
	// If data exists in cache, return it
	if data, found := cache.Get(url); found {
		return data
	}

	resultCh := make(chan []int, 1)
	errorCh := make(chan error, 1)

	go func() {
		var response struct {
			Numbers []int `json:"numbers"`
		}

		resp, err := client.Get(url)
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
					cache.Set(url, response.Numbers, duration)
				}
			}
		}

		resultCh <- response.Numbers
	}()

	// Wait for the fetch to complete or timeout
	select {
	case nums := <-resultCh:
		return nums
	case <-errorCh:
		return nil
	case <-time.After(timeout):
		// Return nil synchronously but the asynchronous fetch will continue and store in cache if possible
		return nil
	}
}
