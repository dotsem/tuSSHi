package tui

import (
	"fmt"
	"strings"
	"tusshi/pkg/tui/style"
)

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
			rows = append(rows, style.RowActive.Render(rowLine))
		} else {
			rowLine = "  " + rowLine[2:]
			rows = append(rows, style.RowInactive.Render(rowLine))
		}
	}

	return strings.Join(rows, "\n")
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
