package main

import (
	"log"
	"log/slog"
	"time"

	"lib-go/pkg/client"
)

const (
	address          = ":8080"
	numberOfRequests = 3
)

func main() {
	tcpClient := client.NewClient(address)
	defer tcpClient.Close()
	for i := 0; i < numberOfRequests; i++ {
		err := tcpClient.Connect()
		if err != nil {
			log.Fatalf("can't connect to server: %v", err)
		}

		localWords, localErr := tcpClient.FetchWords()
		if localErr != nil {
			log.Fatalf("Error fetching words from server: %v", localErr)
		}

		slog.Info("Words from server: " + localWords)
		<-time.After(1 * time.Second) // some delay between requests
	}
}
