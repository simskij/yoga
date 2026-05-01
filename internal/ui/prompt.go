package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	labelStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type promptModel struct {
	label     string
	input     textinput.Model
	submitted bool
	value     string
}

func newPromptModel(label, placeholder, defaultVal string) promptModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(defaultVal)
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	return promptModel{label: label, input: ti}
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.value = m.input.Value()
			if m.value == "" {
				m.value = m.input.Placeholder
			}
			m.submitted = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	return fmt.Sprintf(
		"%s\n%s %s\n",
		labelStyle.Render(m.label),
		promptStyle.Render(">"),
		m.input.View(),
	)
}

func Prompt(label, placeholder, defaultVal string) (string, error) {
	m, err := tea.NewProgram(newPromptModel(label, placeholder, defaultVal)).Run()
	if err != nil {
		return "", err
	}
	pm := m.(promptModel)
	if !pm.submitted {
		return "", fmt.Errorf("cancelled")
	}
	return pm.value, nil
}
