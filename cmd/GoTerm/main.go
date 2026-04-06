package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"GoTerm/internal"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() internal.AppState {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	parentDir := filepath.Dir(cwd)

	ti := textinput.New()
	ti.Placeholder = "Enter filename..."
	ti.CharLimit = 256

	m := internal.AppState{
		Cwd:          cwd,
		Width:        80,
		Height:       24,
		ShowHidden:   false,
		PromptActive: false,
		TextInput:    ti,
		ParentDir:    parentDir,
	}

	m = m.Reload()
	return m
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error has occured: %v", err)
		os.Exit(1)
	}
}
