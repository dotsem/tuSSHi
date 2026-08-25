package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tusshi/internal/config"
	"tusshi/internal/constants"
	tussh "tusshi/internal/ssh"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ServiceStatus holds the live auth check state for one service host.
type ServiceStatus struct {
	Checked bool
	Result  tussh.AuthResult
}

// ServiceActionMsg is dispatched when a user triggers an action (edit, add, delete) from the Services overlay.
type ServiceActionMsg struct {
	Action string // "edit", "add", "delete"
	Host   *config.Host
}

// Services is an interactive overlay listing all IsService hosts with async auth indicators.
type Services struct {
	Hosts         []*config.Host
	SelectedIndex int
	Results       map[string]*ServiceStatus
	Theme         theme.Theme
}

// Init triggers auth checks for all service hosts immediately.
func (s *Services) Init() tea.Cmd {
	return nil
}

// Update handles navigation and action keybindings in the services overlay.
func (s *Services) Update(msg tea.Msg) (tea.Cmd, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEsc):
			return nil, true

		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyUp):
			if s.SelectedIndex > 0 {
				s.SelectedIndex--
			}
			return nil, false

		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyDown):
			if s.SelectedIndex < len(s.Hosts)-1 {
				s.SelectedIndex++
			}
			return nil, false

		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyEdit):
			if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Hosts) {
				selected := s.Hosts[s.SelectedIndex]
				return func() tea.Msg {
					return ServiceActionMsg{Action: "edit", Host: selected}
				}, true
			}

		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyAdd):
			return func() tea.Msg {
				return ServiceActionMsg{Action: "add"}
			}, true

		case utils.MatchesMultipleStringOptions(keyMsg.String(), constants.KeyDelete):
			if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Hosts) {
				selected := s.Hosts[s.SelectedIndex]
				return func() tea.Msg {
					return ServiceActionMsg{Action: "delete", Host: selected}
				}, true
			}
		}
	}
	return nil, false
}

// SetResult stores an auth check result received from a background tea.Cmd.
func (s *Services) SetResult(alias string, result tussh.AuthResult) {
	if s.Results == nil {
		s.Results = make(map[string]*ServiceStatus)
	}
	s.Results[alias] = &ServiceStatus{Checked: true, Result: result}
}

// View renders the interactive services overlay table.
func (s *Services) View(width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(s.Theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(width)

	divider := lipgloss.NewStyle().
		Foreground(s.Theme.Muted).
		Render(strings.Repeat("─", width))

	muteStyle := lipgloss.NewStyle().Foreground(s.Theme.Muted)
	onlineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	offlineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	selectedRowStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("237")).
		Bold(true)

	rows := []string{
		titleStyle.Render("SSH Services"),
		divider,
		"",
	}

	if len(s.Hosts) == 0 {
		rows = append(rows, muteStyle.Align(lipgloss.Center).Width(width).Render("No service hosts configured — press 'a' or use :svc add"))
	} else {
		aliasW := 14
		hostW := 20
		keyW := max(width-aliasW-hostW-7, 24)

		headerStyle := lipgloss.NewStyle().Foreground(s.Theme.Muted).Bold(true)
		header := fmt.Sprintf("    %-*s %-*s %s",
			aliasW, headerStyle.Render("ALIAS"),
			hostW, headerStyle.Render("HOST"),
			headerStyle.Render("KEY"),
		)
		rows = append(rows, header)
		rows = append(rows, muteStyle.Render(strings.Repeat("─", width)))

		for i, h := range s.Hosts {
			isOK := false
			hasError := false
			errMsg := ""

			if st, ok := s.Results[h.Alias]; ok && st.Checked {
				if st.Result.OK {
					isOK = true
				} else {
					hasError = true
					errMsg = st.Result.Error
				}
			}

			aliasCell := truncateStr(h.Alias, aliasW)
			hostCell := truncateStr(h.Name, hostW)
			keyCell := truncateStr(shortenPath(h.IdentityFile), keyW)

			if i == s.SelectedIndex {
				indicatorSymbol := "○"
				if isOK || hasError {
					indicatorSymbol = "●"
				}
				listLine1 := fmt.Sprintf("  %s  %-*s %-*s %s",
					indicatorSymbol,
					aliasW, aliasCell,
					hostW, hostCell,
					keyCell,
				)
				rows = append(rows, selectedRowStyle.Width(width).Render(listLine1))

				if hasError {
					errText := truncateStr(errMsg, max(width-8, 15))
					listLine2 := fmt.Sprintf("     └ %s", errText)
					rows = append(rows, selectedRowStyle.Width(width).Render(listLine2))
				}
			} else {
				indicator := muteStyle.Render("○")
				if isOK {
					indicator = onlineStyle.Render("●")
				} else if hasError {
					indicator = offlineStyle.Render("●")
				}

				line1 := fmt.Sprintf("  %s  %-*s %-*s %s",
					indicator,
					aliasW, aliasCell,
					hostW, hostCell,
					keyCell,
				)
				rows = append(rows, line1)

				if hasError {
					errText := truncateStr(errMsg, max(width-8, 15))
					line2 := fmt.Sprintf("     %s %s", muteStyle.Render("└"), offlineStyle.Render(errText))
					rows = append(rows, line2)
				}
			}
		}

	}

	rows = append(rows, "")
	rows = append(rows, muteStyle.Align(lipgloss.Center).Width(width).Render("↑/↓ navigate • e edit • a add • d delete • Esc close"))

	return strings.Join(rows, "\n")
}

func shortenPath(p string) string {
	if p == "" {
		return "—"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	rel, err := filepath.Rel(home, p)
	if err != nil || strings.HasPrefix(rel, "..") {
		return p
	}
	return "~/" + rel
}

func truncateStr(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit-1]) + "…"
}
