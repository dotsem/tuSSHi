package tui

import (
	"tusshi/internal/constants"
	"tusshi/internal/ssh"
	"tusshi/internal/tui/commands"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

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
	switch {
	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyQuit):
		return m, tea.Quit

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyDown):
		if m.SelectedIndex < len(m.Filtered)-1 {
			m.SelectedIndex++
		}

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyUp):
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyLeft):
		m.navigateTabs(-1)

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyRight):
		m.navigateTabs(1)

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyPing):
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			if m.PingResults == nil {
				m.PingResults = make(map[string]*PingResult)
			}
			m.PingResults[selected.Alias] = &PingResult{Pending: true}
			return m, m.PingHost(selected)
		}

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyPingAll):
		return m, m.PingAll()

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeySearch):
		m.Mode = ModeSearch
		m.SearchInput.SetValue("")
		m.SearchInput.Focus()
		m.FilterHosts()
		return m, textinput.Blink

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyCommand):
		m.Mode = ModeCommand
		m.CommandInput.SetValue("")
		m.CommandInput.Focus()
		return m, textinput.Blink

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyAdd):
		m.FormAction = actionAdd
		m.ActiveComponent = &components.Form{
			Form:     m.BuildHostForm(m.ActiveTab),
			Validate: m.ValidateForm,
			OnSubmit: func() {
				m.executeFormSubmit()
			},
		}
		return m, m.ActiveComponent.Init()

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyEdit):
		if len(m.Filtered) > 0 {
			m.FormAction = actionEdit
			m.ActiveComponent = &components.Form{
				Form:     m.BuildHostForm(m.ActiveTab),
				Validate: m.ValidateForm,
				OnSubmit: func() {
					m.executeFormSubmit()
				},
			}
			return m, m.ActiveComponent.Init()
		}

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyDelete):
		if len(m.Filtered) > 0 {
			ctx := &cmdContext{model: m}
			_ = commands.Dispatch(":d", ctx)
			if m.ActiveComponent != nil {
				return m, m.ActiveComponent.Init()
			}
			return m, nil
		}
	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyHelp):
		m.ActiveComponent = &components.Help{
			Options: commands.HelpOptions(),
			Theme:   theme.Global,
		}
		return m, m.ActiveComponent.Init()

	case utils.MatchesMultipleStringOptions(msg.String(), constants.KeyEnter):
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
	case constants.KeyEsc, constants.KeyEnter:
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
	case constants.KeyEsc:
		m.Mode = ModeNormal
		m.CommandInput.Blur()
		return m, nil
	case constants.KeyEnter:
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
