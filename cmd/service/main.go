package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flexer2006/t-t-ogen-go/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.NewApplication("")
	if err != nil {
		log.Printf("setup: %v", err)

		return
	}

	err = application.Run(ctx)
	if err != nil {
		log.Printf("application: %v", err)
	}
}
