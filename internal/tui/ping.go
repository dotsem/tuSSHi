package tui

import (
	"fmt"
	"strings"
	"time"
	"tusshi/internal/config"
	"tusshi/internal/ping"
	"tusshi/internal/tui/style"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	statusOnline   = "●"
	statusChecking = "○"
	statusOffline  = "✗"

	smallStatusWidth = 10
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

// why: limit concurrency to 15 to avoid fd exhaustion or firewall bans
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

// renderStatusCell outputs a formatted Status string for the host table row,
// prefixed with a symbol: filled circle online, hollow circle checking, cross offline.
// Blends row active backgrounds and column width limits correctly.
func (m *Model) renderStatusCell(alias string, rowActive bool, width int) string {
	res, exists := m.PingResults[alias]
	isSmall := width < smallStatusWidth

	if !exists || res.Pending {
		text := statusChecking + " Checking..."
		if isSmall {
			text = statusChecking
		}
		return renderCell(text, style.Muted, rowActive, width)
	}
	if res.Online {
		text := fmt.Sprintf("%s Online (%.0fms)", statusOnline, res.Latency)
		if isSmall {
			text = fmt.Sprintf("%s %.0fms", statusOnline, res.Latency)
		}
		return renderCell(text, style.StatusOnline, rowActive, width)
	}
	text := statusOffline + " Offline"
	if isSmall {
		text = statusOffline
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
		finalStyle = finalStyle.Background(theme.Global.Highlight)
	}

	styled := finalStyle.Render(text)

	padLen := max(0, width-visibleLen)
	var padding string
	if padLen > 0 {
		padSpace := strings.Repeat(" ", padLen)
		if rowActive {
			padding = lipgloss.NewStyle().Background(theme.Global.Highlight).Render(padSpace)
		} else {
			padding = padSpace
		}
	}

	return styled + padding
}
