package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func expandUserPath(path string) (string, error) {
	if path != "~" && (len(path) < 2 || path[0] != '~' || (path[1] != '/' && path[1] != '\\')) {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locate user home directory: %w", err)
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}
