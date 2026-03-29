package internal

import (
	"os"
	"path/filepath"
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
		out = append(out, entry{name: name, isDir: info.IsDir()})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].isDir != out[j].isDir {
			return out[i].isDir
		}
		return strings.ToLower(out[i].name) < strings.ToLower(out[j].name)
	})
	return out, nil
}
