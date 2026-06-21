package components_test

import (
	"testing"

	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestHelpComponent(t *testing.T) {
	opts := []components.HelpOption{
		{Shortcut: "j/k", Description: "Navigate list"},
		{Shortcut: "a", Description: "Add host"},
	}

	h := &components.Help{
		Options: opts,
		Theme:   theme.Mock,
	}

	t.Run("init", func(t *testing.T) {
		assert.Nil(t, h.Init())
	})

	t.Run("view", func(t *testing.T) {
		rendered := h.View(50)
		assert.Contains(t, rendered, "Available Commands")
		assert.Contains(t, rendered, "j/k")
		assert.Contains(t, rendered, "Navigate list")
		assert.Contains(t, rendered, "a")
		assert.Contains(t, rendered, "Add host")
	})

	t.Run("closing keys", func(t *testing.T) {
		_, done := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
		assert.True(t, done)

		_, done = h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
		assert.True(t, done)

		_, done = h.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, done)
	})

	t.Run("non-closing keys", func(t *testing.T) {
		_, done := h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		assert.False(t, done)
	})
}
