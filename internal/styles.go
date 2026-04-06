package internal

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	DirStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	FileStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	SelStyle   = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("255"))
	ErrStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	HelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)
