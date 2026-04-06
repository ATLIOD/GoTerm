package internal

import (
	"fmt"
	"path/filepath"
	"strings"

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

func (m AppState) mainPanel() string {
	var b strings.Builder

	listHeight := m.Height - 8
	if listHeight < 3 {
		listHeight = 3
	}
	colW := m.Width - 2
	if colW < 10 {
		colW = 10
	}

	start := 0
	if len(m.Entries) > listHeight && m.Cursor >= listHeight/2 {
		start = m.Cursor - listHeight/2
		if start+listHeight > len(m.Entries) {
			start = len(m.Entries) - listHeight
		}
		if start < 0 {
			start = 0
		}
	}

	if len(m.Entries) == 0 {
		b.WriteString(HelpStyle.Render("(empty directory)"))
		b.WriteString("\n")
	} else {
		end := start + listHeight
		if end > len(m.Entries) {
			end = len(m.Entries)
		}
		for i := start; i < end; i++ {
			e := m.Entries[i]
			cursor := " "
			if m.Cursor == i {
				cursor = "›"
			}
			suffix := ""
			if e.IsDir {
				suffix = "/"
			}
			name := e.Name + suffix
			name = Truncate(name, colW-4)

			line := fmt.Sprintf("%s %s", cursor, name)
			var styled string
			if m.Cursor == i {
				styled = SelStyle.Render(line)
			} else {
				if e.IsDir {
					styled = DirStyle.Render(line)
				} else {
					styled = FileStyle.Render(line)
				}
			}
			b.WriteString(styled)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m AppState) leftPanel() string {
	var b strings.Builder

	listHeight := m.Height - 8
	if listHeight < 3 {
		listHeight = 3
	}
	colW := m.Width - 2
	if colW < 10 {
		colW = 10
	}

	start := 0

	if len(m.ParentEntries) == 0 {
		b.WriteString(HelpStyle.Render("(empty directory)"))
		b.WriteString("\n")
	} else {
		end := start + listHeight
		if end > len(m.ParentEntries) {
			end = len(m.ParentEntries)
		}
		for i := start; i < end; i++ {
			e := m.ParentEntries[i]
			suffix := ""
			if e.IsDir {
				suffix = "/"
			}
			name := e.Name + suffix
			name = Truncate(name, colW-4)

			line := fmt.Sprintf("%s ", name)
			var styled string
			if filepath.Base(m.Cwd) == e.Name {
				styled = SelStyle.Render(line)
			} else {
				if e.IsDir {
					styled = DirStyle.Render(line)
				} else {
					styled = FileStyle.Render(line)
				}
			}
			b.WriteString(styled)
			b.WriteString("\n")
		}
	}
	return b.String()
}
