package tui

import (
	"path/filepath"
	"strings"

	"tusshi/internal/config"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the current keyboard input navigation focus state.
type Mode int

const (
	// ModeNormal indicates navigation and shortcut command focus.
	ModeNormal Mode = iota
	// ModeSearch indicates search filter text input focus.
	ModeSearch
	// ModeCommand indicates bottom-bar colon command input focus.
	ModeCommand
)

// Model holds the state machine parameters for the Bubble Tea application loop.
type Model struct {
	Manager       *config.Manager
	Hosts         []*config.Host
	Filtered      []*config.Host
	SelectedIndex int

	// Tabs management
	ActiveTab string
	Tabs      []string

	// Active interaction state
	Mode         Mode
	SearchInput  textinput.Model
	CommandInput textinput.Model

	// Dimensions
	Width  int
	Height int

	// Active forms state
	FormHost          *config.Host
	FormOriginalAlias string
	FormAction        string // "add" or "edit"
	FormDestFile      string
	FormProxyJump     string
	FormForwardAgent  string

	// Alerts
	AlertText string
	ErrorText string

	// Active overlay component
	ActiveComponent components.Component

	PingResults map[string]*PingResult
}

// Init initializes the Bubble Tea application state and returns initial commands.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("tuSSHi"),
		textinput.Blink,
		m.PingAll(),
	)
}

// NewModel initializes and configures a TUSSHI TUI application state model.
func NewModel(mgr *config.Manager) *Model {
	searchIn := textinput.New()
	searchIn.Placeholder = "Type to search..."
	searchIn.Prompt = "/ "

	cmdIn := textinput.New()
	cmdIn.Placeholder = "command..."
	cmdIn.Prompt = ":"

	m := &Model{
		Manager:      mgr,
		Mode:         ModeNormal,
		SearchInput:  searchIn,
		CommandInput: cmdIn,
		PingResults:  make(map[string]*PingResult),
	}
	m.Reload()

	if err := mgr.EnsureFirstRunBackup(); err != nil {
		m.ActiveComponent = &components.Alert{
			Title:   "Backup Failed",
			Message: "Could not create initial pre-tuSSHi backup:\n" + err.Error() + "\n\nPlease ensure SSH directory permissions are correct.",
			IsError: true,
			Theme:   theme.Global,
		}
	}

	return m
}

// Reload synchronizes the model state with the config files on disk.
func (m *Model) Reload() {
	if err := m.Manager.Load(); err != nil {
		m.ErrorText = "Failed to load configs: " + err.Error()
		return
	}
	m.Hosts = m.Manager.GetHosts()

	if m.PingResults != nil {
		activeAliases := make(map[string]bool)
		for _, h := range m.Hosts {
			activeAliases[h.Alias] = true
		}
		for alias := range m.PingResults {
			if !activeAliases[alias] {
				delete(m.PingResults, alias)
			}
		}
	}

	m.Tabs = []string{tabAll}
	for _, f := range m.Manager.FileOrder {
		// why: primary ssh config must not be renamed or deleted via the UI
		if f == m.Manager.PrimaryPath {
			hasConnections := false
			for _, h := range m.Hosts {
				if h.SourceFile == f {
					hasConnections = true
					break
				}
			}
			if hasConnections {
				m.Tabs = append(m.Tabs, f)
			}
			continue
		}
		m.Tabs = append(m.Tabs, f)
	}

	tabValid := false
	for _, t := range m.Tabs {
		if t == m.ActiveTab {
			tabValid = true
			break
		}
	}
	if !tabValid {
		m.ActiveTab = tabAll
	}

	m.FilterHosts()
}

// FilterHosts filters and fuzzy-matches the host list using active tab and search query.
func (m *Model) FilterHosts() {
	var filtered []*config.Host
	searchQ := strings.ToLower(m.SearchInput.Value())

	for _, h := range m.Hosts {
		if m.ActiveTab != "All" && h.SourceFile != m.ActiveTab {
			continue
		}

		// why: wildcard configs (e.g. Host *) are metadata, not connectable hosts
		if h.IsWildcard {
			continue
		}

		aliasMatch := strings.Contains(strings.ToLower(h.Alias), searchQ)
		nameMatch := strings.Contains(strings.ToLower(h.Name), searchQ)
		userMatch := strings.Contains(strings.ToLower(h.User), searchQ)

		if searchQ == "" || aliasMatch || nameMatch || userMatch {
			filtered = append(filtered, h)
		}
	}

	m.Filtered = filtered

	if m.SelectedIndex >= len(m.Filtered) {
		m.SelectedIndex = len(m.Filtered) - 1
	}
	if m.SelectedIndex < 0 {
		m.SelectedIndex = 0
	}
}

// GetTabLabel returns a clean display label (filename) for a config tab path.
func GetTabLabel(tabPath string) string {
	if tabPath == tabAll {
		return tabAll
	}
	return filepath.Base(tabPath)
}
