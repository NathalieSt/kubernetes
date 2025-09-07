package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory")
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
