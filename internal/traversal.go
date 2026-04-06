package internal

import "path/filepath"

func (m AppState) TraverseBack() AppState {
	parent := filepath.Clean(filepath.Join(m.Cwd, ".."))
	if parent == m.Cwd {
		return m
	}
	m.Cwd = parent
	m.ParentDir = filepath.Dir(m.Cwd)
	m.Cursor = 0
	m.Selection = entry{}
	return m.Reload()
}

func (m AppState) enterSelected() AppState {
	if len(m.Entries) == 0 {
		return m
	}
	e := m.Entries[m.Cursor]
	next := filepath.Join(m.Cwd, e.Name)
	if e.IsDir {
		m.ParentDir = m.Cwd
		m.Cwd = next
		m.Cursor = 0
		m.Selection = entry{}
		return m.Reload()
	}
	if err := openWithSystem(next); err != nil {
		m.Err = err.Error()
	}
	return m
}
