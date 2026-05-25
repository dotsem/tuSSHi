package tui

import (
	"regexp"
	"strings"
	"tusshi/pkg/tui/style"

	"github.com/charmbracelet/lipgloss"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// View renders the TUSSHI TUI interface based on state and window constraints.
func (m *Model) View() string {
	if m.Width < 40 || m.Height < 20 {
		return "Terminal is too small."
	}

	bgString := m.renderNormalView(m.Width, m.Height)

	var dialogContent string
	var showDialog bool

	if m.ActiveComponent != nil {
		dialogContent = m.ActiveComponent.View(54)
		showDialog = true
	}

	if showDialog {
		bgLines := strings.Split(stripANSI(bgString), "\n")

		dialogWidth := min(60, m.Width-4)
		dialogHeight := min(20, m.Height-2)

		dialogBox := style.Dialog.Width(dialogWidth).Height(dialogHeight).Render(dialogContent)
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

				mutedLeft := style.Muted.Render(string(leftPart))
				mutedRight := style.Muted.Render(string(rightPart))
				dialogLine := dialogLines[dialogLineIdx]

				finalLines = append(finalLines, mutedLeft+dialogLine+mutedRight)
			} else {
				finalLines = append(finalLines, style.Muted.Render(string(bgRunes)))
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
	headerBox := style.HeaderBox.
		Width(width - 2).
		Height(headerBoxHeight - 2).
		Render(headerContent)

	tableContent := m.renderTable(width-4, bodyBoxHeight-2)
	bodyBox := style.BodyBox.
		Width(width - 2).
		Height(bodyBoxHeight - 2).
		Render(tableContent)

	footerContent := m.renderFooter(width - 4)
	footerBox := style.FooterBox.
		Width(width - 2).
		Height(footerBoxHeight - 2).
		Render(footerContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerBox,
		bodyBox,
		footerBox,
	)
}
