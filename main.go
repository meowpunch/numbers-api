package main

import (
	"log"
	"net/http"
	"numbers-api/cache"
	port "numbers-api/http"
	"numbers-api/usecase"
	"time"
)

func main() {
	purgeDuration := 10 * time.Minute
	inMemory := cache.NewInMemoryCache(purgeDuration)
	timeout := 500 * time.Millisecond
	client := http.DefaultClient
	GetNumbers := usecase.NewGetNumbersFunc(*client, timeout, inMemory)
	handler := port.NewHandlerGetNumbers(GetNumbers)

	// Register endpoint
	http.HandleFunc("/numbers", handler)

	// Start server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
