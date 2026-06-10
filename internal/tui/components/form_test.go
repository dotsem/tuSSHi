package components_test

import (
	"testing"

	"tusshi/internal/tui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func TestFormComponent(t *testing.T) {
	var val string
	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Value(&val),
		),
	)

	submitted := false
	f := &components.Form{
		Form: huhForm,
		OnSubmit: func() {
			submitted = true
		},
	}

	// Test Init
	cmd := f.Init()
	if cmd == nil {
		t.Error("expected Init to return a non-nil command for huh form initialization")
	}

	// Test View rendering
	rendered := f.View(40)
	if rendered == "" {
		t.Error("expected View to return a non-empty string representation of the form")
	}

	// Test cancellation key: Esc
	_, done := f.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !done {
		t.Error("expected Esc to finish/cancel the form component")
	}

	// Test manual completed state transition to trigger callback
	f.Form.State = huh.StateCompleted
	_, done = f.Update(nil)
	if !done {
		t.Error("expected completed form state to return done = true")
	}
	if !submitted {
		t.Error("expected OnSubmit callback to be executed on form completion")
	}

	// Test manual aborted state transition
	f.Form.State = huh.StateAborted
	_, done = f.Update(nil)
	if !done {
		t.Error("expected aborted form state to return done = true")
	}
}
