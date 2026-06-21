package components_test

import (
	"testing"

	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestAlertComponent(t *testing.T) {
	alert := &components.Alert{
		Title:   "Backup Failed",
		Message: "Could not create initial pre-tuSSHi backup",
		IsError: true,
		Theme:   theme.Mock,
	}

	t.Run("init", func(t *testing.T) {
		assert.Nil(t, alert.Init())
	})

	t.Run("view", func(t *testing.T) {
		rendered := alert.View(50)
		assert.Contains(t, rendered, "Backup Failed")
		assert.Contains(t, rendered, "Could not create initial pre-tuSSHi backup")
		assert.Contains(t, rendered, "OK")
	})

	t.Run("dismiss keys", func(t *testing.T) {
		_, done := alert.Update(tea.KeyMsg{Type: tea.KeyEsc})
		assert.True(t, done)

		_, done = alert.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
		assert.True(t, done)

		_, done = alert.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, done)
	})

	t.Run("non-dismiss keys", func(t *testing.T) {
		_, done := alert.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		assert.False(t, done)
	})
}
