package components

import (
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Form wraps a huh.Form and its submission callback.
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
		switch keyMsg.String() {
		case keyEsc:
			return nil, true

		case "alt+enter", "ctrl+s":
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

// View renders the huh form with a custom help footer.
func (f *Form) View(_ int) string {
	formView := f.Form.View()
	if f.Form.State != huh.StateNormal {
		return formView
	}

	var bindings []key.Binding

	bindings = append(bindings, key.NewBinding(
		key.WithKeys("ctrl+s", "alt+enter"),
		key.WithHelp("ctrl+s/alt+enter", "save"),
	))

	if focused := f.Form.GetFocusedField(); focused != nil {
		focusedType := reflect.TypeOf(focused).String()
		if strings.Contains(focusedType, "Select") {
			bindings = append(bindings, key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			))
		} else {
			bindings = append(bindings, key.NewBinding(
				key.WithKeys("enter", "tab"),
				key.WithHelp("enter", "next"),
			))
		}
	}

	bindings = append(bindings, key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "back"),
	))

	bindings = append(bindings, key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit"),
	))

	helpView := f.Form.Help().ShortHelpView(bindings)

	return formView + "\n\n" + helpView
}
