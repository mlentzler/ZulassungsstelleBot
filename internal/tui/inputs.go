package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

func initInputs(m *Model) {
	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "Vollständiger Name"
	m.nameInput.CharLimit = 80
	m.nameInput.Width = 40
	m.nameInput.Focus()

	m.emailInput = textinput.New()
	m.emailInput.Placeholder = "E-Mail-Adresse"
	m.emailInput.Width = 40
	// m.emailInput.SetValue("")

	m.phoneInput = textinput.New()
	m.phoneInput.Placeholder = "Telefonnummer"
	m.phoneInput.Width = 40
}

func focusName(m *Model) {
	m.nameInput.Focus()
	m.emailInput.Blur()
	m.phoneInput.Blur()
}

func focusEmail(m *Model) {
	m.nameInput.Blur()
	m.emailInput.Focus()
	m.phoneInput.Blur()
}

func focusPhone(m *Model) {
	m.nameInput.Blur()
	m.emailInput.Focus()
	m.emailInput.Blur()
	m.phoneInput.Focus()
}

func validateName(s string) error {
	if len(strings.TrimSpace(s)) < 2 {
		return errInvalid("Name zu kurz")
	}
	return nil
}

var reMail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func validateEmail(s string) error {
	if !reMail.MatchString(strings.TrimSpace(s)) {
		return errInvalid("Ungültige E-Mail")
	}
	return nil
}

var rePhone = regexp.MustCompile(`^[0-9 +()/-]{6,}$`)

func validatePhone(s string) error {
	if !rePhone.MatchString(strings.TrimSpace(s)) {
		return errInvalid("Ungültige Telefonnummer")
	}
	return nil
}

type invalidErr string

func (e invalidErr) Error() string { return string(e) }
func errInvalid(msg string) error  { return invalidErr(msg) }
