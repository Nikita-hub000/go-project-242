// Package code provides file and directory size calculation utilities.
package code

import (
	"code/internal/pathsize"
	"code/internal/sizefmt"
)

// GetPathSize returns the formatted size of path.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := pathsize.Calculate(path, all, recursive)
	if err != nil {
		return "", err
	}

	if human {
		return sizefmt.FormatIEC(size), nil
	}

	return sizefmt.FormatRaw(size), nil
}
