package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		envPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(envPath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("failed to get project root")
}
