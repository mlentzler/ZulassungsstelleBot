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

	nameInput  textinput.Model
	emailInput textinput.Model
	phoneInput textinput.Model

	menuRoot   domain.MenuNode
	menuStack  []int
	menuCursor int
	path       []string

	mode     domain.AvailabilityKind
	dateISO  string
	weekday  string
	fromHour int
	toHour   int

	result *domain.BookingRequest

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

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepPerson:
		return updatePerson(m, msg)
	case stepMenu:
		return updateMenu(m, msg) // stub
	case stepAvailabilityMode:
		return updateAvailMode(m, msg) // stub
	case stepAvailabilityDetail:
		return updateAvailDetail(m, msg) // stub
	case stepReview:
		return updateReview(m, msg) // stub
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

func updatePerson(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "enter":
			if m.nameInput.Focused() {
				if err := validateName(m.nameInput.Value()); err != nil {
					m.errMsg = err.Error()
					return m, nil
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
					Name:  m.nameInput.Value(),
					Email: m.emailInput.Value(),
					Phone: m.phoneInput.Value(),
					TZ:    m.cfg.TZ,
				}
				m.step = stepDone
				return m, tea.Quit
			}
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
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

func updateMenu(m Model, msg tea.Msg) (tea.Model, tea.Cmd)        { return m, nil }
func updateAvailMode(m Model, msg tea.Msg) (tea.Model, tea.Cmd)   { return m, nil }
func updateAvailDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func updateReview(m Model, msg tea.Msg) (tea.Model, tea.Cmd)      { return m, nil }
