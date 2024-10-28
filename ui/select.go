package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Select interface {
	Run() (string, error)
}

type selectModel struct {
	options  []string
	selected string

	cursor     int
	pageCursor int
	pageSize   int

	searchMode      bool
	searchTerm      string
	filteredOptions []string
}

func NewSelect(options []string) Select {
	return selectModel{
		options:  options,
		selected: "",

		cursor:     0,
		pageCursor: 0,
		pageSize:   min(len(options), 7),

		searchMode:      false,
		searchTerm:      "",
		filteredOptions: options,
	}
}

func (m selectModel) Run() (string, error) {
	program := tea.NewProgram(m)

	model, err := program.Run()
	if err != nil {
		fmt.Printf("Error selecting from options: %v", err)
		return "", err
	}

	return model.(selectModel).selected, nil
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) updateFilteredOptions() []string {
	filteredOptions := make([]string, 0)
	for _, option := range m.options {
		if strings.Contains(strings.ToLower(option), strings.ToLower(strings.TrimSpace(m.searchTerm))) {
			filteredOptions = append(filteredOptions, option)
		}
	}
	return filteredOptions
}

func (m selectModel) updateSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case " ", "enter":
			m.selected = m.filteredOptions[m.cursor]
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

func (m selectModel) updateSelectMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case " ", "enter":
			m.selected = m.filteredOptions[m.cursor]
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
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

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.searchMode {
		return m.updateSearchMode(msg)
	}

	return m.updateSelectMode(msg)
}

func (m selectModel) View() string {
	s := "What .gitignore template do you want to include?\n\n"

	if m.searchMode {
		s += fmt.Sprintf("Search: %s\n", m.searchTerm)
	}

	for i := m.pageCursor; i < min(m.pageCursor+m.pageSize, len(m.filteredOptions)); i++ {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, m.filteredOptions[i])
	}

	if m.searchMode {
		s += "\n<Press ctrl+c to quit; s to search; space to select; enter to continue>"
	} else {
		s += "\n<Press ctrl+c to quit; esc to exit search; space to select; enter to continue>"
	}

	return s
}
