package components

import (
	"strings"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Alert represents a reusable, self-contained TUI dialog modal for notices or errors.
type Alert struct {
	Title   string
	Message string
	IsError bool
	Theme   theme.Theme
}

// Init initializes the alert component.
func (a *Alert) Init() tea.Cmd {
	return nil
}

// Update processes navigation and dismiss events.
func (a *Alert) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case keyEsc, "q", keyEnter:
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

	okBtn := lipgloss.NewStyle().
		Background(accentColor).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 3).
		Render(" OK ")

	buttonsStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(width)

	rows := []string{
		titleStyle.Render(a.Title),
		divider,
		"",
		msgStyle.Render(a.Message),
		"",
		buttonsStyle.Render(okBtn),
		"",
		lipgloss.NewStyle().Foreground(a.Theme.Muted).Align(lipgloss.Center).Width(width).Render("Press Enter or Esc to dismiss"),
	}

	return strings.Join(rows, "\n")
}
