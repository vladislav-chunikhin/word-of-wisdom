package main

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"wordofwisdom/internal/app"
	"wordofwisdom/pkg/config"
)

const (
	handlerQuote byte = 0x01 // handler ID for the quote handler
)

func main() {
	// parse the config
	cfg, err := config.ClientParse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// set the log level
	app.SetLogger(cfg.LogLevel)

	// create a ticker to control the request rate
	ticker := time.NewTicker(time.Second / time.Duration(cfg.RPS))
	defer ticker.Stop()

	// restrict the number of concurrent requests
	var wg sync.WaitGroup
	wg.Add(cfg.TotalRequests)
	for i := 0; i < cfg.TotalRequests; i++ {
		<-ticker.C

		go func(requestNum int) {
			defer wg.Done()
			err = sendRequest(cfg.ServerAddr, handlerQuote)
			if err != nil {
				slog.Error("request failed", "request_num", requestNum, "error", err)
			} else {
				slog.Info("request succeeded", "request_num", requestNum)
			}
		}(i)
	}

	// wait for all requests to finish
	wg.Wait()
}
