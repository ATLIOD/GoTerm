package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func loadEntries(dir string, showHidden bool) ([]entry, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var out []entry
	for _, name := range names {
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(dir, name)
		info, err := os.Lstat(full)
		if err != nil {
			continue
		}
		out = append(out, entry{Name: name, IsDir: info.IsDir(), Path: full})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func openWithSystem(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func Truncate(s string, max int) string {
	if max <= 3 || len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func ValidateFileName(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("file name cannot be '.' or '..'")
	}
	if strings.ContainsAny(name, "/\x00") {
		return fmt.Errorf("file name cannot contain '/' or null bytes")
	}
	if len(name) > 255 {
		return fmt.Errorf("file name cannot exceed 255 bytes")
	}
	return nil
}

func ValidateDirectoryPath(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	if strings.ContainsRune(path, '\x00') {
		return fmt.Errorf("directory path cannot contain null bytes")
	}
	for _, component := range strings.Split(path, "/") {
		if component == "" {
			continue
		}
		if len(component) > 255 {
			return fmt.Errorf("path component %q exceeds 255 bytes", component)
		}
	}
	return nil
}

func readFileContents(path string, maxWidth int, maxHeight int) (string, error) {
	file, err := os.Open(path) // use the `path` parameter
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var builder strings.Builder

	for i := 0; i < maxHeight; i++ {
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			// Apply maxWidth truncation if specified
			if maxWidth > 0 && len(line) > maxWidth {
				line = line[:maxWidth]
			}
			builder.WriteString(line)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading file: %w", err)
		}
	}

	return builder.String(), nil
}
