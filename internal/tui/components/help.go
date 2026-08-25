package components

import (
	"fmt"
	"strings"
	"tusshi/internal/constants"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpOption defines a shortcut key and description pair for the help menu.
type HelpOption struct {
	Shortcut    string
	Description string
}

// Help is a TUI overlay component displaying available shortcuts and commands.
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
		switch {
		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEsc),
			utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyQuit),
			utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEnter):
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

	divider := lipgloss.NewStyle().Foreground(h.Theme.Muted).Render(strings.Repeat("─", width))
	keyStyle := lipgloss.NewStyle().Foreground(h.Theme.Secondary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	rows := []string{
		titleStyle.Render("Available Commands"),
		divider,
		"",
	}

	for _, opt := range h.Options {
		row := fmt.Sprintf("  %-25s %s", keyStyle.Render(opt.Shortcut), descStyle.Render(opt.Description))
		rows = append(rows, row)
	}

	rows = append(rows, "")
	rows = append(rows, lipgloss.NewStyle().Foreground(h.Theme.Muted).Align(lipgloss.Center).Width(width).Render("Press Esc / q / Enter to close"))

	return strings.Join(rows, "\n")
}
