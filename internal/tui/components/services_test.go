package components_test

import (
	"testing"

	"tusshi/internal/config"
	gossh "tusshi/internal/ssh"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestServicesComponent(t *testing.T) {
	sampleHosts := []*config.Host{
		{Alias: "github", Name: "github.com", User: "git", IdentityFile: "~/.ssh/id_ed25519_github", IsService: true},
		{Alias: "gitlab", Name: "gitlab.com", User: "git", IdentityFile: "~/.ssh/id_ed25519_gitlab", IsService: true},
	}

	t.Run("renders table rows and handles down/up navigation", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		assert.Equal(t, 0, svc.SelectedIndex)

		// Press down arrow
		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyDown})
		assert.Nil(t, cmd)
		assert.False(t, done)
		assert.Equal(t, 1, svc.SelectedIndex)

		// Press up arrow
		cmd, done = svc.Update(tea.KeyMsg{Type: tea.KeyUp})
		assert.Nil(t, cmd)
		assert.False(t, done)
		assert.Equal(t, 0, svc.SelectedIndex)
	})

	t.Run("dispatches edit action on 'e' key", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
		assert.True(t, done)
		assert.NotNil(t, cmd)

		msg := cmd()
		actionMsg, ok := msg.(components.ServiceActionMsg)
		assert.True(t, ok)
		assert.Equal(t, "edit", actionMsg.Action)
		assert.Equal(t, "github", actionMsg.Host.Alias)
	})

	t.Run("dispatches add action on 'a' key", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
		assert.True(t, done)
		assert.NotNil(t, cmd)

		msg := cmd()
		actionMsg, ok := msg.(components.ServiceActionMsg)
		assert.True(t, ok)
		assert.Equal(t, "add", actionMsg.Action)
	})

	t.Run("dispatches delete action on 'd' key", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
		assert.True(t, done)
		assert.NotNil(t, cmd)

		msg := cmd()
		actionMsg, ok := msg.(components.ServiceActionMsg)
		assert.True(t, ok)
		assert.Equal(t, "delete", actionMsg.Action)
		assert.Equal(t, "github", actionMsg.Host.Alias)
	})

	t.Run("stores auth check result and renders status", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		svc.SetResult("github", gossh.AuthResult{OK: true})
		svc.SetResult("gitlab", gossh.AuthResult{OK: false, Error: "Permission denied"})

		viewStr := svc.View(80)
		assert.Contains(t, viewStr, "github")
		assert.Contains(t, viewStr, "gitlab")
		assert.Contains(t, viewStr, "Permission denied")
	})

	t.Run("dispatches ping action on 'p' key", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
		assert.False(t, done)
		assert.NotNil(t, cmd)

		msg := cmd()
		actionMsg, ok := msg.(components.ServiceActionMsg)
		assert.True(t, ok)
		assert.Equal(t, "ping", actionMsg.Action)
		assert.Equal(t, "github", actionMsg.Host.Alias)
	})

	t.Run("dispatches pingall action on 'P' key", func(t *testing.T) {
		svc := &components.Services{
			Hosts: sampleHosts,
			Theme: theme.Global,
		}

		cmd, done := svc.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("P")})
		assert.False(t, done)
		assert.NotNil(t, cmd)

		msg := cmd()
		actionMsg, ok := msg.(components.ServiceActionMsg)
		assert.True(t, ok)
		assert.Equal(t, "pingall", actionMsg.Action)
	})
}
