package components

import (
	"fmt"
	"strings"

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
	Options      []HelpOption
	ColorPrimary lipgloss.TerminalColor
	ColorMuted   lipgloss.TerminalColor
}

// NewHelp creates a new help dialog component.
func NewHelp(options []HelpOption, primary, muted lipgloss.TerminalColor) *Help {
	return &Help{
		Options:      options,
		ColorPrimary: primary,
		ColorMuted:   muted,
	}
}

// Init initializes the help dialog.
func (h *Help) Init() tea.Cmd {
	return nil
}

// Update handles closing the help dialog.
func (h *Help) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc", "q", "enter":
			return nil, true
		}
	}
	return nil, false
}

// View renders the help dialog content.
func (h *Help) View(width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(h.ColorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	header := titleStyle.Render("Available Commands")
	divider := lipgloss.NewStyle().Foreground(h.ColorMuted).Render(strings.Repeat("─", width))

	cmdStyle := lipgloss.NewStyle().Foreground(h.ColorPrimary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	var rows []string
	rows = append(rows, header, divider)
	for _, opt := range h.Options {
		rows = append(rows, fmt.Sprintf("  %-22s %s", cmdStyle.Render(opt.Shortcut), descStyle.Render(opt.Description)))
	}
	rows = append(rows, "", lipgloss.NewStyle().Foreground(h.ColorMuted).Align(lipgloss.Center).Width(width).Render("Press Esc to close"))

	return strings.Join(rows, "\n")
}
