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
