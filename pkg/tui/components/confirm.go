package components

import (
	"strings"
	"tusshi/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Confirm represents a reusable TUI confirmation modal component.
type Confirm struct {
	Title       string
	Message     string
	YesSelected bool
	Theme       theme.Theme
	OnConfirm   func() tea.Cmd
	YesStr      string
	NoStr       string
	Destructive bool
}

// NewConfirm creates a new confirmation component.
func NewConfirm(title, message string, theme theme.Theme, onConfirm func() tea.Cmd) *Confirm {
	return &Confirm{
		Title:       title,
		Message:     message,
		YesSelected: false,
		Theme:       theme,
		OnConfirm:   onConfirm,
		YesStr:      " Yes ",
		NoStr:       " No  ",
		Destructive: false,
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
		Foreground(c.Theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Align(lipgloss.Center).
		Width(width)

	divider := lipgloss.NewStyle().Foreground(c.Theme.Muted).Render(strings.Repeat("─", width))

	btnActive := lipgloss.NewStyle().
		Background(c.Theme.Primary).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 3)

	btnDestructive := lipgloss.NewStyle().
		Background(c.Theme.Error).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 3)

	btnInactive := lipgloss.NewStyle().
		Background(lipgloss.Color("237")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 3)

	var yesBtn, noBtn string
	if c.YesSelected {
		if c.Destructive {
			yesBtn = btnDestructive.Render(c.YesStr)
		} else {
			yesBtn = btnActive.Render(c.YesStr)
		}
		noBtn = btnInactive.Render(c.NoStr)
	} else {
		if c.Destructive {
			noBtn = btnDestructive.Render(c.NoStr)
		} else {
			noBtn = btnActive.Render(c.NoStr)
		}
		yesBtn = btnInactive.Render(c.YesStr)
		noBtn = btnActive.Render(c.NoStr)
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
		lipgloss.NewStyle().Foreground(c.Theme.Muted).Align(lipgloss.Center).Width(width).Render("← / → to switch • Enter to select • Esc to cancel"),
	}

	return strings.Join(rows, "\n")
}
