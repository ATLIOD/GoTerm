package internal

import (
	"fmt"
	"io"
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

func (m AppState) yankSelection() AppState {
	if len(m.Entries) == 0 {
		m.Err = "Nothing to copy"
		return m
	}

	m.Clipboard = []entry{m.Selection}
	m.Err = fmt.Sprintf("Copied %s", m.Selection.Name)
	return m
}

func (m AppState) startPasteIntoCurrentDir() AppState {
	if len(m.Clipboard) == 0 {
		m.Err = "Clipboard is empty"
		return m
	}

	source := m.Clipboard[0]
	destPath := filepath.Join(m.Cwd, source.Name)
	return m.copyWithConflictCheck(source.Path, destPath)
}

func (m AppState) copyWithConflictCheck(sourcePath, destPath string) AppState {
	if sourcePath == "" || destPath == "" {
		m.Err = "Invalid copy paths"
		return m
	}

	if sourcePath == destPath {
		m.Err = "Source and destination are the same"
		return m
	}

	if _, err := os.Lstat(destPath); err == nil {
		m.ConfirmActive = true
		m.ConfirmMessage = fmt.Sprintf("Overwrite %s? (y/n)", filepath.Base(destPath))
		m.PendingSource = sourcePath
		m.PendingDest = destPath
		m.Err = ""
		return m
	} else if !os.IsNotExist(err) {
		m.Err = fmt.Sprintf("Could not access destination: %v", err)
		return m
	}

	return m.executeCopy(sourcePath, destPath, false)
}

func (m AppState) resolveOverwriteConfirm(overwrite bool) AppState {
	sourcePath := m.PendingSource
	destPath := m.PendingDest

	m.ConfirmActive = false
	m.ConfirmMessage = ""
	m.PendingSource = ""
	m.PendingDest = ""

	if !overwrite {
		m.Err = "Paste skipped"
		return m
	}

	return m.executeCopy(sourcePath, destPath, true)
}

func (m AppState) executeCopy(sourcePath, destPath string, overwrite bool) AppState {
	if overwrite {
		if err := os.RemoveAll(destPath); err != nil {
			m.Err = fmt.Sprintf("Could not overwrite destination: %v", err)
			return m
		}
	}

	if err := copyPath(sourcePath, destPath); err != nil {
		m.Err = fmt.Sprintf("Copy failed: %v", err)
		return m
	}

	m.Err = fmt.Sprintf("Pasted %s", filepath.Base(sourcePath))
	return m
}

func copyPath(sourcePath, destPath string) error {
	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	switch {
	case info.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(sourcePath)
		if err != nil {
			return err
		}
		return os.Symlink(target, destPath)
	case info.IsDir():
		return copyDir(sourcePath, destPath, info.Mode().Perm())
	default:
		return copyFile(sourcePath, destPath, info.Mode().Perm())
	}
}

func copyDir(sourceDir, destDir string, mode os.FileMode) error {
	if err := os.MkdirAll(destDir, mode); err != nil {
		return err
	}

	items, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		sourcePath := filepath.Join(sourceDir, item.Name())
		destPath := filepath.Join(destDir, item.Name())
		if err := copyPath(sourcePath, destPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(sourceFile, destFile string, mode os.FileMode) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return err
	}

	return nil
}
