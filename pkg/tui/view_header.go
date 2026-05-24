package tui

import (
	"github.com/charmbracelet/lipgloss"
)

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
