package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// View renders the TUSSHI TUI interface based on state and window constraints.
func (m *Model) View() string {
	if m.Width < 25 || m.Height < 10 {
		return "Terminal is too small."
	}

	bgString := m.renderNormalView(m.Width, m.Height)

	var dialogContent string
	var showDialog bool

	if m.Mode == ModeForm && m.ActiveForm != nil {
		dialogContent = m.ActiveForm.View()
		showDialog = true
	} else if m.Mode == ModeHelp {
		dialogContent = m.renderHelpDialog()
		showDialog = true
	}

	if showDialog {
		bgLines := strings.Split(stripANSI(bgString), "\n")

		dialogWidth := min(60, m.Width-4)
		dialogHeight := min(16, m.Height-2)

		dialogBox := StyleDialog.Width(dialogWidth).Height(dialogHeight).Render(dialogContent)
		dialogLines := strings.Split(dialogBox, "\n")

		dialogW := lipgloss.Width(dialogBox)
		dialogH := len(dialogLines)

		startRow := (len(bgLines) - dialogH) / 2
		startCol := (m.Width - dialogW) / 2

		var finalLines []string
		for i, bgLine := range bgLines {
			bgRunes := []rune(bgLine)
			if len(bgRunes) < m.Width {
				bgRunes = append(bgRunes, []rune(strings.Repeat(" ", m.Width-len(bgRunes)))...)
			} else if len(bgRunes) > m.Width {
				bgRunes = bgRunes[:m.Width]
			}

			if i >= startRow && i < startRow+dialogH {
				dialogLineIdx := i - startRow
				leftPart := bgRunes[:startCol]
				rightPart := bgRunes[startCol+dialogW:]

				mutedLeft := StyleMuted.Render(string(leftPart))
				mutedRight := StyleMuted.Render(string(rightPart))
				dialogLine := dialogLines[dialogLineIdx]

				finalLines = append(finalLines, mutedLeft+dialogLine+mutedRight)
			} else {
				finalLines = append(finalLines, StyleMuted.Render(string(bgRunes)))
			}
		}
		bgString = strings.Join(finalLines, "\n")
	}

	return bgString
}

// renderNormalView compiles the header, table grid, and footer inside inner dimensions.
func (m *Model) renderNormalView(width, height int) string {
	headerBoxHeight := 4

	footerBoxHeight := 3
	if m.ErrorText != "" || m.AlertText != "" {
		footerBoxHeight = 4
	}

	bodyBoxHeight := max(height-headerBoxHeight-footerBoxHeight, 2)

	headerContent := m.renderHeader()
	headerBox := StyleHeaderBox.
		Width(width - 2).
		Height(headerBoxHeight - 2).
		Render(headerContent)

	tableContent := m.renderTable(width-4, bodyBoxHeight-2)
	bodyBox := StyleBodyBox.
		Width(width - 2).
		Height(bodyBoxHeight - 2).
		Render(tableContent)

	footerContent := m.renderFooter(width - 4)
	footerBox := StyleFooterBox.
		Width(width - 2).
		Height(footerBoxHeight - 2).
		Render(footerContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerBox,
		bodyBox,
		footerBox,
	)
}

// renderHeader constructs the config file tabs bar and active search filter.
func (m *Model) renderHeader() string {
	var tabs []string
	for _, t := range m.Tabs {
		label := GetTabLabel(t)
		if t == m.ActiveTab {
			tabs = append(tabs, StyleTabActive.Render(label))
		} else {
			tabs = append(tabs, StyleTabInactive.Render(label))
		}
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	title := StyleTitle.Render(" tuSSHi ")
	headerTop := lipgloss.JoinHorizontal(lipgloss.Left, title, "── ", tabsRow)

	var searchBar string
	if m.Mode == ModeSearch {
		searchBar = m.SearchInput.View()
	} else {
		val := m.SearchInput.Value()
		if val != "" {
			searchBar = StyleNormalPrompt.Render("/ ") + val
		} else {
			searchBar = StyleFooter.Render("/ type to search...")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		headerTop,
		searchBar,
	)
}

// renderTable draws the formatted column grid displaying active connections.
func (m *Model) renderTable(width, maxHeight int) string {
	if len(m.Filtered) == 0 {
		return "\n  No connections found. Press 'a' to add a connection, or ':' for help."
	}

	var headerRow, dividerRow string
	var wAlias, wName, wUser, wPort, wConfig int

	// adaptive column allocation to prevent overflow at small terminal widths
	switch {
	case width >= 61:
		wTotal := max(width-12, 10)
		wAlias = int(float64(wTotal) * 0.20)
		wName = int(float64(wTotal) * 0.30)
		wUser = int(float64(wTotal) * 0.15)
		wPort = int(float64(wTotal) * 0.10)
		wConfig = wTotal - wAlias - wName - wUser - wPort

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
			wPort, "PORT",
			wConfig, "CONFIG",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
			strings.Repeat("─", wPort),
			strings.Repeat("─", wConfig),
		)
	case width >= 41:
		wTotal := max(width-8, 10)
		wAlias = int(float64(wTotal) * 0.30)
		wUser = int(float64(wTotal) * 0.20)
		wName = wTotal - wAlias - wUser

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
		)
	default:
		wTotal := width - 6
		if wTotal < 10 {
			wTotal = 10
		}
		wAlias = int(float64(wTotal) * 0.40)
		wName = wTotal - wAlias

		headerRow = fmt.Sprintf("  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
		)
		dividerRow = fmt.Sprintf("  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
		)
	}

	renderedHeader := StyleTableHeader.Render(headerRow)
	renderedDivider := StyleTableHeader.Render(dividerRow)

	var rows []string
	rows = append(rows, renderedHeader, renderedDivider)

	displayLimit := maxHeight - 2
	if displayLimit <= 0 {
		return ""
	}

	startIndex := 0
	if m.SelectedIndex >= displayLimit {
		startIndex = m.SelectedIndex - displayLimit + 1
	}

	for idx := startIndex; idx < len(m.Filtered) && len(rows) < maxHeight; idx++ {
		h := m.Filtered[idx]
		alias := truncate(h.Alias, wAlias)
		name := truncate(h.Name, wName)

		var rowLine string
		switch {
		case width >= 61:
			user := truncate(h.User, wUser)
			port := h.Port
			if port == "" {
				port = "22"
			}
			cfgNickname := strings.TrimSuffix(GetTabLabel(h.SourceFile), ".conf")
			cfgNickname = strings.TrimSuffix(cfgNickname, "config")
			cfgNickname = truncate(cfgNickname, wConfig)

			rowLine = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
				wAlias, alias,
				wName, name,
				wUser, user,
				wPort, port,
				wConfig, cfgNickname,
			)
		case width >= 41:
			user := truncate(h.User, wUser)
			rowLine = fmt.Sprintf("  %-*s  %-*s  %-*s",
				wAlias, alias,
				wName, name,
				wUser, user,
			)
		default:
			rowLine = fmt.Sprintf("  %-*s  %-*s",
				wAlias, alias,
				wName, name,
			)
		}

		if idx == m.SelectedIndex {
			rowLine = "❯ " + rowLine[2:]
			rows = append(rows, StyleRowActive.Render(rowLine))
		} else {
			rowLine = "  " + rowLine[2:]
			rows = append(rows, StyleRowInactive.Render(rowLine))
		}
	}

	return strings.Join(rows, "\n")
}

// renderFooter generates alert banners and context-aware keyboard mappings.
func (m *Model) renderFooter(width int) string {
	var alertBar string
	if m.ErrorText != "" {
		alertBar = StyleError.Render("ERROR: " + m.ErrorText)
	} else if m.AlertText != "" {
		alertBar = StyleAlert.Render("SUCCESS: " + m.AlertText)
	}

	var cmdBar string
	switch m.Mode {
	case ModeCommand:
		cmdBar = m.CommandInput.View()
	case ModeSearch:
		cmdBar = StyleNormalPrompt.Render("[Search Mode] ") + getSearchShortcuts(width)
	default:
		cmdBar = StyleFooter.Render(getShortcuts(width))
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

func truncate(s string, w int) string {
	if len(s) > w {
		if w > 3 {
			return s[:w-3] + "..."
		}
		return s[:w]
	}
	return s
}

// renderHelpDialog constructs the help dialog text with a title and available commands with descriptions.
func (m *Model) renderHelpDialog() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(50)

	header := titleStyle.Render("Available Commands")
	divider := lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Repeat("─", 50))

	cmdStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	rows := []string{
		header,
		divider,
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":new"), descStyle.Render("Create a new connection")),
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":edit, :e"), descStyle.Render("Edit the selected connection")),
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":delete, :del, :d"), descStyle.Render("Delete the selected connection")),
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":move, :m"), descStyle.Render("Move connection to a file/tab")),
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":quit, :q"), descStyle.Render("Quit the application")),
		fmt.Sprintf("  %-16s %s", cmdStyle.Render(":help, :h"), descStyle.Render("Show this help dialog")),
		"",
		lipgloss.NewStyle().Foreground(ColorMuted).Align(lipgloss.Center).Width(50).Render("Press Esc to close"),
	}

	return strings.Join(rows, "\n")
}
