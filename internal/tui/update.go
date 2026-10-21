package tui

import (
	"fmt"

	"tusshi/internal/tui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// SSHFinishedMsg is returned by Bubble Tea when the subprocess SSH session terminates.
type SSHFinishedMsg struct {
	Err error
}

// Update processes terminal input and updates the state machine model losslessly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := m.handleUpdate(msg)
	if mod, ok := newModel.(*Model); ok {
		return mod, mod.WrapCmd(cmd)
	}
	return newModel, cmd
}

func (m *Model) handleUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.SearchInput.Width = msg.Width - 6
		m.CommandInput.Width = msg.Width - 6
		return m, nil

	case components.ToastExpiredMsg:
		if m.Toasts != nil {
			m.Toasts.Dismiss(msg.ID)
		}
		return m, nil

	case SSHFinishedMsg:
		var toastCmd tea.Cmd
		if msg.Err != nil {
			toastCmd = m.AddToast(components.ActionDelete, fmt.Sprintf("SSH error: %v", msg.Err))
		} else {
			toastCmd = m.AddToast(components.ActionInfo, "SSH session disconnected")
		}
		m.Reload()
		return m, tea.Batch(tea.ClearScreen, m.PingAll(), toastCmd)

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

	case ServiceCheckResult:
		if svc, ok := m.ActiveComponent.(*components.Services); ok {
			svc.SetResult(msg.Alias, msg.Result)
		}
		if m.PendingServiceChecks > 0 {
			m.PendingServiceChecks--
			if !msg.Result.OK {
				m.FailedServiceNames = append(m.FailedServiceNames, msg.Alias)
			}
			if m.PendingServiceChecks == 0 {
				var cmd tea.Cmd
				if len(m.FailedServiceNames) == 0 {
					cmd = m.AddToast(components.ActionAdd, "All services connected")
				} else {
					cmd = m.AddToast(components.ActionDelete, fmt.Sprintf("%d service(s) failed connection", len(m.FailedServiceNames)))
				}
				m.FailedServiceNames = nil
				return m, cmd
			}
		}
		return m, nil

	case components.ServiceActionMsg:
		ctx := &cmdContext{model: m}
		switch msg.Action {
		case "add": // TODO generalize string
			ctx.OpenServiceForm("add", nil)
		case "edit":
			ctx.OpenServiceForm("edit", msg.Host)
		case "delete":
			if msg.Host != nil {
				ctx.DeleteService(msg.Host.Alias)
			}
		case "ping":
			if msg.Host != nil {
				return m, CheckServiceAuth(msg.Host)
			}
		case "pingall":
			return m, m.CheckAllServices()
		}
		return m, ctx.cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		m.AlertText = ""
		m.ErrorText = ""
	}

	if m.ActiveComponent != nil {
		currentComponent := m.ActiveComponent
		var activeCmd tea.Cmd
		activeCmd, done := m.ActiveComponent.Update(msg)
		if done && m.ActiveComponent == currentComponent {
			m.ActiveComponent = nil
			return m, tea.Batch(activeCmd, m.PingAll())
		}
		return m, activeCmd
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(msg)
	}

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
