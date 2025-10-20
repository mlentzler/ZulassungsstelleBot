package chromedpdrv

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Session struct {
	alloc  context.Context
	ctx    context.Context
	cancel context.CancelFunc
}

func New(headless bool) (*Session, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)
	if !headless {
		opts = append(opts,
			chromedp.Flag("auto-open-devtools-for-tabs", true),
		)
	}
	alloc, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	// DEBUG-Logs aus chromedp in dein Terminal
	ctx, cancel := chromedp.NewContext(
		alloc,
		chromedp.WithDebugf(log.Printf),
		// chromedp.WithLogf(log.Printf),
		// chromedp.WithErrorf(log.Printf),
	)
	return &Session{alloc: alloc, ctx: ctx, cancel: cancel}, nil
}

func (s *Session) Context() context.Context { return s.ctx }
func (s *Session) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}
func Sleep(ms int) chromedp.Action { return chromedp.Sleep(time.Duration(ms) * time.Millisecond) }
