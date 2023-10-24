package usecase

import (
	"net/http"
	"numbers-api/cache"
	"numbers-api/usecase/internal/fetcher"
	"numbers-api/usecase/internal/merger"
	"time"
)

type GetNumbersFunc func(urls []string) []int

func NewGetNumbersFunc(client http.Client, timeout time.Duration, cache cache.Cache) GetNumbersFunc {
	return func(urls []string) []int {
		finalNumbers := []int{}
		ch := make(chan []int, len(urls))
		timeoutCh := time.After(timeout)

		for _, url := range urls {
			go func(u string) {
				nums := fetcher.FetchNumbers(client, u, timeout, cache) // We're ignoring errors here
				ch <- nums
			}(url)
		}

		// Count for number of responses we've processed
		count := 0

		// Main loop for processing responses
		for {
			select {
			case nums := <-ch:
				finalNumbers = merger.MergeAndSortUnique(finalNumbers, nums)
				count++
				if count == len(urls) { // all urls processed
					return finalNumbers
				}
			case <-timeoutCh:
				// Timeout occurred, return whatever numbers we have so far
				return finalNumbers
			}
		}
	}
}
