package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type MultiSelect interface {
	Run() ([]string, error)
}

type multiSelectModel struct {
	options  []string
	selected map[string]bool

	cursor     int
	pageCursor int
	pageSize   int

	searchMode      bool
	searchTerm      string
	filteredOptions []string
}

func NewMultiSelect(options []string) MultiSelect {
	return multiSelectModel{
		options:  options,
		selected: make(map[string]bool),

		cursor:     0,
		pageCursor: 0,
		pageSize:   min(len(options), 7),

		searchMode:      false,
		searchTerm:      "",
		filteredOptions: options,
	}
}

func (m multiSelectModel) Run() ([]string, error) {
	program := tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error selecting from options: %v", err)
		return nil, err
	}

	selected := make([]string, 0)
	for key := range m.selected {
		selected = append(selected, key)
	}

	return selected, nil
}

func (m multiSelectModel) Init() tea.Cmd {
	return nil
}

func (m multiSelectModel) updateFilteredOptions() []string {
	filteredOptions := make([]string, 0)
	for _, option := range m.options {
		if strings.Contains(strings.ToLower(option), strings.ToLower(strings.TrimSpace(m.searchTerm))) {
			filteredOptions = append(filteredOptions, option)
		}
	}
	return filteredOptions
}

func (m multiSelectModel) updateSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		case "esc":
			m.searchTerm = ""
			m.searchMode = false
			m.filteredOptions = m.updateFilteredOptions()
			m.cursor = 0
			m.pageCursor = 0
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "left", "right":
		case " ":
			if _, ok := m.selected[m.filteredOptions[m.cursor]]; ok {
				m.selected[m.filteredOptions[m.cursor]] = !m.selected[m.filteredOptions[m.cursor]]
			} else {
				m.selected[m.filteredOptions[m.cursor]] = true
			}
		case "backspace":
			if len(m.searchTerm) > 0 {
				m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
				m.filteredOptions = m.updateFilteredOptions()
				m.cursor = 0
				m.pageCursor = 0
			}
		default:
			m.searchTerm += msg.String()
			m.filteredOptions = m.updateFilteredOptions()
			m.cursor = 0
			m.pageCursor = 0
		}
	}

	return m, nil
}

func (m multiSelectModel) updateSelectMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case " ":
			if _, ok := m.selected[m.filteredOptions[m.cursor]]; ok {
				m.selected[m.filteredOptions[m.cursor]] = !m.selected[m.filteredOptions[m.cursor]]
			} else {
				m.selected[m.filteredOptions[m.cursor]] = true
			}
		case "s":
			m.searchMode = true
		}
	}

	if m.cursor-m.pageCursor >= m.pageSize-1 {
		m.pageCursor = min(m.pageCursor+1, len(m.options)-m.pageSize)
	} else if m.cursor-m.pageCursor <= 0 {
		m.pageCursor = max(0, m.pageCursor-1)
	}

	return m, nil
}

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.searchMode {
		return m.updateSearchMode(msg)
	}

	return m.updateSelectMode(msg)
}

func (m multiSelectModel) View() string {
	s := "What .gitignore template do you want to include?\n\n"

	if m.searchMode {
		s += fmt.Sprintf("Search: %s\n", m.searchTerm)
	}

	for i := m.pageCursor; i < min(m.pageCursor+m.pageSize, len(m.filteredOptions)); i++ {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		selected := " "
		if isSelected, ok := m.selected[m.filteredOptions[i]]; ok && isSelected {
			selected = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, selected, m.filteredOptions[i])
	}

	if m.searchMode {
		s += "\n<Press ctrl+c to quit; s to search; space to select; enter to continue>"
	} else {
		s += "\n<Press ctrl+c to quit; esc to exit search; space to select; enter to continue>"
	}

	return s
}
