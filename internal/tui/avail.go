package tui

import "time"

func validateDateISO(s string) (time.Time, error) {
	// TODO: time.Parse("2006-01-02", s)
	return time.Time{}, nil
}

func validateHours(from, to int) error {
	// TODO: 0<=from<to<=24
	return nil
}

// updateAvailMode(m, msg) & updateAvailDetail(m, msg) implementieren
// (Keybindings: Pfeile/Tab/Enter, Fehler → m.errMsg)
