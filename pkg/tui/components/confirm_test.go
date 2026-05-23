package components

import (
	"testing"
)

func TestConfirmComponent(t *testing.T) {
	c := NewConfirm("Test Confirm", "Are you sure?")

	if c.Title != "Test Confirm" {
		t.Errorf("expected Title 'Test Confirm', got %q", c.Title)
	}

	if c.YesSelected {
		t.Error("expected YesSelected to be false by default")
	}

	// Move left
	done, confirmed := c.Update("left")
	if done || confirmed {
		t.Error("expected navigation to not finalize selection")
	}
	if !c.YesSelected {
		t.Error("expected YesSelected to be true after left key press")
	}

	// Move right
	done, confirmed = c.Update("right")
	if done || confirmed {
		t.Error("expected navigation to not finalize selection")
	}
	if c.YesSelected {
		t.Error("expected YesSelected to be false after right key press")
	}

	// Confirm 'No'
	done, confirmed = c.Update("enter")
	if !done {
		t.Error("expected done to be true after enter key press")
	}
	if confirmed {
		t.Error("expected confirmed to be false since 'No' was focused")
	}

	// Esc test
	done, confirmed = c.Update("esc")
	if !done {
		t.Error("expected done to be true after esc key press")
	}
	if confirmed {
		t.Error("expected confirmed to be false on cancel/esc")
	}
}
