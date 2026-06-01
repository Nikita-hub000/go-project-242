package pathsize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Calculate(path string, includeHidden, recursive bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("stat path %q: %w", path, err)
	}

	return calculatePath(path, info, includeHidden, recursive)
}

func calculatePath(path string, info os.FileInfo, includeHidden, recursive bool) (int64, error) {
	if isSymlink(info) || !info.IsDir() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("read directory %q: %w", path, err)
	}

	var total int64

	for _, entry := range entries {
		if shouldSkipEntry(entry, includeHidden, recursive) {
			continue
		}

		size, err := calculateEntry(filepath.Join(path, entry.Name()), includeHidden, recursive)
		if err != nil {
			return 0, err
		}

		total += size
	}

	return total, nil
}

func calculateEntry(path string, includeHidden, recursive bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("stat entry %q: %w", path, err)
	}

	return calculatePath(path, info, includeHidden, recursive)
}

func shouldSkipEntry(entry os.DirEntry, includeHidden, recursive bool) bool {
	return (isHidden(entry) && !includeHidden) || (isDirectory(entry) && !recursive)
}

func isHidden(entry os.DirEntry) bool {
	return strings.HasPrefix(entry.Name(), ".")
}

func isDirectory(entry os.DirEntry) bool {
	return entry.IsDir()
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}
