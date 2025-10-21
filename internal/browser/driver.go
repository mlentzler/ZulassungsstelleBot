package browser

import (
	"context"
	"time"

	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

type Slot struct {
	Start time.Time
	Ref   any //TODO set
}

type Driver interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error

	StartFlow(ctx context.Context, baseURL string, titles []string, selectors []string) error
	PickDate(ctx context.Context, date time.Time) error
	ListSlots(ctx context.Context) ([]Slot, error)
	BookSlot(ctx context.Context, s Slot, form map[string]string) error
	FillFromMap(ctx context.Context, form map[string]string) error
}

func SlotMatches(av domain.Availability, slot time.Time, loc *time.Location) bool {
	switch av.Kind {
	case domain.AvailOneOff:
		d := slot.In(loc)
		return d.Format("2006-01-02") == av.OneOff.DateISO &&
			d.Hour() >= av.OneOff.FromHour && d.Hour() < av.OneOff.ToHour

	case domain.AvailRecurring:
		d := slot.In(loc)
		goToDe := [...]string{"SO", "MO", "DI", "MI", "DO", "FR", "SA"} // time.Sunday=0
		wd := goToDe[d.Weekday()]
		for _, w := range av.Recurring.Days {
			if w.Weekday == wd && d.Hour() >= w.FromHour && d.Hour() < w.ToHour {
				return true
			}
		}
	}
	return false
}
