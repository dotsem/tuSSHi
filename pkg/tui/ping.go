package tui

import (
	"fmt"
	"strings"
	"time"
	"tusshi/pkg/config"
	"tusshi/pkg/ping"
	"tusshi/pkg/tui/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PingResult represents the status and latency of a single host connection test.
type PingResult struct {
	Online  bool
	Latency float64
	Pending bool
}

// PingResultMsg carries the outcome of a background ping check.
type PingResultMsg struct {
	Alias   string
	Online  bool
	Latency float64
}

// Limit concurrency to 15 to avoid fd exhaustion or firewall bans.
var pingSemaphore = make(chan struct{}, 15)

// PingHost performs a non-blocking TCP dial check against the host.
// It resolves the target hostname and defaults the port to 22.
func (m *Model) PingHost(h *config.Host) tea.Cmd {
	return func() tea.Msg {
		pingSemaphore <- struct{}{}
		defer func() { <-pingSemaphore }()

		target := h.Name
		if target == "" {
			target = h.Alias
		}

		res := ping.Ping(target, h.Port, 1500*time.Millisecond)

		return PingResultMsg{
			Alias:   h.Alias,
			Online:  res.Online,
			Latency: res.Latency,
		}
	}
}

// PingAll returns a batch command to ping all active hosts in the TUI list.
// It pre-allocates and clears prior results, setting them to pending.
func (m *Model) PingAll() tea.Cmd {
	if len(m.Hosts) == 0 {
		return nil
	}

	var cmds []tea.Cmd
	for _, h := range m.Hosts {
		if h.IsWildcard {
			continue
		}
		if m.PingResults == nil {
			m.PingResults = make(map[string]*PingResult)
		}
		m.PingResults[h.Alias] = &PingResult{Pending: true}
		cmds = append(cmds, m.PingHost(h))
	}
	return tea.Batch(cmds...)
}

// renderStatusCell outputs a formatted Status string for the host table row.
// Blends row active backgrounds and column width limits correctly.
func (m *Model) renderStatusCell(alias string, rowActive bool, width int) string {
	res, exists := m.PingResults[alias]
	isSmall := width < 10

	if !exists {
		text := "Pending"
		if isSmall {
			text = "○"
		}
		return renderCell(text, style.Muted, rowActive, width)
	}
	if res.Pending {
		text := "Checking..."
		if isSmall {
			text = "○"
		}
		return renderCell(text, style.Muted, rowActive, width)
	}
	if res.Online {
		text := fmt.Sprintf("Online (%.0fms)", res.Latency)
		if isSmall {
			text = fmt.Sprintf("● %.0fms", res.Latency)
		}
		return renderCell(text, style.StatusOnline, rowActive, width)
	}
	text := "Offline"
	if isSmall {
		text = "ø"
	}
	return renderCell(text, style.StatusOffline, rowActive, width)
}

// renderCell renders padded text, preserving background styling for selected rows.
func renderCell(text string, cellStyle lipgloss.Style, rowActive bool, width int) string {
	runes := []rune(text)
	if len(runes) > width {
		text = truncate(text, width)
		runes = []rune(text)
	}
	visibleLen := len(runes)

	finalStyle := cellStyle
	if rowActive {
		finalStyle = finalStyle.Background(lipgloss.Color("237"))
	}

	styled := finalStyle.Render(text)

	padLen := max(0, width-visibleLen)
	var padding string
	if padLen > 0 {
		padSpace := strings.Repeat(" ", padLen)
		if rowActive {
			padding = lipgloss.NewStyle().Background(lipgloss.Color("237")).Render(padSpace)
		} else {
			padding = padSpace
		}
	}

	return styled + padding
}
