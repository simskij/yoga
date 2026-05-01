package ui

import "github.com/charmbracelet/lipgloss"

var (
	Green  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	Red    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Bold   = lipgloss.NewStyle().Bold(true)
	Dim    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
