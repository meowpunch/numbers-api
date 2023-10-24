package http

import (
	"net/http"
	"numbers-api/usecase"
)

func NewRouter(getNumbers usecase.GetNumbersFunc) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/numbers", NewHandlerGetNumbers(getNumbers))
	return mux
}
