package internal

func (m AppState) HandleAction(input string) {
	switch m.Action {
	case NewFile:
		m = m.newFile(input)
	case NewDirectory:
		m = m.newDirectory(input)
	}
	m.Action = None
}

func (m AppState) GetPrompt() string {
	switch m.Action {
	case NewFile:
		return "Enter new file name:"
	case NewDirectory:
		return "Enter new directory name:"
	default:
		return ""
	}
}
