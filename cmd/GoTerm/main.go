package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"GoTerm/internal"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() internal.AppState {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	m := internal.AppState{
		Cwd:        cwd,
		Width:      80,
		Height:     24,
		ShowHidden: false,
	}
	return m
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
