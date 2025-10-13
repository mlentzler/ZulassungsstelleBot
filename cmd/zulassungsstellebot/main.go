package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mlentzler/ZulassungsstelleBot/internal/config"
	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
	"github.com/mlentzler/ZulassungsstelleBot/internal/tui"
)

func main() {
	// TODO: optional .env laden
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var req domain.BookingRequest
	var err error

	req, err = tui.Run(ctx, cfg)
	if err != nil {
		if err == context.Canceled {
			log.Println("abgebrochen")
			return
		}
		log.Fatal(err)
	}

	// TODO: hier später Watcher starten / Request an Backend geben
	log.Printf("TUI fertig, BookingRequest: %+v\n", req)
}
