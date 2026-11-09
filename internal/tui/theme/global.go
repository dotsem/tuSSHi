// Package theme defines the visual color configurations used across all components.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Why: colors reference slots of the terminal's own ANSI palette instead of fixed
// RGB values, so the user's terminal theme decides the final appearance on both
// dark and light backgrounds.
const (
	ansiBlack       = lipgloss.ANSIColor(0)
	ansiRed         = lipgloss.ANSIColor(1)
	ansiGreen       = lipgloss.ANSIColor(2)
	ansiYellow      = lipgloss.ANSIColor(3)
	ansiBlue        = lipgloss.ANSIColor(4)
	ansiBrightBlack = lipgloss.ANSIColor(8)
)

// Global bundles the system theme colors for UI components.
var Global = Theme{
	Primary:   ansiRed,
	Secondary: ansiYellow,
	Muted:     ansiBrightBlack,
	Success:   ansiGreen,
	Error:     ansiRed,
	Info:      ansiBlue,
	Text:      lipgloss.NoColor{},
	OnAccent:  ansiBlack,
	Highlight: ansiBrightBlack,
}
