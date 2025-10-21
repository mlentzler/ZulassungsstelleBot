package watcher

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/mlentzler/ZulassungsstelleBot/internal/browser"
	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

type Config struct {
	BaseURL    string
	Headless   bool
	PollMinSec int
	PollMaxSec int
}

func Run(ctx context.Context, drv browser.Driver, cfg Config, req domain.BookingRequest) error {
	loc, _ := time.LoadLocation(req.TZ)
	form := map[string]string{
		"name":    req.Name,
		"email":   req.Email,
		"telefon": req.Phone,
	}

	if err := drv.Open(ctx); err != nil {
		return err
	}
	defer drv.Close(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := drv.StartFlow(ctx, cfg.BaseURL, req.Menu.Path, req.Menu.Selectors); err != nil {
			// WICHTIG: Log mit selector/titel
			// (import "log")
			log.Printf("StartFlow error: %v", err)
			sleep(ctx, jitter(cfg.PollMinSec, cfg.PollMaxSec))
			continue
		}

		if req.Avail.Kind == domain.AvailOneOff && req.Avail.OneOff != nil {
			dt, _ := time.ParseInLocation("2006-01-02", req.Avail.OneOff.DateISO, loc)
			_ = drv.PickDate(ctx, dt)
		}

		slots, err := drv.ListSlots(ctx)
		if err != nil || len(slots) == 0 {
			sleep(ctx, jitter(cfg.PollMinSec, cfg.PollMaxSec))
			continue
		}

		var chosen *browser.Slot
		for i := range slots {
			if browser.SlotMatches(req.Avail, slots[i].Start, loc) {
				chosen = &slots[i]
				break
			}
		}

		if chosen == nil {
			sleep(ctx, jitter(cfg.PollMinSec, cfg.PollMaxSec))
			continue
		}

		if err := drv.BookSlot(ctx, *chosen, form); err != nil {
			sleep(ctx, 3*time.Second)
			continue
		}
		return nil
	}
}

func jitter(minSec, maxSec int) time.Duration {
	if minSec <= 0 {
		minSec = 20
	}
	if maxSec < minSec {
		maxSec = minSec + 25
	}
	sec := rand.Intn(maxSec-minSec+1) + minSec
	return time.Duration(sec) * time.Second
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
