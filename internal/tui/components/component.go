// Package components defines modular, self-contained interactive TUI overlay components.
package components

import tea "github.com/charmbracelet/bubbletea"

const (
	keyEsc   = "esc"
	keyEnter = "enter"
)

// Component represents a self-contained interactive UI overlay.
type Component interface {
	// Init initializes the component and returns any setup commands.
	Init() tea.Cmd

	// Update processes a Bubble Tea message and returns an optional command and whether the component is done.
	Update(msg tea.Msg) (tea.Cmd, bool)

	// View renders the component given the available width constraint.
	View(width int) string
}
