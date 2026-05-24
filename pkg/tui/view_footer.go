package tui

import (
	"tusshi/pkg/tui/style"

	"github.com/charmbracelet/lipgloss"
)

// renderFooter generates alert banners and context-aware keyboard mappings.
func (m *Model) renderFooter(width int) string {
	var alertBar string
	if m.ErrorText != "" {
		alertBar = style.Error.Render("ERROR: " + m.ErrorText)
	} else if m.AlertText != "" {
		alertBar = style.Alert.Render("SUCCESS: " + m.AlertText)
	}

	var cmdBar string
	switch m.Mode {
	case ModeCommand:
		cmdBar = m.CommandInput.View()
	case ModeSearch:
		cmdBar = style.NormalPrompt.Render("[Search Mode] ") + getSearchShortcuts(width)
	default:
		cmdBar = style.Footer.Render(getShortcuts(width))
	}

	if alertBar != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			alertBar,
			cmdBar,
		)
	}

	return cmdBar
}

func getShortcuts(width int) string {
	full := "[Normal] j/k: Nav • h/l: Tabs • a/e: Add/Edit • d: Del • Enter: Connect • : Cmd"
	medium := "j/k: Nav • h/l: Tabs • a/e: Add/Edit • d: Del • Enter: Connect • : Cmd"
	short := "j/k: Nav • h/l: Tabs • a/e: Edit • Enter: Connect"
	minimal := "j/k: Nav • Enter: Connect"

	if width >= len(full) {
		return full
	}
	if width >= len(medium) {
		return medium
	}
	if width >= len(short) {
		return short
	}
	if width >= len(minimal) {
		return minimal
	}
	return "j/k/Enter"
}

func getSearchShortcuts(width int) string {
	full := "[Search] Type to filter. Esc/Enter: Done"
	short := "Type to filter • Esc: Exit"
	if width >= len(full) {
		return full
	}
	if width >= len(short) {
		return short
	}
	return "Search..."
}
