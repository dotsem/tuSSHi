package tui

import "github.com/charmbracelet/lipgloss"

// ThemeColors defines the gorgeous HSL-based palette used for visual hierarchy.
var (
	ColorPrimary   = lipgloss.Color("#FF5500") // TODO: follow system color theme (GTK4/QT6?)
	ColorSecondary = lipgloss.Color("#1F1F1F")
	ColorMuted     = lipgloss.Color("#757575")
	ColorBg        = lipgloss.Color("#121212")
	ColorSuccess   = lipgloss.Color("#FF7851")
	ColorError     = lipgloss.Color("#ff5050")
)

// Style definitions for clean visual boundaries and typography.
var (
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	StyleTabActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 2)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Padding(0, 2)

	StyleHeader = lipgloss.NewStyle().
			Padding(0, 1)

	StyleTableHeader = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Bold(true)

	StyleRowActive = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	StyleRowInactive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	StyleCell = lipgloss.NewStyle().
			Padding(0, 1)

	StyleFooter = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleAlert = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleNormalPrompt = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Bold(true)

	StyleCommandPrompt = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true)

	StyleHeaderBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	StyleBodyBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	StyleFooterBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	StyleDialog = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(1, 2)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
