package tui

import (
	"fmt"
	"strings"

	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
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
	s += "Enter: weiter · Esc: abbrechen\n"
	return s
}

func viewMenu(m Model) string {
	var b strings.Builder
	b.WriteString("📂 Auswahl-Menü\n")
	b.WriteString(breadcrumb(&m) + "\n\n")

	nodes := currentLevelNodes(&m)
	if len(nodes) == 0 {
		b.WriteString("(Keine Einträge auf dieser Ebene)\n\n")
		b.WriteString("↑/↓: bewegen · Enter/L: wählen · H: zurück · Esc: beenden\n")
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

	b.WriteString("\n↑/↓: bewegen · Enter/L: wählen · H: zurück · Esc: beenden\n")
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
	s += "←/→ oder ↑/↓: wählen · Enter/L: weiter · H: zurück · Esc: beenden\n"
	return s
}

func viewAvailDetail(m Model) string {
	if m.mode == domain.AvailOneOff {
		f0, f1, f2 := "  ", "  ", "  "
		if m.detailFocus == 0 {
			f0 = "➤ "
		}
		if m.detailFocus == 1 {
			f1 = "➤ "
		}
		if m.detailFocus == 2 {
			f2 = "➤ "
		}

		s := "📅 Einmaliger Termin\n\n"
		s += f0 + "Datum (DD.MM.YYYY): " + m.dateInput.View() + "\n\n"
		s += f1 + "From Stunde (0-23): " + m.fromInput.View() + "\n\n"
		s += f2 + "To Stunde (1-24):   " + m.toInput.View() + "\n\n"
		if m.errMsg != "" {
			s += "⚠️  " + m.errMsg + "\n\n"
		}
		s += "Tab/Ctrl+Tab: Feld wechseln · Enter: weiter · H: zurück · Esc: beenden\n"
		return s
	}

	// Recurring
	f0, f1, f2 := "  ", "  ", "  "
	if m.detailFocus == 0 {
		f0 = "➤ "
	}
	if m.detailFocus == 1 {
		f1 = "➤ "
	}
	if m.detailFocus == 2 {
		f2 = "➤ "
	}

	s := "🔁 Wöchentliche Verfügbarkeit\n\n"
	// Weekday-List (einfach)
	for i, wd := range weekdays {
		cur := "  "
		if i == m.weekdayCursor {
			cur = "● "
		}
		if m.detailFocus == 0 && i == m.weekdayCursor {
			s += f0 + cur + wd + "\n"
		} else {
			s += "  " + cur + wd + "\n"
		}
	}
	s += "\n"
	s += f1 + "From Stunde (0-23): " + m.fromInput.View() + "\n\n"
	s += f2 + "To Stunde (1-24):   " + m.toInput.View() + "\n\n"
	if m.errMsg != "" {
		s += "⚠️  " + m.errMsg + "\n\n"
	}
	s += "↑/↓: Wochentag (wenn markiert) · Tab: Feld wechseln · Enter: weiter · H: zurück · Esc: beenden\n"
	return s
}

/*
func viewReview(m Model) string {
	// TODO: Zusammenfassung + "Enter: bestätigen · Backspace: zurück"
	return ""
}:
*/

func viewReview(m Model) string { return "Review (kommt später)\n" }
