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
		s += f0 + "Datum (YYYY-MM-DD): " + m.dateInput.View() + "\n\n"
		s += f1 + "From Stunde (0-23): " + m.fromInput.View() + "\n\n"
		s += f2 + "To   Stunde (1-24): " + m.toInput.View() + "\n\n"
		if m.errMsg != "" {
			s += "⚠️  " + m.errMsg + "\n\n"
		}
		s += "Tab/Shift+Tab: Feld wechseln · Enter: weiter · Esc: zurück · q/Ctrl+C: beenden\n"
		return s
	}

	var b strings.Builder
	b.WriteString("🔁 Wöchentliche Verfügbarkeit\n\n")
	b.WriteString("  Leertaste: Tag an/aus · ↑/↓: Zeile · ←/→ oder Tab: Feld · Enter: weiter · Esc: zurück · q/Ctrl+C: beenden\n\n")

	for i, wd := range weekdays {
		rowSel := "  "
		if i == m.recCursor {
			rowSel = "> "
		}
		toggle := "[ ]"
		if m.recSelected[i] {
			toggle = "[x]"
		}
		tf := "  "
		ff := "  "
		of := "  "
		if i == m.recCursor {
			switch m.recField {
			case 0:
				tf = "★ "
			case 1:
				ff = "★ "
			case 2:
				of = "★ "
			}
		}

		fmt.Fprintf(&b, "%s%s%s  %-2s   %sFrom: %s   %sTo: %s\n",
			rowSel, tf, toggle, wd,
			ff, m.recFromInputs[i].View(),
			of, m.recToInputs[i].View(),
		)
	}
	if m.errMsg != "" {
		b.WriteString("\n⚠️  " + m.errMsg + "\n")
	}
	return b.String()
}

/*
func viewReview(m Model) string {
	// TODO: Zusammenfassung + "Enter: bestätigen · Backspace: zurück"
	return ""
}:
*/

func viewReview(m Model) string { return "Review (kommt später)\n" }
