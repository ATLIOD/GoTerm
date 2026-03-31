// Package internal holds all the special bits
package internal

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m AppState) Init() tea.Cmd {
	return nil
}

func (m AppState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Entries)-1 {
				m.Cursor++
			}

		case "enter", "l", "o":
			return m.enterSelected(), nil

		case "backspace", "left", "h":
			return m.TraverseBack(), nil

		case "home", "g":
			m.Cursor = 0

		case "end", "G":
			if len(m.Entries) > 0 {
				m.Cursor = len(m.Entries) - 1
			}

		case "r":
			return m.Reload(), nil

		case ".":
			m.ShowHidden = !m.ShowHidden
			m.Cursor = 0
			return m.Reload(), nil

		case "ctrl+n":
			return m.newFile(), nil

		case "ctrl+shift+n":
			return m.newFolder(), nil

		}
	}

	return m, nil
}

func (m AppState) View() string {
	pathLine := m.Cwd
	maxPath := m.Width - 4
	if maxPath < 8 {
		maxPath = 8
	}
	pathLine = Truncate(pathLine, maxPath)

	var b strings.Builder
	b.WriteString(TitleStyle.Render(" GoTerm — file manager "))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(pathLine))
	b.WriteString("\n\n")

	if m.Err != "" {
		b.WriteString(ErrStyle.Render(m.Err))
		b.WriteString("\n\n")
	}

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

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(
		Truncate("j/k move · Enter/l dir or open · h/← parent · r refresh · . hidden · q quit", m.Width-1),
	))
	b.WriteString("\n")

	return b.String()
}

func (m AppState) Reload() AppState {
	entries, err := loadEntries(m.Cwd, m.ShowHidden)
	if err != nil {
		m.Err = err.Error()
		return m
	}
	m.Entries = entries
	m.Err = ""
	if len(m.Entries) == 0 {
		m.Cursor = 0
	} else if m.Cursor >= len(m.Entries) {
		m.Cursor = len(m.Entries) - 1
	}
	return m
}
