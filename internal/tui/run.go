package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mlentzler/ZulassungsstelleBot/internal/config"
	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

func Run(ctx context.Context, cfg config.Config) (domain.BookingRequest, error) {
	root, err := config.LoadMenu(cfg.MenuPath)
	if err != nil {
		return domain.BookingRequest{}, fmt.Errorf("lade menu: %w", err)
	}

	m := NewModel(root, cfg)
	p := tea.NewProgram(m)

	// Kontext-Abbruch behandeln
	go func() {
		<-ctx.Done()
		p.Quit()
	}()

	res, err := p.Run()
	if err != nil {
		return domain.BookingRequest{}, err
	}
	final, ok := res.(Model)
	if !ok {
		return domain.BookingRequest{}, fmt.Errorf("unexpected model type")
	}
	if final.step != stepDone || final.result == nil {
		return domain.BookingRequest{}, context.Canceled
	}
	return *final.result, nil
}
