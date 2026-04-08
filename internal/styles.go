package internal

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	DirStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	FileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	SelStyle  = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("255"))
	ErrStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	ColStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("240")).
			PaddingLeft(1)

	PreviewStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("244"))
)
