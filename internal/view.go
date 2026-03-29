package internal

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
	s := "GoTerm File Manager\n\n"

	for i, entry := range m.Entries {

		cursor := " "
		// Is the cursor pointing at this choice?
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, entry)
	}

	s += "\nPress q to quit.\n"

	return s
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
