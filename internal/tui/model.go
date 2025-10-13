package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mlentzler/ZulassungsstelleBot/internal/config"
	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

type step int

const (
	stepPerson step = iota
	stepMenu
	stepAvailabilityMode
	stepAvailabilityDetail
	stepReview
	stepDone
)

type Model struct {
	step step

	// Inputs
	nameInput  textinput.Model
	emailInput textinput.Model
	phoneInput textinput.Model

	// Menu
	menuRoot   domain.MenuNode
	menuStack  []int // Index je Ebene
	menuCursor int   // aktueller Index auf aktueller Ebene
	path       []string

	// Availability
	mode     domain.AvailabilityKind
	dateISO  string
	weekday  string
	fromHour int
	toHour   int

	// Ergebnis
	result *domain.BookingRequest

	// Misc
	errMsg string
	cfg    config.Config
}

func NewModel(root domain.MenuNode, cfg config.Config) Model {
	var m Model
	m.cfg = cfg
	m.menuRoot = root
	m.step = stepPerson
	initInputs(&m)
	return m
}

func (m Model) Init() tea.Cmd {
	// optional: Fokus erstes Feld setzen
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepPerson:
		return updatePerson(m, msg)
	case stepMenu:
		return updateMenu(m, msg)
	case stepAvailabilityMode:
		return updateAvailMode(m, msg)
	case stepAvailabilityDetail:
		return updateAvailDetail(m, msg)
	case stepReview:
		return updateReview(m, msg)
	case stepDone:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func updatePerson(m Model, msg tea.Msg) (tea.Model, tea. Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.nameInput.Focused() {
				if err := validateName(m.nameInput.Value()); err != nil {
					m.errMsg = err.Error()
					return m , nil
				}
				m.errMsg = ""
				m.nameInput.Blur()
				m.emailInput.Focus()
				return m, nil
			}
			if m.emailInput.Focused() {
				if err := validateEmail(m.emailInput.Value()); err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				m.errMsg = ""
				m.emailInput.Blur()
				m.phoneInput.Focus()
				return m, nil
			}
			if m.phoneInput.Focused() {
				if err := validatePhone(m.phoneInput.Value()); err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				m.errMsg = ""
				m.result = &domain.BookingRequest{
					Name: m.nameInput.Value(),
					Email: m.emailInput.Value(),
					Phone: m.phoneInput.Value(),
					TZ: m.cfg.TZ,
				}
				m.step = stepDone
				return m, tea.Quit()
			}
		case "ctrl+c", "esc":
			return m, tea.Quit()
	}

	var cmd tea.Cmd
		if m.nameInput.Focused() {
			m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
		}
		if m.emailInput.Focused() {
			m.emailInput, cmd = m.emailInput.Update(msg)
			return m, cmd
		}
		m.phoneInput, cmd = m.phoneInput.Update(msg)
		return m, cmd
}

func (m Model) View() string {
	switch m.step {
	case stepPerson:
		return viewPerson(m)
	case stepMenu:
		return viewMenu(m)
	case stepAvailabilityMode:
		return viewAvailMode(m)
	case stepAvailabilityDetail:
		return viewAvailDetail(m)
	case stepReview:
		return viewReview(m)
	case stepDone:
		return "fertig\n"
	default:
		return ""
	}
}

// --- TEMP: No-ops bis wir diese Steps bauen ---
func updateMenu(m Model, msg tea.Msg) (tea.Model, tea.Cmd)        { return m, nil }
func updateAvailMode(m Model, msg tea.Msg) (tea.Model, tea.Cmd)   { return m, nil }
func updateAvailDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func updateReview(m Model, msg tea.Msg) (tea.Model, tea.Cmd)      { return m, nil }
