package tui

import (
	"fmt"
	"strings"
	"time"

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

const (
	carArt = "" +
		"                         _.-=\"_-         _\n" +
		"                         _.-=\"   _-          | ||\"\"\"\"\"\"\"---._______     __..\n" +
		"             ___.===\"\"\"\"-.______-,,,,,,,,,,,,`-''----\" \"\"\"\"\"       \"\"\"\"\"  __'\n" +
		"      __.--\"\"     __        ,'                   o \\           __        [__|\n" +
		" __-\"\"=======.--\"\"  \"\"--.=================================.--\"\"  \"\"--.=======:\n" +
		"]       [w] : /        \\ : |========================|    : /        \\ :  [w] :\n" +
		"V___________:|          |: |========================|    :|          |:   _-\"\n" +
		" V__________: \\        / :_|=======================/_____: \\        / :__-\"\n" +
		" -----------'  \"-_____-\"  `-------------------------------'  \"-_____-\"\n"
)

type Model struct {
	step step

	nameInput  textinput.Model
	emailInput textinput.Model
	phoneInput textinput.Model

	menuRoot      domain.MenuNode
	menuStack     []int
	menuCursor    int
	path          []string
	menuSelectors []string

	mode        domain.AvailabilityKind
	availCursor int

	// ---- One-Off detail ----
	detailFocus   int
	dateISO       string
	dateInput     textinput.Model
	weekday       string
	weekdayCursor int
	fromHour      int
	toHour        int
	fromInput     textinput.Model
	toInput       textinput.Model

	// ---- Recurring detail ----
	recCursor     int
	recField      int
	recSelected   [7]bool
	recFromInputs [7]textinput.Model
	recToInputs   [7]textinput.Model
	recDays       []domain.DayWindow

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
	initAvailInputs(&m)
	initAvailInputs(&m)
	m.availCursor = 0
	m.detailFocus = 0
	m.recCursor = 0
	m.recField = 0
	return m
}

func (m Model) Init() tea.Cmd { return nil }

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
	var s strings.Builder
	s.WriteString(carArt)
	s.WriteString("\n\n")

	switch m.step {
	case stepPerson:
		s.WriteString(viewPerson(m))
	case stepMenu:
		s.WriteString(viewMenu(m))
	case stepAvailabilityMode:
		s.WriteString(viewAvailMode(m))
	case stepAvailabilityDetail:
		s.WriteString(viewAvailDetail(m))
	case stepReview:
		s.WriteString(viewReview(m))
	case stepDone:
		s.WriteString("fertig\n")
	}
	return s.String()
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
		case "esc":
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
		case "shift + tab":
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
		case "enter", "l", "right":
			if len(nodes) == 0 {
				return m, nil
			}
			n := nodes[m.menuCursor]
			if len(n.Children) > 0 {
				pushMenu(&m, m.menuCursor)
				return m, nil
			}
			m.path = append(currentPathTitles(&m), n.Title)
			selectors := currentPathSelectors(&m, m.menuCursor)
			m.menuSelectors = selectors
			m.step = stepAvailabilityMode
			return m, nil
		case "left", "h":
			popMenu(&m)
			return m, nil
		case "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateAvailMode(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dateInput.Placeholder == "" && m.fromInput.Placeholder == "" && m.toInput.Placeholder == "" {
		initAvailInputs(&m)
	}

	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "up", "k":
			if m.availCursor > 0 {
				m.availCursor--
			}
			return m, nil
		case "down", "j":
			if m.availCursor < 1 {
				m.availCursor++
			}
			return m, nil
		case "enter":
			//One-Off
			if m.availCursor == 0 {
				m.mode = domain.AvailOneOff
				m.detailFocus = 0
				cmd := m.dateInput.Focus()
				m.fromInput.Blur()
				m.toInput.Blur()
				m.errMsg = ""
				m.step = stepAvailabilityDetail
				return m, cmd
			}
			// Recurring
			m.mode = domain.AvailRecurring
			m.recCursor = 0
			m.recField = 0
			for i := 0; i < 7; i++ {
				m.recFromInputs[i].Blur()
				m.recToInputs[i].Blur()
			}
			m.errMsg = ""
			m.step = stepAvailabilityDetail
			return m, nil

		case "left", "h":
			m.step = stepMenu
			return m, nil
		case "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateAvailDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case domain.AvailOneOff:
		return updateOneOffDetail(m, msg)
	case domain.AvailRecurring:
		return updateRecurringDetail(m, msg)
	default:
		return m, nil
	}
}

func updateOneOffDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "tab":
			m.detailFocus = (m.detailFocus + 1) % 3
			switch m.detailFocus {
			case 0:
				m.dateInput.Focus()
				m.fromInput.Blur()
				m.toInput.Blur()
			case 1:
				m.dateInput.Blur()
				m.fromInput.Focus()
				m.toInput.Blur()
			case 2:
				m.dateInput.Blur()
				m.fromInput.Blur()
				m.toInput.Focus()
			}
			return m, nil

		case "shift+tab":
			m.detailFocus = (m.detailFocus + 2) % 3
			switch m.detailFocus {
			case 0:
				m.dateInput.Focus()
				m.fromInput.Blur()
				m.toInput.Blur()
			case 1:
				m.dateInput.Blur()
				m.fromInput.Focus()
				m.toInput.Blur()
			case 2:
				m.dateInput.Blur()
				m.fromInput.Blur()
				m.toInput.Focus()
			}
			return m, nil

		case "enter":
			var err error
			if _, err = validateDateEU(m.dateInput.Value()); err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			fh, fe := parseHour(m.fromInput.Value())
			th, te := parseHour(m.toInput.Value())
			if fe != nil {
				m.errMsg = fe.Error()
				return m, nil
			}
			if te != nil {
				m.errMsg = te.Error()
				return m, nil
			}
			if err = validateHours(fh, th); err != nil {
				m.errMsg = err.Error()
				return m, nil
			}

			m.fromHour, m.toHour = fh, th
			m.dateISO = func(s string) string {
				s = strings.TrimSpace(s)
				if t, err := time.Parse("02.01.2006", s); err == nil {
					return t.Format("2006-01-02")
				}
				if t, err := time.Parse("2.1.2006", s); err == nil {
					return t.Format("2006-01-02")
				}
				return ""
			}(m.dateInput.Value())
			m.errMsg = ""
			m.step = stepReview
			return m, nil

		case "esc":
			m.step = stepAvailabilityMode
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.detailFocus {
	case 0:
		m.dateInput, cmd = m.dateInput.Update(msg)
	case 1:
		m.fromInput, cmd = m.fromInput.Update(msg)
	default:
		m.toInput, cmd = m.toInput.Update(msg)
	}
	return m, cmd
}

func updateRecurringDetail(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "up", "k":
			if m.recCursor > 0 {
				m.recCursor--
			}
			return m, nil
		case "down", "j":
			if m.recCursor < 6 {
				m.recCursor++
			}
			return m, nil

		case "left", "h":
			if m.recField > 0 {
				m.recField--
			}
			setRecFocus(&m)
			return m, nil
		case "right", "l", "tab":
			m.recField = (m.recField + 1) % 3
			setRecFocus(&m)
			return m, nil
		case "shift+tab":
			m.recField = (m.recField + 2) % 3
			setRecFocus(&m)
			return m, nil

		case " ", "space":
			if m.recField == 0 {
				m.recSelected[m.recCursor] = !m.recSelected[m.recCursor]
				return m, nil
			}

		case "enter":
			has := false
			for i := 0; i < 7; i++ {
				if m.recSelected[i] {
					has = true
					break
				}
			}
			if !has {
				m.errMsg = "Bitte mindestens einen Wochentag auswählen"
				return m, nil
			}

			var out []domain.DayWindow
			for i := 0; i < 7; i++ {
				if !m.recSelected[i] {
					continue
				}
				fh, fe := parseHour(m.recFromInputs[i].Value())
				th, te := parseHour(m.recToInputs[i].Value())
				if fe != nil {
					m.errMsg = fmt.Sprintf("%s: %s", weekdays[i], fe.Error())
					return m, nil
				}
				if te != nil {
					m.errMsg = fmt.Sprintf("%s: %s", weekdays[i], te.Error())
					return m, nil
				}
				if err := validateHours(fh, th); err != nil {
					m.errMsg = fmt.Sprintf("%s: %s", weekdays[i], err.Error())
					return m, nil
				}
				out = append(out, domain.DayWindow{
					Weekday:  weekdays[i],
					FromHour: fh,
					ToHour:   th,
				})
			}
			m.recDays = out
			m.errMsg = ""
			m.step = stepReview
			return m, nil

		case "esc":
			m.step = stepAvailabilityMode
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	if m.recField == 1 {
		m.recFromInputs[m.recCursor], _ = m.recFromInputs[m.recCursor].Update(msg)
		return m, nil
	}
	if m.recField == 2 {
		m.recToInputs[m.recCursor], _ = m.recToInputs[m.recCursor].Update(msg)
		return m, nil
	}
	return m, nil
}

func setRecFocus(m *Model) {
	for i := 0; i < 7; i++ {
		m.recFromInputs[i].Blur()
		m.recToInputs[i].Blur()
	}
	switch m.recField {
	case 1:
		m.recFromInputs[m.recCursor].Focus()
	case 2:
		m.recToInputs[m.recCursor].Focus()
	}
}

func updateReview(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "enter":
			if m.mode == domain.AvailRecurring && len(m.recDays) == 0 {
				m.errMsg = "Keine Tage ausgewählt – bitte im vorherigen Schritt mindestens einen Tag aktivieren."
				return m, nil
			}

			br := domain.BookingRequest{
				Name:  m.nameInput.Value(),
				Email: m.emailInput.Value(),
				Phone: m.phoneInput.Value(),
				Menu: domain.MenuChoice{
					Path:      append([]string{}, m.path...),
					Selectors: append([]string{}, m.menuSelectors...),
				},
				TZ: m.cfg.TZ,
			}

			if m.mode == domain.AvailOneOff {
				br.Avail = domain.Availability{
					Kind: domain.AvailOneOff,
					OneOff: &domain.OneOff{
						DateISO:  m.dateISO,
						FromHour: m.fromHour,
						ToHour:   m.toHour,
					},
				}
			} else {
				br.Avail = domain.Availability{
					Kind:      domain.AvailRecurring,
					Recurring: &domain.Recurring{Days: append([]domain.DayWindow{}, m.recDays...)},
				}
			}

			m.result = &br
			m.step = stepDone
			return m, tea.Quit

		case "esc":
			m.step = stepAvailabilityDetail
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}
