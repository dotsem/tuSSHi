// Package theme defines the visual color configurations used across all components.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Mock provides a blank/zero color palette used to keep component unit test assertions stable and color-independent.
var Mock = Theme{
	Primary:   lipgloss.NoColor{},
	Secondary: lipgloss.NoColor{},
	Muted:     lipgloss.NoColor{},
	Success:   lipgloss.NoColor{},
	Error:     lipgloss.NoColor{},
	Info:      lipgloss.NoColor{},
	Text:      lipgloss.NoColor{},
	OnAccent:  lipgloss.NoColor{},
	Highlight: lipgloss.NoColor{},
}
