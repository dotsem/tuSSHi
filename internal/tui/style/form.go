package style

import (
	"tusshi/internal/tui/theme"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Form builds the huh form theme from the global terminal palette.
func Form() *huh.Theme {
	t := huh.ThemeBase()
	applyFocusedForm(t, theme.Global)
	applyBlurredForm(t)
	return t
}

func applyFocusedForm(t *huh.Theme, th theme.Theme) {
	f := &t.Focused
	f.Base = f.Base.BorderForeground(th.Muted)
	f.Card = f.Base
	f.Title = f.Title.Foreground(th.Primary).Bold(true)
	f.NoteTitle = f.NoteTitle.Foreground(th.Primary).Bold(true).MarginBottom(1)
	f.Directory = f.Directory.Foreground(th.Primary)
	f.Description = f.Description.Foreground(th.Muted)
	f.ErrorIndicator = f.ErrorIndicator.Foreground(th.Error)
	f.ErrorMessage = f.ErrorMessage.Foreground(th.Error)
	f.SelectSelector = f.SelectSelector.Foreground(th.Primary)
	f.NextIndicator = f.NextIndicator.Foreground(th.Primary)
	f.PrevIndicator = f.PrevIndicator.Foreground(th.Primary)
	f.Option = f.Option.Foreground(th.Text)
	f.MultiSelectSelector = f.MultiSelectSelector.Foreground(th.Primary)
	f.SelectedOption = f.SelectedOption.Foreground(th.Success)
	f.SelectedPrefix = lipgloss.NewStyle().Foreground(th.Success).SetString("✓ ")
	f.UnselectedPrefix = lipgloss.NewStyle().Foreground(th.Muted).SetString("• ")
	f.UnselectedOption = f.UnselectedOption.Foreground(th.Text)
	f.FocusedButton = f.FocusedButton.Foreground(th.OnAccent).Background(th.Primary).Bold(true)
	f.Next = f.FocusedButton
	f.BlurredButton = f.BlurredButton.Foreground(th.Text).Background(th.Highlight)
	f.TextInput.Cursor = f.TextInput.Cursor.Foreground(th.Primary)
	f.TextInput.Placeholder = f.TextInput.Placeholder.Foreground(th.Muted)
	f.TextInput.Prompt = f.TextInput.Prompt.Foreground(th.Primary)
}

func applyBlurredForm(t *huh.Theme) {
	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
}
