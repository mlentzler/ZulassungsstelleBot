package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mlentzler/ZulassungsstelleBot/internal/config"
	"github.com/mlentzler/ZulassungsstelleBot/internal/tui"

	drvcdp "github.com/mlentzler/ZulassungsstelleBot/internal/browser/chromedp"
	"github.com/mlentzler/ZulassungsstelleBot/internal/watcher"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	req, err := tui.Run(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if b, e := json.MarshalIndent(req, "", "  "); e == nil {
		log.Printf("TUI fertig, BookingRequest:\n%s\n", string(b))
	} else {
		log.Printf("TUI fertig, BookingRequest: %+v\n", req)
	}

	loc, _ := time.LoadLocation(req.TZ)
	drv, err := drvcdp.NewDriver(cfg.Headless, loc)
	if err != nil {
		log.Fatal(err)
	}

	wcfg := watcher.Config{
		BaseURL:    cfg.BaseURL,
		Headless:   cfg.Headless,
		PollMinSec: cfg.PollMin,
		PollMaxSec: cfg.PollMax,
	}

	if err := watcher.Run(ctx, drv, wcfg, req); err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Termin gebucht!")
}
