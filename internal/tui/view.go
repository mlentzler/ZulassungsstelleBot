package tui

import (
	"fmt"
	"strings"
)

func viewPerson(m Model) string {
	var s string
	s += "📝 Persönliche Daten\n\n"
	s += "Name:\n" + m.nameInput.View() + "\n\n"
	s += "E-Mail:\n" + m.emailInput.View() + "\n\n"
	s += "Telefon:\n" + m.phoneInput.View() + "\n\n"
	if m.errMsg != "" {
		s += "⚠️  " + m.errMsg + "\n\n"
	}
	s += "Enter: weiter · Esc/Ctrl+C: abbrechen\n"
	return s
}

func viewMenu(m Model) string {
	// TODO: breadcrumb + Liste aktuelle Ebene + footer
	var b strings.Builder
	b.WriteString("📂 Auswahl-Menü\n")
	b.WriteString(breadcrumb(&m) + "\n\n")

	nodes := currentLevelNodes(&m)
	if len(nodes) == 0 {
		b.WriteString("(Keine Einträge auf dieser Ebene)\n\n")
		b.WriteString("↑/↓: bewegen · Enter: wählen · Backspace: zurück · Ctrl+C: beenden\n")
		return b.String()
	}

	for i, n := range nodes {
		cursor := "  "
		if i == m.menuCursor {
			cursor = "➤ "
		}
		leaf := " ⤷"
		if len(n.Children) > 0 {
			leaf = " ▸"
		}
		fmt.Fprintf(&b, "%s%s%s\n", cursor, n.Title, leaf)
	}

	b.WriteString("\n↑/↓: bewegen · Enter: wählen · Backspace: zurück · Ctrl+C: beenden\n")
	return b.String()
}

/*
func viewAvailMode(m Model) string {
	// TODO: radio-like Auswahl (Einmalig/Wöchentlich)
	return ""
}

func viewAvailDetail(m Model) string {
	// TODO: je nach mode die Felder anzeigen (dateISO / weekday + from/to)
	return ""
}

func viewReview(m Model) string {
	// TODO: Zusammenfassung + "Enter: bestätigen · Backspace: zurück"
	return ""
}:
*/

func viewAvailMode(m Model) string   { return "Verfügbarkeitsmodus (kommt später)\n" }
func viewAvailDetail(m Model) string { return "Verfügbarkeitsdetails (kommen später)\n" }
func viewReview(m Model) string      { return "Review (kommt später)\n" }
