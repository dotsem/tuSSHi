package ssh_test

import (
	"testing"

	"tusshi/internal/ssh"

	"github.com/stretchr/testify/assert"
)

func TestFindPreset(t *testing.T) {
	t.Run("returns built-in github preset", func(t *testing.T) {
		preset, ok := ssh.FindPreset("github")
		assert.True(t, ok)
		assert.Equal(t, "GitHub", preset.Name)
		assert.Equal(t, "github.com", preset.HostName)
		assert.Equal(t, "git", preset.User)
	})

	t.Run("returns built-in gitlab preset", func(t *testing.T) {
		preset, ok := ssh.FindPreset("gitlab")
		assert.True(t, ok)
		assert.Equal(t, "GitLab", preset.Name)
		assert.Equal(t, "gitlab.com", preset.HostName)
		assert.Equal(t, "git", preset.User)
	})

	t.Run("returns false for unknown preset", func(t *testing.T) {
		_, ok := ssh.FindPreset("unknown-service")
		assert.False(t, ok)
	})
}
