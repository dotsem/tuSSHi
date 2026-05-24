package tui

import (
	"fmt"

	"tusshi/pkg/ssh"
	"tusshi/pkg/tui/commands"
	"tusshi/pkg/tui/components"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyMsg routes and processes key presses based on the active UI mode.
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.Mode {
	case ModeSearch:
		return m.handleSearchKey(msg)
	case ModeCommand:
		return m.handleCommandKey(msg)
	default:
		return m.handleNormalKey(msg)
	}
}

// handleNormalKey processes keyboard shortcuts when the application is in normal mode.
func (m *Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "j", "down":
		if m.SelectedIndex < len(m.Filtered)-1 {
			m.SelectedIndex++
		}

	case "k", "up":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}

	case "h", "left":
		m.navigateTabs(-1)

	case "l", "right":
		m.navigateTabs(1)

	case "/":
		m.Mode = ModeSearch
		m.SearchInput.SetValue("")
		m.SearchInput.Focus()
		m.FilterHosts()
		return m, textinput.Blink

	case ":":
		m.Mode = ModeCommand
		m.CommandInput.SetValue("")
		m.CommandInput.Focus()
		return m, textinput.Blink

	case "a":
		m.FormAction = actionAdd
		m.ActiveComponent = components.NewForm(m.BuildHostForm(m.ActiveTab), func() {
			m.executeFormSubmit()
		})
		return m, m.ActiveComponent.Init()

	case "e":
		if len(m.Filtered) > 0 {
			m.FormAction = actionEdit
			m.ActiveComponent = components.NewForm(m.BuildHostForm(m.ActiveTab), func() {
				m.executeFormSubmit()
			})
			return m, m.ActiveComponent.Init()
		}

	case "d":
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			m.ActiveComponent = components.NewConfirm(
				"Delete Connection?",
				fmt.Sprintf("Are you sure you want to delete host '%s'?", selected.Alias),
				func() tea.Cmd {
					ctx := &cmdContext{model: m}
					action := commands.Delete(m.Manager, selected)
					action(ctx)
					return ctx.cmd
				},
			)
			return m, m.ActiveComponent.Init()
		}
	case "?", ",":
		m.ActiveComponent = components.NewHelp(helpOptions, ColorPrimary, ColorMuted)
		return m, m.ActiveComponent.Init()

	case keyEnter:
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			sshCmd := ssh.NewSSHCommand(selected.Alias)
			return m, tea.ExecProcess(sshCmd, func(err error) tea.Msg {
				return SSHFinishedMsg{Err: err}
			})
		}
	}

	return m, nil
}

// handleSearchKey processes keyboard input when performing a text search.
func (m *Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyEsc, keyEnter:
		m.Mode = ModeNormal
		m.SearchInput.Blur()
		return m, nil
	}

	var searchCmd tea.Cmd
	m.SearchInput, searchCmd = m.SearchInput.Update(msg)
	m.FilterHosts()
	return m, searchCmd
}

// handleCommandKey processes keyboard input when typing command-line instructions.
func (m *Model) handleCommandKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyEsc:
		m.Mode = ModeNormal
		m.CommandInput.Blur()
		return m, nil
	case keyEnter:
		rawCmd := m.CommandInput.Value()
		m.Mode = ModeNormal
		m.CommandInput.Blur()
		m.CommandInput.SetValue("")
		return m.executeCommand(rawCmd)
	}

	var cmdCmd tea.Cmd
	m.CommandInput, cmdCmd = m.CommandInput.Update(msg)
	return m, cmdCmd
}
