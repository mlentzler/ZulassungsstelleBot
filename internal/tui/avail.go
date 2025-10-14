package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

var weekdays = []string{"MO", "TU", "WE", "TH", "FR"}

func initAvailInputs(m *Model) {
	m.dateInput = textinput.New()
	m.dateInput.Placeholder = "DD.MM.YYYY"
	m.dateInput.CharLimit = 10
	m.dateInput.Width = 12

	m.fromInput = textinput.New()
	m.fromInput.Placeholder = "From 0-23Uhr"
	m.fromInput.CharLimit = 2
	m.fromInput.Width = 5

	m.toInput = textinput.New()
	m.toInput.Placeholder = "To 1-24Uhr"
	m.toInput.CharLimit = 2
	m.toInput.Width = 5
}

func validateDateISO(s string) (time.Time, error) {
	t, err := time.Parse("02.01.2006", strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("Datum muss DD.MM.YYYY sein")
	}
	return t, nil
}

func parseHour(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Stunde muss eine Zahl sein")
	}
	return n, nil
}

func validateHours(from, to int) error {
	if from < 0 || from > 23 {
		return fmt.Errorf("From außerhalb von 0-23Uhr")
	}
	if to < 1 || to > 24 {
		return fmt.Errorf("To außerhalb von 1-24Uhr")
	}
	if to <= from {
		return fmt.Errorf("To muss > From sein")
	}
	return nil
}

// updateAvailMode(m, msg) & updateAvailDetail(m, msg) implementieren
// (Keybindings: Pfeile/Tab/Enter, Fehler → m.errMsg)
