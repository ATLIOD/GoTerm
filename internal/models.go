package internal

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type entry struct {
	Name  string
	IsDir bool
}

type AppState struct {
	Cwd          string
	Entries      []entry
	Cursor       int
	Err          string
	Width        int
	Height       int
	ShowHidden   bool
	PromptActive bool
	TextInput    textinput.Model
}
