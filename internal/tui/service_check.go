package tui

import (
	"tusshi/internal/config"
	tussh "tusshi/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

// ServiceCheckResult carries the auth check outcome for a single service host.
type ServiceCheckResult struct {
	Alias  string
	Result tussh.AuthResult
}

// CheckServiceAuth runs the SSH auth check for a service host as a background tea.Cmd.
func CheckServiceAuth(h *config.Host) tea.Cmd {
	return func() tea.Msg {
		return ServiceCheckResult{
			Alias:  h.Alias,
			Result: tussh.CheckAuth(h.Alias),
		}
	}
}

// CheckAllServices returns a batch command to auth-check every service host.
func (m *Model) CheckAllServices() tea.Cmd {
	var cmds []tea.Cmd
	for _, h := range m.Hosts {
		if h.IsService {
			cmds = append(cmds, CheckServiceAuth(h))
		}
	}
	return tea.Batch(cmds...)
}
