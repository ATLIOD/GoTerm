package internal

import "path/filepath"

func (m AppState) TraverseBack() AppState {
	parent := filepath.Clean(filepath.Join(m.Cwd, ".."))
	if parent == m.Cwd {
		return m
	}
	m.Cwd = parent
	m.Cursor = 0
	return m.Reload()
}
