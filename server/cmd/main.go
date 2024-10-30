package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"lib-go/pkg/app"
	"lib-go/pkg/pow"
	"server/internal/service"

	"server/internal/repository"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	server := app.NewServer()
	quoteRep := repository.NewQuote()

	provider, err := pow.NewProvider(server.Config().POW.Complexity)
	if err != nil {
		log.Fatalf("can't init pow provider: %v", err)
	}

	powService := service.NewPOW(quoteRep, provider)
	if err = server.Run(ctx, powService.Handle); err != nil {
		slog.Error(err.Error())
	}
}
