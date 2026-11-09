package theme

import "github.com/charmbracelet/lipgloss"

// Theme bundles the semantic color roles used for consistent component styling.
type Theme struct {
	Primary   lipgloss.TerminalColor
	Secondary lipgloss.TerminalColor
	Muted     lipgloss.TerminalColor
	Success   lipgloss.TerminalColor
	Error     lipgloss.TerminalColor
	Info      lipgloss.TerminalColor
	Text      lipgloss.TerminalColor
	OnAccent  lipgloss.TerminalColor
	Highlight lipgloss.TerminalColor
}
