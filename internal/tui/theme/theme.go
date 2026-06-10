package theme

import "github.com/charmbracelet/lipgloss"

// Theme bundles the color configuration for consistent component styling.
type Theme struct {
	Primary   lipgloss.TerminalColor
	Secondary lipgloss.TerminalColor
	Muted     lipgloss.TerminalColor
	Bg        lipgloss.TerminalColor
	Success   lipgloss.TerminalColor
	Error     lipgloss.TerminalColor
}
