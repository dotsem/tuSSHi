package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Form wraps a huh.Form and its submission callback.
type Form struct {
	Form     *huh.Form
	OnSubmit func()
}

// NewForm creates a new form component.
func NewForm(f *huh.Form, onSubmit func()) *Form {
	return &Form{
		Form:     f,
		OnSubmit: onSubmit,
	}
}

// Init initializes the huh form.
func (f *Form) Init() tea.Cmd {
	return f.Form.Init()
}

// Update delegates key inputs to Huh and triggers submission.
func (f *Form) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return nil, true
	}

	newForm, cmd := f.Form.Update(msg)
	f.Form = newForm.(*huh.Form)

	switch f.Form.State {
	case huh.StateCompleted:
		if f.OnSubmit != nil {
			f.OnSubmit()
		}
		return nil, true
	case huh.StateAborted:
		return nil, true
	}

	return cmd, false
}

// View renders the huh form.
func (f *Form) View(width int) string {
	return f.Form.View()
}
