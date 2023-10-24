//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tsenart/vegeta/v12/lib"
	"log"
	app "numbers-api/server"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	// Start the coding-challenge-test-server container
	reqChallenge := testcontainers.ContainerRequest{
		Image:        "emanuelschmoczer/coding-challenge-test-server:latest",
		ExposedPorts: []string{"8090/tcp"},
		WaitingFor:   wait.ForListeningPort("8090/tcp"),
	}

	challengeServer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: reqChallenge,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer challengeServer.Terminate(ctx)

	challengeServerIP, err := challengeServer.Host(ctx)
	challengeServerPort, err := challengeServer.MappedPort(ctx, "8090")
	challengeServerURL := fmt.Sprintf("http://%s:%s", challengeServerIP, challengeServerPort.Port())

	log.Println("Started challengeServer at:", challengeServerURL)

	// Here, you can start your application's container or directly run main.go, depending on your preference
	// In this example, I'm assuming you run your application via main.go:
	go func() {
		application := app.NewApp()

		log.Println("Server started at :8080")
		log.Fatal(application.Run())
	}()
	time.Sleep(5 * time.Second) // giving some time for your application to start up

	// Run the load test against your API
	loadTest(t, challengeServerURL)
}

func loadTest(t *testing.T, challengeServerURL string) {
	rate := vegeta.Rate{Freq: 100, Per: time.Second} // Adjust this rate as needed
	duration := 10 * time.Second
	endpointURL := fmt.Sprintf(
		"http://127.0.0.1:8080/numbers?u=%s/primes&u=%s/fibo&u=%s/rand&u=%s/rand",
		challengeServerURL, challengeServerURL, challengeServerURL, challengeServerURL,
	)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    endpointURL,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	if metrics.Latencies.P99 > 500*time.Millisecond {
		t.Errorf("99th percentile latency is more than 500ms: %v", metrics.Latencies.P99)
	}
}
