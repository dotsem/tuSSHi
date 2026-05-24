// Package tui implements the Bubble Tea terminal user interface and interaction loops.
package tui

import (
	"fmt"
	"strings"

	"tusshi/pkg/tui/commands"
	"tusshi/pkg/tui/components"
	"tusshi/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	quitCmd         = "q, quit"
	newCmd          = "n, new"
	editCmd         = "e, edit"
	deleteCmd       = "d, rm"
	moveCmd         = "m, mv"
	helpCmd         = "h, help, ?"
	addConfigCmd    = "addconf, add-config"
	renameConfigCmd = "mvconf, rename-config"
	deleteConfigCmd = "rmconf, delete-config"
)

// helpOptions centralizes all interactive command shortcuts and their help text
var helpOptions = []components.HelpOption{
	{Shortcut: newCmd, Description: "Create a new connection"},
	{Shortcut: editCmd, Description: "Edit selected connection"},
	{Shortcut: deleteCmd, Description: "Delete selected connection"},
	{Shortcut: moveCmd, Description: "Move connection to a file/tab"},
	{Shortcut: addConfigCmd, Description: "Add a new config file"},
	{Shortcut: renameConfigCmd, Description: "Rename a config file"},
	{Shortcut: deleteConfigCmd, Description: "Delete empty config file"},
	{Shortcut: quitCmd, Description: "Quit the application"},
	{Shortcut: helpCmd, Description: "Show this help dialog"},
}

func matchesCommand(cmd string, shouldMatch string) bool {
	cmds := strings.SplitSeq(shouldMatch, ",")
	for s := range cmds {
		if cmd == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}

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
		Options: helpOptions,
		Theme:   theme.Global,
	}
}

// OpenForm sets up and opens the add/edit interactive form.
func (c *cmdContext) OpenForm(action string) {
	c.model.FormAction = action
	c.model.ActiveComponent = &components.Form{
		Form: c.model.BuildHostForm(c.model.ActiveTab),
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

// executeCommand runs commands typed into the command mode bar.
func (m *Model) executeCommand(raw string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(strings.TrimPrefix(raw, ":"))
	if len(parts) == 0 {
		return m, nil
	}

	cmd := parts[0]
	var action func(commands.Context)

	switch {
	case matchesCommand(cmd, quitCmd):
		action = commands.Quit()

	case matchesCommand(cmd, newCmd):
		action = commands.New()

	case matchesCommand(cmd, editCmd):
		action = commands.Edit(len(m.Filtered) > 0)

	case matchesCommand(cmd, deleteCmd):
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			m.ActiveComponent = &components.Confirm{
				Title:       "Delete Connection?",
				Message:     fmt.Sprintf("Are you sure you want to delete host '%s'?", selected.Alias),
				Theme:       theme.Global,
				Destructive: true,
				OnConfirm: func() tea.Cmd {
					ctx := &cmdContext{model: m}
					action := commands.Delete(m.Manager, selected)
					action(ctx)
					return ctx.cmd
				},
			}
			return m, nil
		} else {
			return m, nil
		}

	case matchesCommand(cmd, moveCmd):
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			action = commands.Move(m.Manager, selected, parts)
		} else {
			return m, nil
		}

	case matchesCommand(cmd, helpCmd):
		action = commands.Help()

	case matchesCommand(cmd, addConfigCmd):
		action = commands.AddConfig(m.Manager, parts)

	case matchesCommand(cmd, renameConfigCmd):
		action = commands.RenameConfig(m.Manager, parts)

	case matchesCommand(cmd, deleteConfigCmd):
		action = commands.DeleteConfig(m.Manager, parts)

	default:
		m.ErrorText = "Unknown command: " + cmd
		return m, nil
	}

	ctx := &cmdContext{model: m}
	action(ctx)

	return m, ctx.cmd
}
