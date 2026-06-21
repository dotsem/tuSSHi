package components_test

import (
	"errors"
	"testing"

	"tusshi/internal/tui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
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

	t.Run("init", func(t *testing.T) {
		cmd := f.Init()
		assert.NotNil(t, cmd)
	})

	t.Run("view rendering", func(t *testing.T) {
		rendered := f.View(40)
		assert.NotEmpty(t, rendered)
	})

	t.Run("cancellation", func(t *testing.T) {
		_, done := f.Update(tea.KeyMsg{Type: tea.KeyEsc})
		assert.True(t, done)
	})

	t.Run("completed state transition", func(t *testing.T) {
		f.Form.State = huh.StateCompleted
		_, done := f.Update(nil)
		assert.True(t, done)
		assert.True(t, submitted)
	})

	t.Run("aborted state transition", func(t *testing.T) {
		f.Form.State = huh.StateAborted
		_, done := f.Update(nil)
		assert.True(t, done)
	})
}

func TestFormValidationAndSubmission(t *testing.T) {
	var val string
	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Value(&val),
		),
	)

	submitted := false
	validated := false
	var validationErr error

	f := &components.Form{
		Form: huhForm,
		OnSubmit: func() {
			submitted = true
		},
		Validate: func() error {
			validated = true
			return validationErr
		},
	}

	t.Run("validation fails", func(t *testing.T) {
		f.Form.State = huh.StateNormal
		validationErr = errors.New("invalid field")
		_, done := f.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
		assert.False(t, done)
		assert.True(t, validated)
		assert.False(t, submitted)
	})

	t.Run("validation succeeds", func(t *testing.T) {
		validated = false
		submitted = false
		validationErr = nil

		_, done := f.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
		assert.True(t, done)
		assert.True(t, validated)
		assert.True(t, submitted)
		assert.Equal(t, huh.StateCompleted, f.Form.State)
	})
}
