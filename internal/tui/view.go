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
	var b strings.Builder
	b.WriteString("📂 Auswahl-Menü\n")
	b.WriteString(breadcrumb(&m) + "\n\n")

	nodes := currentLevelNodes(&m)
	if len(nodes) == 0 {
		b.WriteString("(Keine Einträge auf dieser Ebene)\n\n")
		b.WriteString("↑/↓: bewegen · Enter: wählen · Backspace: zurück · Esc/Ctrl+C: beenden\n")
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

func viewAvailMode(m Model) string {
	o1 := "  Einmaliger Termin"
	o2 := "  Wöchentlich (z. B. Mi 10–13)"
	if m.availCursor == 0 {
		o1 = "➤ Einmaliger Termin"
	} else {
		o2 = "➤ Wöchentlich (z. B. Mi 10–13)"
	}
	s := "⏱️  Verfügbarkeitsmodus wählen\n\n"
	s += o1 + "\n" + o2 + "\n\n"
	s += "←/→ oder ↑/↓: wählen · Enter: weiter · Backspace: zurück · Esc/Ctrl+C: beenden\n"
	return s
}

/*
func viewAvailDetail(m Model) string {
	// TODO: je nach mode die Felder anzeigen (dateISO / weekday + from/to)
	return ""
}

func viewReview(m Model) string {
	// TODO: Zusammenfassung + "Enter: bestätigen · Backspace: zurück"
	return ""
}:
*/

func viewAvailDetail(m Model) string { return "Verfügbarkeitsdetails (kommen später)\n" }
func viewReview(m Model) string      { return "Review (kommt später)\n" }
