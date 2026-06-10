package components_test

import (
	"strings"
	"testing"

	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
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

	if c.Title != "Test Confirm" {
		t.Errorf("expected Title 'Test Confirm', got %q", c.Title)
	}

	if c.YesSelected {
		t.Error("expected YesSelected to be false by default")
	}

	// Move left
	_, done := c.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if done {
		t.Error("expected navigation to not finalize selection")
	}
	if !c.YesSelected {
		t.Error("expected YesSelected to be true after left key press")
	}

	// Move right
	_, done = c.Update(tea.KeyMsg{Type: tea.KeyRight})
	if done {
		t.Error("expected navigation to not finalize selection")
	}
	if c.YesSelected {
		t.Error("expected YesSelected to be false after right key press")
	}

	// Confirm 'No'
	_, done = c.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !done {
		t.Error("expected done to be true after enter key press")
	}
	if confirmedCalled {
		t.Error("expected confirmedCalled to be false since 'No' was focused")
	}

	// Select Yes and Confirm 'Yes'
	c.YesSelected = true
	_, done = c.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !done {
		t.Error("expected done to be true after enter key press on Yes")
	}
	if !confirmedCalled {
		t.Error("expected confirmedCalled to be true since 'Yes' was focused")
	}

	// Esc test
	_, done = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !done {
		t.Error("expected done to be true after esc key press")
	}
}

func TestConfirmCustomLabels(t *testing.T) {
	c := &components.Confirm{
		Theme: theme.Mock,
	}

	// Verify defaults render when fields are left empty
	viewEmpty := c.View(40)
	if !strings.Contains(viewEmpty, " Yes ") || !strings.Contains(viewEmpty, " No  ") {
		t.Error("expected view to render default button labels when empty")
	}

	// Apply Option 2 (direct mutation)
	c.YesStr = " Delete "
	c.NoStr = " Cancel "

	viewCustom := c.View(40)
	if !strings.Contains(viewCustom, " Delete ") || !strings.Contains(viewCustom, " Cancel ") {
		t.Error("expected view to render custom button labels after mutation")
	}
}
