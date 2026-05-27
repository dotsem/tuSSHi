package components

import (
	"fmt"
	"strings"
	"tusshi/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpOption represents a single command shortcut and its description.
type HelpOption struct {
	Shortcut    string
	Description string
}

// Help represents the interactive help dialog component.
type Help struct {
	Options []HelpOption
	Theme   theme.Theme
}

// Init initializes the help dialog.
func (h *Help) Init() tea.Cmd {
	return nil
}

// Update handles closing the help dialog.
func (h *Help) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case keyEsc, "q", keyEnter:
			return nil, true
		}
	}
	return nil, false
}

// View renders the help dialog content.
func (h *Help) View(width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(h.Theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	header := titleStyle.Render("Available Commands")
	divider := lipgloss.NewStyle().Foreground(h.Theme.Muted).Render(strings.Repeat("─", width))

	cmdStyle := lipgloss.NewStyle().Foreground(h.Theme.Primary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	var rows []string
	rows = append(rows, header, divider)
	for _, opt := range h.Options {
		rows = append(rows, fmt.Sprintf("  %-22s %s", cmdStyle.Render(opt.Shortcut), descStyle.Render(opt.Description)))
	}
	rows = append(rows, "", lipgloss.NewStyle().Foreground(h.Theme.Muted).Align(lipgloss.Center).Width(width).Render("Press Esc to close"))

	return strings.Join(rows, "\n")
}
