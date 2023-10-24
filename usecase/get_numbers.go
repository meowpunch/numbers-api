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

		for _, url := range urls {
			go func(u string) {
				nums := fetcher.FetchNumbers(client, u, timeout, cache)
				ch <- nums
			}(url)
		}

		for range urls {
			nums := <-ch
			finalNumbers = merger.MergeAndSortUnique(finalNumbers, nums)
		}

		close(ch)
		return finalNumbers
	}
}
