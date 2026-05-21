package tui

import (
	"fmt"

	"tusshi/pkg/ssh"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// SSHFinishedMsg is returned by Bubble Tea when the subprocess SSH session terminates.
type SSHFinishedMsg struct {
	Err error
}

// Update processes terminal input and updates the state machine model losslessly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.SearchInput.Width = msg.Width - 6
		m.CommandInput.Width = msg.Width - 6
		return m, nil

	case SSHFinishedMsg:
		// SSH terminated; clear screen and restore terminal
		if msg.Err != nil {
			m.ErrorText = fmt.Sprintf("SSH session error: %v", msg.Err)
		} else {
			m.AlertText = "SSH session disconnected."
		}
		m.Reload()
		return m, tea.ClearScreen

	case tea.KeyMsg:
		// Reset temporary notifications on keypress
		m.AlertText = ""
		m.ErrorText = ""
	}

	// Delegate based on active mode
	switch m.Mode {
	case ModeForm:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == keyEsc {
			m.Mode = ModeNormal
			m.ActiveForm = nil
			return m, nil
		}
		if m.ActiveForm != nil {
			var formModel tea.Model
			formModel, cmd = m.ActiveForm.Update(msg)
			m.ActiveForm = formModel.(*huh.Form)

			switch m.ActiveForm.State {
			case huh.StateCompleted:
				m.Mode = ModeNormal
				m.executeFormSubmit()
			case huh.StateAborted:
				m.Mode = ModeNormal
				m.ActiveForm = nil
			}
		}
		return m, cmd

	case ModeSearch:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case keyEsc, keyEnter:
				m.Mode = ModeNormal
				m.SearchInput.Blur()
				return m, nil
			}
		}

		var searchCmd tea.Cmd
		m.SearchInput, searchCmd = m.SearchInput.Update(msg)
		m.FilterHosts()
		return m, searchCmd

	case ModeCommand:
		if msg, ok := msg.(tea.KeyMsg); ok {
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
		}

		var cmdCmd tea.Cmd
		m.CommandInput, cmdCmd = m.CommandInput.Update(msg)
		return m, cmdCmd

	default: // ModeNormal
		if msg, ok := msg.(tea.KeyMsg); ok {
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
				m.ActiveForm = m.BuildHostForm(m.ActiveTab)
				m.Mode = ModeForm
				return m, m.ActiveForm.Init()

			case "e":
				if len(m.Filtered) > 0 {
					m.FormAction = actionEdit
					m.ActiveForm = m.BuildHostForm(m.ActiveTab)
					m.Mode = ModeForm
					return m, m.ActiveForm.Init()
				}

			case "d":
				// Safe confirmation entry using command bar pre-population
				m.Mode = ModeCommand
				m.CommandInput.SetValue("delete")
				m.CommandInput.Focus()
				return m, textinput.Blink

			case keyEnter:
				if len(m.Filtered) > 0 {
					selected := m.Filtered[m.SelectedIndex]
					sshCmd := ssh.NewSSHCommand(selected.Alias)
					return m, tea.ExecProcess(sshCmd, func(err error) tea.Msg {
						return SSHFinishedMsg{Err: err}
					})
				}
			}
		}
	}

	return m, nil
}

// navigateTabs switches the current configuration tab focus left or right.
func (m *Model) navigateTabs(direction int) {
	if len(m.Tabs) <= 1 {
		return
	}

	currentIndex := 0
	for i, t := range m.Tabs {
		if t == m.ActiveTab {
			currentIndex = i
			break
		}
	}

	newIndex := (currentIndex + direction + len(m.Tabs)) % len(m.Tabs)
	m.ActiveTab = m.Tabs[newIndex]
	m.SelectedIndex = 0
	m.FilterHosts()
}

// executeFormSubmit saves the completed CRUD form details back to disk AST.
func (m *Model) executeFormSubmit() {
	var err error

	// Map advanced properties to the FormHost structure
	m.FormHost.SourceFile = m.FormDestFile
	if m.FormProxyJump != "" {
		m.FormHost.Properties["ProxyJump"] = m.FormProxyJump
	} else {
		delete(m.FormHost.Properties, "ProxyJump")
	}
	if m.FormForwardAgent != "" {
		m.FormHost.Properties["ForwardAgent"] = m.FormForwardAgent
	} else {
		delete(m.FormHost.Properties, "ForwardAgent")
	}

	switch m.FormAction {
	case actionAdd:
		err = m.Manager.AddHost(m.FormHost.SourceFile, m.FormHost)
		if err == nil {
			m.AlertText = fmt.Sprintf("Host %q added successfully!", m.FormHost.Alias)
		}
	case actionEdit:
		err = m.Manager.UpdateHost(m.FormOriginalAlias, m.FormHost)
		if err == nil {
			m.AlertText = fmt.Sprintf("Host %q updated successfully!", m.FormHost.Alias)
		}
	}

	if err != nil {
		m.ErrorText = "Error saving host: " + err.Error()
	}

	m.Reload()
}
