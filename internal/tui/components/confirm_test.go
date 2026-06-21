package components_test

import (
	"testing"

	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestConfirmComponent(t *testing.T) {
	confirmedCalled := false
	c := &components.Confirm{
		Title:   "Test Confirm",
		Message: "Are you sure?",
		Theme:   theme.Mock,
		OnConfirm: func() tea.Cmd {
			confirmedCalled = true
			return nil
		},
	}

	t.Run("initial state", func(t *testing.T) {
		assert.Equal(t, "Test Confirm", c.Title)
		assert.False(t, c.YesSelected)
	})

	t.Run("move left", func(t *testing.T) {
		_, done := c.Update(tea.KeyMsg{Type: tea.KeyLeft})
		assert.False(t, done)
		assert.True(t, c.YesSelected)
	})

	t.Run("move right", func(t *testing.T) {
		_, done := c.Update(tea.KeyMsg{Type: tea.KeyRight})
		assert.False(t, done)
		assert.False(t, c.YesSelected)
	})

	t.Run("confirm no", func(t *testing.T) {
		_, done := c.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, done)
		assert.False(t, confirmedCalled)
	})

	t.Run("select yes and confirm yes", func(t *testing.T) {
		c.YesSelected = true
		_, done := c.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, done)
		assert.True(t, confirmedCalled)
	})

	t.Run("esc close", func(t *testing.T) {
		_, done := c.Update(tea.KeyMsg{Type: tea.KeyEsc})
		assert.True(t, done)
	})
}

func TestConfirmCustomLabels(t *testing.T) {
	c := &components.Confirm{
		Theme: theme.Mock,
	}

	t.Run("default labels", func(t *testing.T) {
		viewEmpty := c.View(40)
		assert.Contains(t, viewEmpty, " Yes ")
		assert.Contains(t, viewEmpty, " No  ")
	})

	t.Run("custom labels", func(t *testing.T) {
		c.YesStr = " Delete "
		c.NoStr = " Cancel "

		viewCustom := c.View(40)
		assert.Contains(t, viewCustom, " Delete ")
		assert.Contains(t, viewCustom, " Cancel ")
	})
}
