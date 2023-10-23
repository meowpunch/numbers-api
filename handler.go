package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	fetcher *Fetcher
	timeout time.Duration
}

func NewHandler(fetcher *Fetcher, timeout time.Duration) *Handler {
	return &Handler{
		fetcher: fetcher,
		timeout: timeout,
	}
}

func (h *Handler) NumbersEndpoint(w http.ResponseWriter, r *http.Request) {
	urls := r.URL.Query()["u"]

	// Fetch in parallel
	chans := h.fetcher.FetchNumbersInParallel(urls)

	finalNumbers := []int{}
	for range urls {
		select {
		case nums := <-chans:
			finalNumbers = MergeAndSortUnique(finalNumbers, nums)
		case <-time.After(h.timeout):
			continue
		}
	}

	response := map[string][]int{
		"numbers": finalNumbers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
