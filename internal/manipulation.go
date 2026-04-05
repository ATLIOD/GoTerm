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

	if err := ValidateFileName(fileName); err != nil {
		m.Err = fmt.Sprintf("Invalid file name: %v", err)
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

func (m AppState) newDirectory(directoryName string) AppState {
	filepath := filepath.Join(m.Cwd, directoryName)

	if _, err := os.Stat(filepath); err == nil {
		m.Err = "Directory already exists"
		return m
	}

	if err := ValidateDirectoryPath(directoryName); err != nil {
		m.Err = fmt.Sprintf("Invalid directory name: %v", err)
		return m
	}

	err := os.MkdirAll(filepath, 0777)
	if err != nil {
		m.Err = fmt.Sprintf("Error creating directory: %v", err)
		return m
	}

	return m
}
