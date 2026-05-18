package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive bool, all bool, human bool) (string, error) {
	size, err := calculateSize(path, all, recursive)
	if err != nil {
		return "", err
	}

	if human {
		return byteCountIEC(size), nil
	}

	return fmt.Sprintf("%dB", size), nil
}

func calculateSize(path string, all bool, recursive bool) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64

	for _, entry := range entries {
		if !all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			if !recursive {
				continue
			}
		}

		size, err := calculateSize(fullPath, all, recursive)
		if err != nil {
			return 0, err
		}

		total += size
	}

	return total, nil
}

func byteCountIEC(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%dB", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.1fGB", float64(size)/(1024*1024*1024))
	}
}
