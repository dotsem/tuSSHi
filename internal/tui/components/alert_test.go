package components_test

import (
	"strings"
	"testing"

	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAlertComponent(t *testing.T) {
	alert := &components.Alert{
		Title:   "Backup Failed",
		Message: "Could not create initial pre-tuSSHi backup",
		IsError: true,
		Theme:   theme.Mock,
	}

	if cmd := alert.Init(); cmd != nil {
		t.Errorf("expected Init to return nil, got %v", cmd)
	}

	rendered := alert.View(50)
	if !strings.Contains(rendered, "Backup Failed") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(rendered, "Could not create initial pre-tuSSHi backup") {
		t.Error("expected view to contain message")
	}
	if !strings.Contains(rendered, "OK") {
		t.Error("expected view to contain OK button")
	}

	_, done := alert.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !done {
		t.Error("expected Esc to dismiss alert")
	}
	_, done = alert.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if !done {
		t.Error("expected 'q' to dismiss alert")
	}
	_, done = alert.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !done {
		t.Error("expected Enter to dismiss alert")
	}
	_, done = alert.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if done {
		t.Error("expected other keys to not dismiss alert")
	}
}
