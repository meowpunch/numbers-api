package app

import (
	"net/http"
	"numbers-api/cache"
	port "numbers-api/http"
	"numbers-api/usecase"
	"time"
)

type App struct {
	Server *http.Server
}

func NewApp() *App {
	purgeDuration := 10 * time.Minute
	inMemory := cache.NewInMemoryCache(purgeDuration)
	timeout := 450 * time.Millisecond
	client := http.DefaultClient
	getNumbers := usecase.NewGetNumbersFunc(*client, timeout, inMemory)

	handler := port.NewRouter(getNumbers)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return &App{
		Server: server,
	}
}

func (a *App) Run() error {
	return a.Server.ListenAndServe()
}
