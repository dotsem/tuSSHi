// Package theme defines the visual color configurations used across all components.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Mock provides a blank/zero color palette used to keep component unit test assertions stable and color-independent.
var Mock = Theme{
	Primary:   lipgloss.Color("0"),
	Secondary: lipgloss.Color("0"),
	Muted:     lipgloss.Color("0"),
	Bg:        lipgloss.Color("0"),
	Success:   lipgloss.Color("0"),
	Error:     lipgloss.Color("0"),
}
