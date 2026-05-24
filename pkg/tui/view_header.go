package tui

import (
	"tusshi/pkg/tui/style"

	"github.com/charmbracelet/lipgloss"
)

// renderHeader constructs the config file tabs bar and active search filter.
func (m *Model) renderHeader() string {
	var tabs []string
	for _, t := range m.Tabs {
		label := GetTabLabel(t)
		if t == m.ActiveTab {
			tabs = append(tabs, style.TabActive.Render(label))
		} else {
			tabs = append(tabs, style.TabInactive.Render(label))
		}
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	title := style.Title.Render(" tuSSHi ")
	headerTop := lipgloss.JoinHorizontal(lipgloss.Left, title, "── ", tabsRow)

	var searchBar string
	if m.Mode == ModeSearch {
		searchBar = m.SearchInput.View()
	} else {
		val := m.SearchInput.Value()
		if val != "" {
			searchBar = style.NormalPrompt.Render("/ ") + val
		} else {
			searchBar = style.Footer.Render("/ type to search...")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		headerTop,
		searchBar,
	)
}
