// main.go
package main

import (
	"log"
	"numbers-api/server"
)

func main() {
	application := app.NewApp()

	log.Println("Server started at :8080")
	log.Fatal(application.Run())
}
