package components

import (
	"tusshi/internal/constants"
	"tusshi/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Form is an interactive form component wrapping Huh form logic.
type Form struct {
	Form     *huh.Form
	OnSubmit func()
	Validate func() error
}

// Init initializes the huh form.
func (f *Form) Init() tea.Cmd {
	return f.Form.Init()
}

// Update delegates key inputs to Huh and triggers submission.
func (f *Form) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEsc):
			return nil, true

		case utils.MatchesMultipleStringOptions(keyMsg.String(), "alt+enter, ctrl+s"):
			if focused := f.Form.GetFocusedField(); focused != nil {
				_ = focused.Blur()
			}

			if f.Validate != nil {
				if err := f.Validate(); err != nil {
					if focused := f.Form.GetFocusedField(); focused != nil {
						_ = focused.Focus()
					}
					return nil, false
				}
			}

			f.Form.State = huh.StateCompleted
			if f.OnSubmit != nil {
				f.OnSubmit()
			}
			return nil, true
		}
	}

	_, formCmd := f.Form.Update(msg)
	if f.Form.State == huh.StateCompleted || f.Form.State == huh.StateAborted {
		if f.Form.State == huh.StateCompleted {
			if f.Validate != nil {
				if err := f.Validate(); err != nil {
					return formCmd, false
				}
			}
			if f.OnSubmit != nil {
				f.OnSubmit()
			}
		}
		return formCmd, true
	}

	return formCmd, false

}

// View renders the interactive Huh form.
func (f *Form) View(_ int) string {
	return f.Form.View()
}
