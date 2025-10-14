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

	mode        domain.AvailabilityKind
	availCursor int

	detailFocus int

	dateISO   string
	dateInput textinput.Model

	weekday       string
	weekdayCorsor int

	fromHour  int
	toHour    int
	fromInput textinput.Model
	toInput   textinput.Model

	result *domain.BookingRequest

	errMsg string
	cfg    config.Config
}

func NewModel(root domain.MenuNode, cfg config.Config) Model {
	var m Model
	m.cfg = cfg
	m.menuRoot = root
	m.step = stepPerson
	m.availCursor = 0
	m.detailFocus = 0
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
				m.step = stepMenu
				return m, nil
			}
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.nameInput.Focused() {
				focusEmail(&m)
				return m, nil
			}
			if m.emailInput.Focused() {
				focusPhone(&m)
				return m, nil
			}
			focusName(&m)
			return m, nil
		case "shift+tab":
			if m.phoneInput.Focused() {
				focusEmail(&m)
				return m, nil
			}
			if m.emailInput.Focused() {
				focusName(&m)
				return m, nil
			}
			focusPhone(&m)
			return m, nil
		case "up":
			if m.emailInput.Focused() {
				focusName(&m)
				return m, nil
			}
			if m.phoneInput.Focused() {
				focusEmail(&m)
				return m, nil
			}
		case "down":
			if m.nameInput.Focused() {
				focusEmail(&m)
				return m, nil
			}
			if m.emailInput.Focused() {
				focusPhone(&m)
				return m, nil
			}
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

func updateMenu(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	nodes := currentLevelNodes(&m)

	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
			return m, nil
		case "down", "j":
			if len(nodes) > 0 && m.menuCursor < len(nodes)-1 {
				m.menuCursor++
			}
			return m, nil
		case "enter":
			if len(nodes) == 0 {
				return m, nil
			}
			n := nodes[m.menuCursor]
			if len(n.Children) > 0 {
				pushMenu(&m, m.menuCursor)
				return m, nil
			}
			m.path = append(currentPathTitles(&m), n.Title)
			m.step = stepAvailabilityMode
			return m, nil
		case "backspace", "h":
			popMenu(&m)
			return m, nil
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateAvailMode(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "left", "h", "up", "k":
			if m.availCursor > 0 {
				m.availCursor--
			}
			return m, nil
		case "right", "l", "down", "j":
			if m.availCursor < 1 {
				m.availCursor++
			}
			return m, nil
		case "enter":
			if m.availCursor == 0 {
				m.mode = domain.AvailOneOff
				m.detailFocus = 0
				m.dateInput.Focus()
				m.fromInput.Blur()
				m.toInput.Blur()
			} else {
				m.mode = domain.AvailRecurring
				m.detailFocus = 0
				m.dateInput.Blur()
				m.fromInput.Blur()
				m.toInput.Blur()
			}
			m.errMsg = ""
			m.step = stepAvailabilityDetail
			return m, nil
		case "backspace":
			m.step = stepMenu
			return m, nil
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateAvailDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func updateReview(m Model, msg tea.Msg) (tea.Model, tea.Cmd)      { return m, nil }
