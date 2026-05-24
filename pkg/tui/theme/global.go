// Package theme defines the visual color configurations used across all components.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Global bundles the system theme colors for UI components.
var Global = Theme{
	Primary:   primary,
	Secondary: secondary,
	Muted:     muted,
	Bg:        bg,
	Success:   success,
	Error:     errColor,
}

const (
	primary   = lipgloss.Color("#FF5500")
	secondary = lipgloss.Color("#1F1F1F")
	muted     = lipgloss.Color("#757575")
	bg        = lipgloss.Color("#121212")
	success   = lipgloss.Color("#FF7851")
	errColor  = lipgloss.Color("#ff5050")
)
