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
		if m.PromptActive {
			switch msg.String() {

			case "enter":
				input := m.TextInput.Value()
				m.PromptActive = false
				m.TextInput.Reset()
				m.HandleAction(input)
				return m.Reload(), nil

			case "esc":
				m.PromptActive = false
				m.TextInput.Reset()
				return m, nil
			}
			var cmd tea.Cmd
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
				m.Selection = m.Entries[m.Cursor]
			}

		case "down", "j":
			if m.Cursor < len(m.Entries)-1 {
				m.Cursor++
				m.Selection = m.Entries[m.Cursor]
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
			m.Action = NewFile
			m.PromptActive = true
			m.TextInput.Focus()
			return m, nil

		case "alt+n":
			m.Action = NewDirectory
			m.PromptActive = true
			m.TextInput.Focus()
			return m, nil
		}
	}

	return m, nil
}

func (m AppState) View() string {
	if m.PromptActive {
		return fmt.Sprintf(m.GetPrompt()+"\n%s", m.TextInput.View())
	}

	pathLine := m.Selection.Path
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

	c := lipgloss.JoinHorizontal(0.1, m.leftPanel(), m.mainPanel(), m.rightPanel())
	b.WriteString(c)

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(
		Truncate("j/k move · Enter/l dir or open · h/← parent · r refresh · . hidden · q quit", m.Width-1),
	))
	b.WriteString("\n")

	return b.String()
}

func (m AppState) Reload() AppState {
	// Update current directory entries for main panel
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

	// Update parent entries for left panel
	entries, err = loadEntries(m.ParentDir, m.ShowHidden)
	if err != nil {
		m.Err = err.Error()
		return m
	}
	m.ParentEntries = entries

	// Update selection and its entries for right panel
	m.Selection = m.Entries[m.Cursor]
	if m.Selection.IsDir {
		entries, err = loadEntries(m.Selection.Path, m.ShowHidden)
		if err != nil {
			m.Err = err.Error()
			return m
		}
		m.SelectionEntries = entries
	}
	m.Err = ""
	return m
}
