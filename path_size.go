package code

import (
	"os"
	"path/filepath"
	"strings"
)
func GetPathSize(path string, all bool, recursive bool, human bool) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	answer := int64(0)

	for _, file := range files {
		if !all && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, file.Name())

		if file.IsDir() && !recursive {
			continue
		}

		size, err := GetPathSize(fullPath, all, recursive, human)
		if err != nil {
			return 0, err
		}

		answer += size
	}

	return answer, nil
}
