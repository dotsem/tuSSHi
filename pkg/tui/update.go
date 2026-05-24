package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// SSHFinishedMsg is returned by Bubble Tea when the subprocess SSH session terminates.
type SSHFinishedMsg struct {
	Err error
}

// Update processes terminal input and updates the state machine model losslessly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, tea.Batch(tea.ClearScreen, m.PingAll())

	case PingResultMsg:
		if m.PingResults == nil {
			m.PingResults = make(map[string]*PingResult)
		}
		m.PingResults[msg.Alias] = &PingResult{
			Online:  msg.Online,
			Latency: msg.Latency,
			Pending: false,
		}
		return m, nil

	case tea.KeyMsg:
		// Reset temporary notifications on keypress
		m.AlertText = ""
		m.ErrorText = ""
	}

	// Delegate to active overlay component if one is open
	if m.ActiveComponent != nil {
		var activeCmd tea.Cmd
		activeCmd, done := m.ActiveComponent.Update(msg)
		if done {
			m.ActiveComponent = nil
			return m, tea.Batch(activeCmd, m.PingAll())
		}
		return m, activeCmd
	}

	// Delegate based on active mode
	if msg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(msg)
	}

	// Forward non-key messages (such as blink timers) to active text inputs
	switch m.Mode {
	case ModeSearch:
		var searchCmd tea.Cmd
		m.SearchInput, searchCmd = m.SearchInput.Update(msg)
		return m, searchCmd
	case ModeCommand:
		var cmdCmd tea.Cmd
		m.CommandInput, cmdCmd = m.CommandInput.Update(msg)
		return m, cmdCmd
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
