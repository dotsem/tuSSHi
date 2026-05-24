package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Confirm represents a reusable TUI confirmation modal component.
type Confirm struct {
	Title       string
	Message     string
	YesSelected bool
	OnConfirm   func() tea.Cmd
}

// NewConfirm creates a new confirmation component.
func NewConfirm(title, message string, onConfirm func() tea.Cmd) *Confirm {
	return &Confirm{
		Title:       title,
		Message:     message,
		YesSelected: false,
		OnConfirm:   onConfirm,
	}
}

// Init initializes the confirmation dialog.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update processes navigation and selection events.
func (c *Confirm) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "left", "h":
			c.YesSelected = true
		case "right", "l":
			c.YesSelected = false
		case "enter":
			if c.YesSelected && c.OnConfirm != nil {
				return c.OnConfirm(), true
			}
			return nil, true
		case "esc", "q":
			return nil, true
		}
	}
	return nil, false
}

// View renders the confirmation modal nicely styled with Lip Gloss.
func (c *Confirm) View(width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Center).
		Width(width)

	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", width))

	btnActive := lipgloss.NewStyle().
		Background(lipgloss.Color("205")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 3)

	btnInactive := lipgloss.NewStyle().
		Background(lipgloss.Color("237")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 3)

	var yesBtn, noBtn string
	if c.YesSelected {
		yesBtn = btnActive.Render(" Yes ")
		noBtn = btnInactive.Render(" No  ")
	} else {
		yesBtn = btnInactive.Render(" Yes ")
		noBtn = btnActive.Render(" No  ")
	}

	buttonsRow := lipgloss.JoinHorizontal(lipgloss.Center,
		yesBtn,
		"     ",
		noBtn,
	)
	buttonsStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(width)

	rows := []string{
		titleStyle.Render(c.Title),
		divider,
		"",
		msgStyle.Render(c.Message),
		"",
		buttonsStyle.Render(buttonsRow),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center).Width(width).Render("← / → to switch • Enter to select • Esc to cancel"),
	}

	return strings.Join(rows, "\n")
}
