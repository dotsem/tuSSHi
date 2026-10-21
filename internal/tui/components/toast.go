package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// ToastAction defines the semantic action type which governs border coloring.
type ToastAction int

// ActionAdd and subsequent constants define the semantic action types for toast notifications.
const (
	ActionAdd ToastAction = iota
	ActionEdit
	ActionDelete
	ActionInfo
)

const (
	colorAdd         = lipgloss.Color("#22c55e")
	colorEdit        = lipgloss.Color("#3b82f6")
	colorDelete      = lipgloss.Color("#ef4444")
	colorInfo        = lipgloss.Color("#a1a1aa")
	maxToasts        = 3
	toastTTL         = 3 * time.Second
	anchorHorizontal = 1
	anchorVertical   = 2
)

// Toast represents a single ephemeral notification item.
type Toast struct {
	ID        string
	Action    ToastAction
	Message   string
	CreatedAt time.Time
}

// ToastExpiredMsg is emitted by tea.Tick when a toast lifecycle ends.
type ToastExpiredMsg struct {
	ID string
}

// ToastQueue manages a capped stack of active toasts.
type ToastQueue struct {
	items []Toast
}

// NewToastQueue initializes an empty toast queue.
func NewToastQueue() *ToastQueue {
	return &ToastQueue{
		items: make([]Toast, 0, maxToasts),
	}
}

// Items returns a slice copy of all currently active toasts.
func (q *ToastQueue) Items() []Toast {
	res := make([]Toast, len(q.items))
	copy(res, q.items)
	return res
}

// Add pushes a new toast to the queue and returns an auto-dismiss timer command.
func (q *ToastQueue) Add(action ToastAction, message string) (Toast, tea.Cmd) {
	for len(q.items) >= maxToasts {
		q.items = q.items[1:]
	}

	t := Toast{
		ID:        fmt.Sprintf("toast-%d", time.Now().UnixNano()),
		Action:    action,
		Message:   message,
		CreatedAt: time.Now(),
	}

	q.items = append(q.items, t)

	cmd := tea.Tick(toastTTL, func(_ time.Time) tea.Msg {
		return ToastExpiredMsg{ID: t.ID}
	})

	return t, cmd
}

// Dismiss removes a specific toast from the queue by ID.
func (q *ToastQueue) Dismiss(id string) {
	filtered := q.items[:0]
	for _, item := range q.items {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}
	q.items = filtered
}

// BorderColor returns the border color mapped to the toast action.
func (t *Toast) BorderColor() lipgloss.TerminalColor {
	switch t.Action {
	case ActionAdd:
		return colorAdd
	case ActionEdit:
		return colorEdit
	case ActionDelete:
		return colorDelete
	default:
		return colorInfo
	}
}

// Render compiles an individual toast into a styled bordered string.
func (t *Toast) Render(maxWidth int) string {
	contentWidth := max(maxWidth-4, 10)
	msg := ansi.Truncate(t.Message, contentWidth, "…")

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderColor()).
		Padding(0, 1).
		Render(msg)
}

// RenderStack renders all active toasts vertically aligned to the right.
func (q *ToastQueue) RenderStack(maxWidth int) string {
	if len(q.items) == 0 {
		return ""
	}

	rendered := make([]string, len(q.items))
	for i := range q.items {
		rendered[i] = q.items[i].Render(maxWidth)
	}

	return lipgloss.JoinVertical(lipgloss.Right, rendered...)
}

// Overlay composits the active toasts onto the bottom-right corner of baseView.
func (q *ToastQueue) Overlay(baseView string, screenWidth, _ int) string {
	if len(q.items) == 0 {
		return baseView
	}

	toastBox := q.RenderStack(min(45, max(screenWidth/2, 20)))
	if toastBox == "" {
		return baseView
	}

	stackLines := strings.Split(toastBox, "\n")
	stackHeight := len(stackLines)
	stackWidth := lipgloss.Width(toastBox)

	baseLines := strings.Split(baseView, "\n")
	if len(baseLines) == 0 {
		return baseView
	}

	startRow := len(baseLines) - stackHeight - anchorVertical
	startCol := screenWidth - stackWidth - anchorHorizontal

	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	for i, stackLine := range stackLines {
		targetRow := startRow + i
		if targetRow >= len(baseLines) {
			break
		}

		baseLine := baseLines[targetRow]
		leftPart := ansi.Truncate(baseLine, startCol, "")
		leftWidth := ansi.StringWidth(leftPart)

		var pad string
		if leftWidth < startCol {
			pad = strings.Repeat(" ", startCol-leftWidth)
		}

		rightPart := ansi.TruncateLeft(baseLine, startCol+stackWidth, "")
		baseLines[targetRow] = leftPart + pad + stackLine + rightPart
	}

	return strings.Join(baseLines, "\n")
}
