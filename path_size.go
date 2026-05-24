package code

import (
	"code/internal/pathsize"
	"code/internal/sizefmt"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := pathsize.Calculate(path, all, recursive)
	if err != nil {
		return "", err
	}

	if human {
		return sizefmt.BytesIEC(size), nil
	}

	return sizefmt.BytesRaw(size), nil
}
