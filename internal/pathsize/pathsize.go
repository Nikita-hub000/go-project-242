package pathsize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Calculate(path string, includeHidden, recursive bool) (int64, error) {
	return calculate(path, includeHidden, recursive)
}

func calculate(path string, includeHidden, recursive bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("stat path %q: %w", path, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return info.Size(), nil
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("read directory %q: %w", path, err)
	}

	var total int64

	for _, entry := range entries {
		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(path, entry.Name())
		entryInfo, err := os.Lstat(entryPath)
		if err != nil {
			return 0, fmt.Errorf("stat entry %q: %w", entryPath, err)
		}

		if entryInfo.Mode()&os.ModeSymlink != 0 {
			total += entryInfo.Size()
			continue
		}

		if entryInfo.IsDir() {
			if !recursive {
				continue
			}

			size, err := calculate(entryPath, includeHidden, recursive)
			if err != nil {
				return 0, err
			}

			total += size
			continue
		}

		total += entryInfo.Size()
	}

	return total, nil
}
