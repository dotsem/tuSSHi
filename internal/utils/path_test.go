package utils_test

import (
	"os"
	"path/filepath"
	"testing"
	"tusshi/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path",
			input:    "/var/log/baboon",
			expected: "/var/log/baboon",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "home directory prefix",
			input:    "~/Documents",
			expected: filepath.Join(home, "Documents"),
		},
		{
			name:     "ssh config path",
			input:    "~/.ssh/config",
			expected: filepath.Join(home, ".ssh/config"),
		},
		{
			name:     "tilde without slash",
			input:    "~",
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.ExpandTilde(tt.input))
		})
	}
}
