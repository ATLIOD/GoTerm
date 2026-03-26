package internal

type entry struct {
	name  string
	isDir bool
}

type AppState struct {
	Cwd        string
	Entries    []entry
	Cursor     int
	Err        string
	Width      int
	Height     int
	ShowHidden bool
}
