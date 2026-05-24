package components_test

import (
	"strings"
	"testing"

	"tusshi/pkg/tui/components"
	"tusshi/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
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

	// Test Init
	if cmd := h.Init(); cmd != nil {
		t.Error("expected Init to return nil")
	}

	// Test View rendering
	rendered := h.View(50)
	if !strings.Contains(rendered, "Available Commands") {
		t.Error("expected view to contain header")
	}
	if !strings.Contains(rendered, "j/k") || !strings.Contains(rendered, "Navigate list") {
		t.Error("expected view to contain 'j/k' shortcut and description")
	}
	if !strings.Contains(rendered, "a") || !strings.Contains(rendered, "Add host") {
		t.Error("expected view to contain 'a' shortcut and description")
	}

	// Test closing keys: Esc
	_, done := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !done {
		t.Error("expected Esc to finish help dialog")
	}

	// Test closing keys: q
	_, done = h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if !done {
		t.Error("expected 'q' to finish help dialog")
	}

	// Test closing keys: Enter
	_, done = h.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !done {
		t.Error("expected Enter to finish help dialog")
	}

	// Test non-closing keys
	_, done = h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if done {
		t.Error("expected other keys to not finish help dialog")
	}
}
