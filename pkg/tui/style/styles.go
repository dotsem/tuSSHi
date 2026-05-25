// Package style defines visual and typography styles for rendering the TUSSHI TUI.
package style

import (
	"tusshi/pkg/tui/theme"

	"github.com/charmbracelet/lipgloss"
)

// Structural and typography styles for the user interface, bound to theme colors.
var (
	Title = lipgloss.NewStyle().
		Foreground(theme.Global.Primary).
		Bold(true).
		Padding(0, 1)

	TabActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(theme.Global.Primary).
			Bold(true).
			Padding(0, 2)

	TabInactive = lipgloss.NewStyle().
			Foreground(theme.Global.Muted).
			Padding(0, 2)

	Header = lipgloss.NewStyle().
		Padding(0, 1)

	TableHeader = lipgloss.NewStyle().
			Foreground(theme.Global.Muted).
			Bold(true)

	RowActive = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	RowInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	Cell = lipgloss.NewStyle().
		Padding(0, 1)

	Footer = lipgloss.NewStyle().
		Foreground(theme.Global.Muted)

	Alert = lipgloss.NewStyle().
		Foreground(theme.Global.Success).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(theme.Global.Error).
		Bold(true)

	NormalPrompt = lipgloss.NewStyle().
			Foreground(theme.Global.Secondary).
			Bold(true)

	CommandPrompt = lipgloss.NewStyle().
			Foreground(theme.Global.Primary).
			Bold(true)

	HeaderBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Global.Primary).
			Padding(0, 1)

	BodyBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Global.Primary).
		Padding(0, 1)

	FooterBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Global.Primary).
			Padding(0, 1)

	Dialog = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Global.Primary).
		Background(lipgloss.Color("#1A1A1A")).
		Padding(1, 2)

	Muted = lipgloss.NewStyle().
		Foreground(theme.Global.Muted)

	StatusOnline = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	StatusOffline = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)
