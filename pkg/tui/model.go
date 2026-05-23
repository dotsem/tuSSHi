package tui

import (
	"path/filepath"
	"strings"

	"tusshi/pkg/config"
	"tusshi/pkg/tui/components"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
	// ModeForm indicates active interactive CRUD form focus.
	ModeForm
	// ModeHelp indicates interactive help dialog focus.
	ModeHelp
	// ModeConfirm indicates active confirmation dialog focus.
	ModeConfirm
)

const (
	actionAdd  = "add"
	actionEdit = "edit"
	tabAll     = "All"
	keyEsc     = "esc"
	keyEnter   = "enter"
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

	// Active forms
	ActiveForm        *huh.Form
	FormHost          *config.Host
	FormOriginalAlias string
	FormAction        string // "add" or "edit"
	FormDestFile      string
	FormProxyJump     string
	FormForwardAgent  string

	// Alerts
	AlertText string
	ErrorText string

	// Reusable components
	ConfirmComponent *components.Confirm
}

// Init initializes the Bubble Tea application state and returns initial commands.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("tuSSHi"),
		textinput.Blink,
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
	}
	m.Reload()
	return m
}

// Reload synchronizes the model state with the config files on disk.
func (m *Model) Reload() {
	if err := m.Manager.Load(); err != nil {
		m.ErrorText = "Failed to load configs: " + err.Error()
		return
	}
	m.Hosts = m.Manager.GetHosts()

	m.Tabs = []string{tabAll}
	for _, f := range m.Manager.FileOrder {
		// Filter out the "config" file as it can't be renamed or deleted and just creates confusion.
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
		// All other config files are added as tabs
		m.Tabs = append(m.Tabs, f)
	}

	// Default active tab to All if not set or invalid
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
		// Tab matching
		if m.ActiveTab != "All" && h.SourceFile != m.ActiveTab {
			continue
		}

		// Skip wildcards from the main listing to keep it connection-focused
		if h.IsWildcard {
			continue
		}

		// Substring match on Alias, Address (Name), or User
		aliasMatch := strings.Contains(strings.ToLower(h.Alias), searchQ)
		nameMatch := strings.Contains(strings.ToLower(h.Name), searchQ)
		userMatch := strings.Contains(strings.ToLower(h.User), searchQ)

		if searchQ == "" || aliasMatch || nameMatch || userMatch {
			filtered = append(filtered, h)
		}
	}

	m.Filtered = filtered

	// Clamp selected index within range
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
