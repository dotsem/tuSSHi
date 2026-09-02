package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandTilde replaces ~/ prefix with the user home directory path.
func ExpandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
