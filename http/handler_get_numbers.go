package http

import (
	"encoding/json"
	"net/http"
	"numbers-api/usecase"
)

func NewHandlerGetNumbers(getNumbers usecase.GetNumbersFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := r.URL.Query()["u"]

		numbers := getNumbers(urls)

		response := map[string][]int{
			"numbers": numbers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
