// Package tui implements the Bubble Tea terminal user interface and interaction loops.
package tui

import (
	"tusshi/internal/config"
	"tusshi/internal/tui/commands"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

// cmdContext implements commands.Context to proxy actions to the Model.
type cmdContext struct {
	model *Model
	cmd   tea.Cmd
}

// Quit proxies the quit command to Bubble Tea runtime.
func (c *cmdContext) Quit() {
	c.cmd = tea.Quit
}

// OpenHelp sets the active component to help overlay.
func (c *cmdContext) OpenHelp() {
	c.model.ActiveComponent = &components.Help{
		Options: commands.HelpOptions(),
		Theme:   theme.Global,
	}
}

// OpenForm sets up and opens the add/edit interactive form for connections.
func (c *cmdContext) OpenForm(action string) {
	c.model.FormAction = action
	c.model.ActiveComponent = &components.Form{
		Form:     c.model.BuildHostForm(c.model.ActiveTab),
		Validate: c.model.ValidateForm,
		OnSubmit: func() {
			c.model.executeFormSubmit()
		},
	}
	c.cmd = c.model.ActiveComponent.Init()
}

// SetAlert sets the model alert text banner.
func (c *cmdContext) SetAlert(text string) {
	c.model.AlertText = text
}

// SetError sets the model error text banner.
func (c *cmdContext) SetError(text string) {
	c.model.ErrorText = text
}

// Reload reloads the configurations.
func (c *cmdContext) Reload() {
	c.model.Reload()
}

// GetActiveTab returns the model's active tab path.
func (c *cmdContext) GetActiveTab() string {
	return c.model.ActiveTab
}

// SetActiveTab sets the model's active tab path.
func (c *cmdContext) SetActiveTab(tab string) {
	c.model.ActiveTab = tab
}

// GetSelectedHost returns the currently selected host from the filtered list.
func (c *cmdContext) GetSelectedHost() *config.Host {
	if len(c.model.Filtered) > 0 && c.model.SelectedIndex >= 0 && c.model.SelectedIndex < len(c.model.Filtered) {
		return c.model.Filtered[c.model.SelectedIndex]
	}
	return nil
}

// GetManager returns the configuration manager instance.
func (c *cmdContext) GetManager() *config.Manager {
	return c.model.Manager
}

// Confirm sets the active component to a confirmation dialog overlay.
func (c *cmdContext) Confirm(title, message string, destructive bool, onConfirm func()) {
	c.model.ActiveComponent = &components.Confirm{
		Title:       title,
		Message:     message,
		Theme:       theme.Global,
		Destructive: destructive,
		OnConfirm: func() tea.Cmd {
			onConfirm()
			return nil
		},
	}
}

// PingHost triggers a background ping check for a single host.
func (c *cmdContext) PingHost(host *config.Host) {
	c.cmd = c.model.PingHost(host)
}

// PingAll triggers background ping checks for all hosts.
func (c *cmdContext) PingAll() {
	c.cmd = c.model.PingAll()
}

// executeCommand runs commands typed into the command mode bar via the commands registry.
func (m *Model) executeCommand(raw string) (tea.Model, tea.Cmd) {
	ctx := &cmdContext{model: m}
	if err := commands.Dispatch(raw, ctx); err != nil {
		m.ErrorText = err.Error()
		return m, nil
	}
	return m, ctx.cmd
}
