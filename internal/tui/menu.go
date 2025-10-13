package tui

import (
	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

func currentLevelNodes(m *Model) []domain.MenuNode {
	// TODO: aus m.menuRoot & m.menuStack die aktuelle Ebene ermitteln
	return nil
}

func pushMenu(m *Model, idx int) {
	// TODO: ausgewählten Title in m.path pushen, neue Ebene öffnen
}

func popMenu(m *Model) {
	// TODO: eine Ebene zurück (stack & path anpassen)
}

// updateMenu(m, msg) & viewMenu(m) bitte in model.go / view.go aufrufen
