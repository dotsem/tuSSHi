package tui_test

import (
	"testing"

	tussh "tusshi/internal/ssh"
	"tusshi/internal/tui"

	"github.com/stretchr/testify/assert"
)

func TestServiceFormState(t *testing.T) {
	t.Run("ApplyPreset populates github defaults when fields empty", func(t *testing.T) {
		state := &tui.ServiceFormState{
			Action:      "add",
			PresetAlias: "github",
			KeySource:   "generate",
			KeyType:     tussh.KeyTypeED25519,
		}

		state.ApplyPreset()

		assert.Equal(t, "github", state.HostAlias)
		assert.Equal(t, "github.com", state.HostName)
		assert.Equal(t, "git", state.HostUser)
		assert.Contains(t, state.KeyPath, "id_ed25519_github")
	})

	t.Run("ApplyPreset populates gitlab defaults when fields empty", func(t *testing.T) {
		state := &tui.ServiceFormState{
			Action:      "add",
			PresetAlias: "gitlab",
			KeySource:   "generate",
			KeyType:     tussh.KeyTypeED25519,
		}

		state.ApplyPreset()

		assert.Equal(t, "gitlab", state.HostAlias)
		assert.Equal(t, "gitlab.com", state.HostName)
		assert.Equal(t, "git", state.HostUser)
		assert.Contains(t, state.KeyPath, "id_ed25519_gitlab")
	})

	t.Run("ProvideDefaultKeyPath generates non-colliding key path", func(t *testing.T) {
		state := &tui.ServiceFormState{
			HostAlias: "github-work",
			KeyType:   tussh.KeyTypeED25519,
		}

		path := state.ProvideDefaultKeyPath()
		assert.Contains(t, path, "id_ed25519_github-work")
	})

	t.Run("BuildServiceForm returns valid non-nil form", func(t *testing.T) {
		state := &tui.ServiceFormState{
			Action:      "add",
			PresetAlias: "github",
		}

		form := tui.BuildServiceForm(state)
		assert.NotNil(t, form)
	})
}
