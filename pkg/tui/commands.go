// Package tui implements the Bubble Tea terminal user interface and interaction loops.
package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// executeCommand runs commands typed into the command mode bar.
func (m *Model) executeCommand(raw string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(strings.TrimPrefix(raw, ":"))
	if len(parts) == 0 {
		return m, nil
	}

	cmd := parts[0]
	switch cmd {
	case "q", "quit":
		return m, tea.Quit

	case "new":
		m.FormAction = actionAdd
		m.ActiveForm = m.BuildHostForm(m.ActiveTab)
		m.Mode = ModeForm
		return m, m.ActiveForm.Init()

	case actionEdit, "e":
		if len(m.Filtered) > 0 {
			m.FormAction = actionEdit
			m.ActiveForm = m.BuildHostForm(m.ActiveTab)
			m.Mode = ModeForm
			return m, m.ActiveForm.Init()
		}

	case "delete", "del", "d":
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			if err := m.Manager.DeleteHost(selected.Alias); err == nil {
				m.AlertText = fmt.Sprintf("Deleted connection %q.", selected.Alias)
			} else {
				m.ErrorText = "Delete error: " + err.Error()
			}
			m.Reload()
		}

	case "move", "m":
		if len(m.Filtered) > 0 {
			if len(parts) < 2 {
				m.ErrorText = "Usage: :move <target-file-nickname>"
				return m, nil
			}

			// Find matched filepath from file order nicknames
			targetNickname := parts[1]
			var matchedFile string
			for _, file := range m.Manager.FileOrder {
				if filepath.Base(file) == targetNickname || strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) == targetNickname {
					matchedFile = file
					break
				}
			}

			if matchedFile == "" {
				// If target doesn't exist, create it relative to primary dir
				matchedFile = filepath.Join(filepath.Dir(m.Manager.PrimaryPath), targetNickname)
			}

			selected := m.Filtered[m.SelectedIndex]
			if err := m.Manager.MoveHost(selected.Alias, matchedFile); err == nil {
				m.AlertText = fmt.Sprintf("Moved %q to %s.", selected.Alias, filepath.Base(matchedFile))
			} else {
				m.ErrorText = "Move error: " + err.Error()
			}
			m.Reload()
		}

	default:
		m.ErrorText = "Unknown command: " + cmd
	}

	return m, nil
}
