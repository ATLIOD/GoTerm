package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m AppState) newFile(fileName string) AppState {
	filepath := filepath.Join(m.Cwd, fileName)

	if _, err := os.Stat(filepath); err == nil {
		m.Err = "File already exists"
		return m
	}

	file, err := os.Create(filepath)
	if err != nil {
		m.Err = fmt.Sprintf("Error creating file: %v", err)
		return m
	}
	defer file.Close()

	return m
}
