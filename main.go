package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Configuration
	timeout := 500 * time.Millisecond

	// Initialize fetcher, merger, and handler
	client := http.DefaultClient
	cache := NewCache()
	fetcher := NewFetcher(client, cache, timeout)
	handler := NewHandler(fetcher, timeout)

	// Register endpoint
	http.HandleFunc("/numbers", handler.NumbersEndpoint)

	// Start server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
