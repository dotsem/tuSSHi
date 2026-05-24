package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Global bundles the system theme colors for UI components.
var Global = Theme{
	Primary:   Primary,
	Secondary: Secondary,
	Muted:     Muted,
	Bg:        Bg,
	Success:   Success,
	Error:     Error,
}

const (
	Primary   = lipgloss.Color("#FF5500") // TODO: follow system color theme (GTK4/QT6?)
	Secondary = lipgloss.Color("#1F1F1F")
	Muted     = lipgloss.Color("#757575")
	Bg        = lipgloss.Color("#121212")
	Success   = lipgloss.Color("#FF7851")
	Error     = lipgloss.Color("#ff5050")
)
