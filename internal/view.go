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

		case "enter", " ":
			// _, ok := m.Selected[m.Cursor]
			// if ok {
			// 	delete(m.Selected, m.Cursor)
			// } else {
			// 	m.Selected[m.Cursor] = struct{}{}
			// }
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
	pathLine = pathLine

	var b strings.Builder
	b.WriteString(" GoTerm — file manager ")
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(pathLine))
	b.WriteString("\n\n")

	if m.Err != "" {
		b.WriteString(m.Err)
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
		b.WriteString("(empty directory)")
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
			if e.isDir {
				suffix = "/"
			}
			name := e.name + suffix
			name = name

			line := fmt.Sprintf("%s %s", cursor, name)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString("j/k move · Enter/l dir or open · h/← parent · r refresh · . hidden · q quit")

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
