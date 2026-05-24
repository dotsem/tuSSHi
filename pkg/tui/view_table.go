package tui

import (
	"fmt"
	"strings"
	"tusshi/pkg/config"
	"tusshi/pkg/tui/style"

	"github.com/charmbracelet/lipgloss"
)

// renderTable draws the formatted column grid displaying active connections.
func (m *Model) renderTable(width, maxHeight int) string {
	if len(m.Filtered) == 0 {
		return "\n  No connections found. Press 'a' to add a connection, or ':' for help."
	}

	var headerRow, dividerRow string
	var wAlias, wName, wUser, wPort, wStatus, wConfig int

	// adaptive column allocation to prevent overflow at small terminal widths
	switch {
	case width >= 85:
		wTotal := max(width-14, 10)
		wAlias = int(float64(wTotal) * 0.15)
		wName = int(float64(wTotal) * 0.20)
		wUser = int(float64(wTotal) * 0.12)
		wPort = int(float64(wTotal) * 0.08)
		wStatus = int(float64(wTotal) * 0.25)
		wConfig = wTotal - wAlias - wName - wUser - wPort - wStatus

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
			wPort, "PORT",
			wStatus, "STATUS",
			wConfig, "CONFIG",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
			strings.Repeat("─", wPort),
			strings.Repeat("─", wStatus),
			strings.Repeat("─", wConfig),
		)
	case width >= 65:
		wTotal := max(width-12, 10)
		wAlias = int(float64(wTotal) * 0.20)
		wName = int(float64(wTotal) * 0.25)
		wUser = int(float64(wTotal) * 0.15)
		wStatus = int(float64(wTotal) * 0.20)
		wConfig = wTotal - wAlias - wName - wUser - wStatus

		headerRow = fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
			wAlias, "ALIAS",
			wName, "NAME / ADDRESS",
			wUser, "USER",
			wStatus, "STATUS",
			wConfig, "CONFIG",
		)
		dividerRow = fmt.Sprintf("  %s  %s  %s  %s  %s",
			strings.Repeat("─", wAlias),
			strings.Repeat("─", wName),
			strings.Repeat("─", wUser),
			strings.Repeat("─", wStatus),
			strings.Repeat("─", wConfig),
		)
	case width >= 45:
		wTotal := max(width-8, 10)
		wAlias = int(float64(wTotal) * 0.30)
		wStatus = int(float64(wTotal) * 0.30)
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
		rows = append(rows, m.renderRow(h, idx, width, wAlias, wName, wUser, wPort, wStatus))
	}

	return strings.Join(rows, "\n")
}

// renderRow constructs a formatted row, applying specific colors for the status column
// and blending background colors correctly when the row is active/selected.
func (m *Model) renderRow(h *config.Host, idx int, wAlias, wName, wUser, wPort, wStatus, wConfig int) string {
	rowActive := idx == m.SelectedIndex

	var cells []string

	// Alias cell
	alias := truncate(h.Alias, wAlias)
	var aliasStyle lipgloss.Style
	if rowActive {
		aliasStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	} else {
		aliasStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	}
	cells = append(cells, renderCell(alias, aliasStyle, rowActive, wAlias))

	// Name cell
	name := truncate(h.Name, wName)
	var nameStyle lipgloss.Style
	if rowActive {
		nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	} else {
		nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	}
	cells = append(cells, renderCell(name, nameStyle, rowActive, wName))

	// User cell
	if wUser > 0 {
		user := truncate(h.User, wUser)
		var userStyle lipgloss.Style
		if rowActive {
			userStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
		} else {
			userStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		}
		cells = append(cells, renderCell(user, userStyle, rowActive, wUser))
	}

	// Port cell
	if wPort > 0 {
		port := h.Port
		if port == "" {
			port = "22"
		}
		port = truncate(port, wPort)
		var portStyle lipgloss.Style
		if rowActive {
			portStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
		} else {
			portStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
		}
		cells = append(cells, renderCell(port, portStyle, rowActive, wPort))
	}

	// Status cell
	if wStatus > 0 {
		statusCell := m.renderStatusCell(h.Alias, rowActive, wStatus)
		cells = append(cells, statusCell)
	}

	// Config cell
	if wConfig > 0 {
		cfgNickname := strings.TrimSuffix(GetTabLabel(h.SourceFile), ".conf")
		cfgNickname = strings.TrimSuffix(cfgNickname, "config")
		cfgNickname = truncate(cfgNickname, wConfig)
		var cfgStyle lipgloss.Style
		if rowActive {
			cfgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
		} else {
			cfgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		}
		cells = append(cells, renderCell(cfgNickname, cfgStyle, rowActive, wConfig))
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
	if len(s) > w {
		if w > 3 {
			return s[:w-3] + "..."
		}
		return s[:w]
	}
	return s
}
