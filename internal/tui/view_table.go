package tui

import (
	"fmt"
	"strings"
	"tusshi/internal/config"
	"tusshi/internal/tui/style"

	"github.com/charmbracelet/lipgloss"
)

// renderTable draws the formatted column grid displaying active connections.
func (m *Model) renderTable(width, maxHeight int) string {
	if len(m.Filtered) == 0 {
		return "\n  No connections found. Press 'a' to add a connection, or ':' for help."
	}

	var headerRow, dividerRow string
	var wAlias, wName, wUser, wPort, wStatus, wConfig int

	switch {
	case width >= 85:
		wTotal := max(width-14, 10)
		wAlias = int(float64(wTotal) * 0.15)
		wName = int(float64(wTotal) * 0.30)
		wUser = int(float64(wTotal) * 0.12)
		wPort = int(float64(wTotal) * 0.08)
		wConfig = int(float64(wTotal) * 0.23)
		wStatus = wTotal - wAlias - wName - wUser - wPort - wConfig

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
			wPort, "PORT",
			wConfig, "CONFIG",
			wStatus, "STATUS",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
			strings.Repeat("─", wPort),
			strings.Repeat("─", wConfig),
			strings.Repeat("─", wStatus),
		)
	case width >= 65:
		wTotal := max(width-12, 10)
		wAlias = int(float64(wTotal) * 0.20)
		wName = int(float64(wTotal) * 0.35)
		wUser = int(float64(wTotal) * 0.15)
		wConfig = int(float64(wTotal) * 0.18)
		wStatus = wTotal - wAlias - wName - wUser - wConfig

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
			wConfig, "CONFIG",
			wStatus, "STATUS",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
			strings.Repeat("─", wConfig),
			strings.Repeat("─", wStatus),
		)
	case width >= 45:
		wTotal := max(width-8, 10)
		wAlias = int(float64(wTotal) * 0.30)
		wStatus = 6
		wName = wTotal - wAlias - wStatus

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wStatus, "STATUS",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wStatus),
		)
	default:
		wTotal := max(width-6, 10)
		wAlias = int(float64(wTotal) * 0.35)
		wStatus = 2
		wName = wTotal - wAlias - wStatus

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wStatus, "S",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wStatus),
		)
	}

	renderedHeader := style.TableHeader.Render(headerRow)
	renderedDivider := style.TableHeader.Render(dividerRow)

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
		rows = append(rows, m.renderRow(h, idx, wAlias, wName, wUser, wPort, wStatus, wConfig))
	}

	return strings.Join(rows, "\n")
}

// renderRow constructs a formatted row, applying specific colors for the status column
// and blending background colors correctly when the row is active/selected.
func (m *Model) renderRow(h *config.Host, idx int, wAlias, wName, wUser, wPort, wStatus, wConfig int) string {
	rowActive := idx == m.SelectedIndex

	var cells []string

	alias := truncate(h.Alias, wAlias)
	cells = append(cells, renderCell(alias, rowCellStyle(rowActive, "252"), rowActive, wAlias))

	name := truncate(h.Name, wName)
	cells = append(cells, renderCell(name, rowCellStyle(rowActive, "250"), rowActive, wName))

	if wUser > 0 {
		user := truncate(h.User, wUser)
		cells = append(cells, renderCell(user, rowCellStyle(rowActive, "245"), rowActive, wUser))
	}

	if wPort > 0 {
		port := h.Port
		if port == "" {
			port = "22"
		}
		port = truncate(port, wPort)
		cells = append(cells, renderCell(port, rowCellStyle(rowActive, "242"), rowActive, wPort))
	}

	if wConfig > 0 {
		cfgNickname := strings.TrimSuffix(GetTabLabel(h.SourceFile), ".conf")
		cfgNickname = strings.TrimSuffix(cfgNickname, "config")
		cfgNickname = truncate(cfgNickname, wConfig)
		cells = append(cells, renderCell(cfgNickname, rowCellStyle(rowActive, "240"), rowActive, wConfig))
	}

	if wStatus > 0 {
		statusCell := m.renderStatusCell(h.Alias, rowActive, wStatus)
		cells = append(cells, statusCell)
	}

	rowContent := strings.Join(cells, "  ")

	prefix := "  "
	if rowActive {
		prefix = "❯ "
	}

	var prefixStyle lipgloss.Style
	if rowActive {
		prefixStyle = lipgloss.NewStyle().Background(lipgloss.Color("237")).Foreground(lipgloss.Color("#FF5500")).Bold(true)
	} else {
		prefixStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	}

	return prefixStyle.Render(prefix) + rowContent
}

func truncate(s string, w int) string {
	runes := []rune(s)
	if len(runes) > w {
		if w > 3 {
			return string(runes[:w-3]) + "..."
		}
		return string(runes[:w])
	}
	return s
}

func rowCellStyle(rowActive bool, normalColor string) lipgloss.Style {
	if rowActive {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(normalColor))
}
