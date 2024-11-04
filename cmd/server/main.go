package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"wordofwisdom/internal/app"
	"wordofwisdom/internal/handler"
	"wordofwisdom/internal/repository"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create the router and add the quote handler
	router := handler.NewRouter()
	router.AddRoute(handler.HandlerQuote, handler.HandleQuote)

	// create the PoW server
	srv := app.NewPoWServer(ctx, repository.NewQuote(), router)

	// start the server and listen for shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(
		stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("failed to start server", err)
		}
	}()

	slog.Info("server is running...")

	<-stopChan
	slog.Info("shutdown signal received")

	if err := srv.Shutdown(); err != nil {
		slog.Error("error during server shutdown", "error", err)
	} else {
		slog.Info("server shutdown completed")
	}
}
