// Package tui implements the Bubble Tea terminal user interface and interaction loops.
package tui

import (
	"fmt"
	"strings"

	"tusshi/pkg/tui/commands"
	"tusshi/pkg/tui/components"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	quitCmd         = "q, quit"
	newCmd          = "new"
	editCmd         = "e, edit"
	deleteCmd       = "d, del, delete, rm"
	moveCmd         = "m, move, mv"
	helpCmd         = "h, help"
	addConfigCmd    = "add-config, config-add, new-config"
	renameConfigCmd = "rename-config, config-rename, mvconfig, mvconf"
	deleteConfigCmd = "delete-config, config-delete, rmconfig, rmconf"
)

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

// OpenHelp sets the model mode to help.
func (c *cmdContext) OpenHelp() {
	c.model.Mode = ModeHelp
}

// OpenForm sets up and opens the add/edit interactive form.
func (c *cmdContext) OpenForm(action string) {
	c.model.FormAction = action
	c.model.ActiveForm = c.model.BuildHostForm(c.model.ActiveTab)
	c.model.Mode = ModeForm
	c.cmd = c.model.ActiveForm.Init()
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
			m.ConfirmComponent = components.NewConfirm(
				"Delete Connection?",
				fmt.Sprintf("Are you sure you want to delete host '%s'?", selected.Alias),
			)
			m.Mode = ModeConfirm
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
