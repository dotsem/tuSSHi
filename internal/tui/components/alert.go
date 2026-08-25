package components

import (
	"strings"
	"tusshi/internal/constants"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Alert represents a reusable TUI alert overlay component.
type Alert struct {
	Title   string
	Message string
	IsError bool
	Theme   theme.Theme
}

// Init initializes the alert dialog.
func (a *Alert) Init() tea.Cmd {
	return nil
}

// Update processes navigation and dismiss events.
func (a *Alert) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEsc),
			utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyQuit),
			utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEnter):
			return nil, true
		}
	}
	return nil, false
}

// View renders the alert box styled with Lip Gloss.
func (a *Alert) View(width int) string {
	accentColor := a.Theme.Primary
	if a.IsError {
		accentColor = a.Theme.Error
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Center).
		Width(width)

	divider := lipgloss.NewStyle().Foreground(a.Theme.Muted).Render(strings.Repeat("─", width))

	rows := []string{
		titleStyle.Render(a.Title),
		divider,
		"",
		msgStyle.Render(a.Message),
		"",
		lipgloss.NewStyle().Foreground(a.Theme.Muted).Align(lipgloss.Center).Width(width).Render("Press OK / Enter / Esc / q to dismiss"),
	}

	return strings.Join(rows, "\n")
}
