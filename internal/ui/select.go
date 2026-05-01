package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

type selectModel struct {
	label    string
	options  []string
	cursor   int
	selected string
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.selected = m.options[m.cursor]
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	out := labelStyle.Render(m.label) + "\n"
	for i, opt := range m.options {
		if i == m.cursor {
			out += cursorStyle.Render("> ") + opt + "\n"
		} else {
			out += "  " + opt + "\n"
		}
	}
	return out
}

func Select(label string, options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}
	m, err := tea.NewProgram(selectModel{label: label, options: options}).Run()
	if err != nil {
		return "", err
	}
	sm := m.(selectModel)
	if sm.selected == "" {
		return "", fmt.Errorf("cancelled")
	}
	return sm.selected, nil
}
