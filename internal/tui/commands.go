// Package tui implements the Bubble Tea terminal user interface and interaction loops.
package tui

import (
	"fmt"
	"os"
	"strings"

	"tusshi/internal/config"
	"tusshi/internal/ssh"
	"tusshi/internal/tui/commands"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

	"github.com/atotto/clipboard"
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
	pingCmd         = "p, ping"
	pingAllCmd      = "P, pingall"
	tagCmd          = "tag"
	untagCmd        = "untag"
	serviceCmd      = "service, services, svc"
)

// helpOptions centralizes all interactive command shortcuts and their help text
var helpOptions = []components.HelpOption{
	{Shortcut: newCmd, Description: "Create a new connection"},
	{Shortcut: editCmd, Description: "Edit selected connection"},
	{Shortcut: deleteCmd, Description: "Delete selected connection"},
	{Shortcut: moveCmd, Description: "Move connection to a file/tab"},
	{Shortcut: tagCmd, Description: "Add tags to connection (:tag [alias] <tags...>)"},
	{Shortcut: untagCmd, Description: "Remove tags from connection (:untag [alias] <tags...>)"},
	{Shortcut: pingCmd, Description: "Ping selected connection"},
	{Shortcut: pingAllCmd, Description: "Ping all connections"},
	{Shortcut: serviceCmd, Description: "Manage SSH services (:svc [add|edit|rm])"},
	{Shortcut: addConfigCmd, Description: "Add a new config file"},
	{Shortcut: renameConfigCmd, Description: "Rename a config file"},
	{Shortcut: deleteConfigCmd, Description: "Delete empty config file"},
	{Shortcut: quitCmd, Description: "Quit the application"},
	{Shortcut: helpCmd, Description: "Show this help dialog"},
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

// OpenServiceForm opens the SSH service form for adding or editing a service host.
func (c *cmdContext) OpenServiceForm(action string, targetHost *config.Host) {
	state := &ServiceFormState{
		Action:      action,
		KeySource:   keySourceGenerate,
		KeyType:     ssh.KeyTypeED25519,
		PresetAlias: "github",
	}

	if action == actionEdit && targetHost != nil {
		state.OriginalAlias = targetHost.Alias
		state.HostAlias = targetHost.Alias
		state.HostName = targetHost.Name
		state.HostUser = targetHost.User
		state.KeyPath = targetHost.IdentityFile
		state.KeySource = keySourceExisting
		if preset, ok := ssh.FindPreset(targetHost.Alias); ok {
			state.PresetAlias = preset.HostName
		} else {
			state.PresetAlias = ssh.PresetCustom
		}

	}

	c.model.ActiveComponent = &components.Form{
		Form: BuildServiceForm(state),
		OnSubmit: func() {
			c.model.executeServiceFormSubmit(state)
		},
	}
	c.cmd = c.model.ActiveComponent.Init()
}

// OpenServiceEdit locates a service host by alias and opens its edit form.
func (c *cmdContext) OpenServiceEdit(alias string) {
	var found *config.Host
	for _, h := range c.model.Hosts {
		if h.IsService && h.Alias == alias {
			found = h
			break
		}
	}
	if found == nil {
		c.SetError(fmt.Sprintf("Service host %q not found", alias))
		return
	}
	c.OpenServiceForm(actionEdit, found)
}

// DeleteService prompts for confirmation and deletes a service host by alias, with optional SSH key cleanup.
func (c *cmdContext) DeleteService(alias string) {
	var found *config.Host
	for _, h := range c.model.Hosts {
		if h.IsService && h.Alias == alias {
			found = h
			break
		}
	}
	if found == nil {
		c.SetError(fmt.Sprintf("Service host %q not found", alias))
		return
	}

	keyPath := expandTildePath(found.IdentityFile)
	var hasKeyFile bool
	if keyPath != "" {
		if _, err := os.Stat(keyPath); err == nil {
			hasKeyFile = true
		}
	}

	c.model.ActiveComponent = &components.Confirm{
		Title:       "Delete Service Connection?",
		Message:     fmt.Sprintf("Are you sure you want to delete service host '%s'?", alias),
		Theme:       theme.Global,
		Destructive: true,
		OnConfirm: func() tea.Cmd {
			if !hasKeyFile {
				if err := c.model.Manager.DeleteHost(alias); err != nil {
					c.model.ErrorText = "Failed to delete service host: " + err.Error()
				} else {
					c.model.AlertText = fmt.Sprintf("Service host %q deleted", alias)
				}
				c.model.Reload()
				return nil
			}

			c.model.ActiveComponent = &components.Confirm{
				Title:       "Delete Associated SSH Key Files?",
				Message:     fmt.Sprintf("Do you also want to remove key files from disk?\n\n• Private: %s\n• Public: %s.pub", found.IdentityFile, found.IdentityFile),
				Theme:       theme.Global,
				YesStr:      " Delete Keys ",
				NoStr:       " Keep Keys ",
				Destructive: true,
				OnConfirm: func() tea.Cmd {
					_ = os.Remove(keyPath)
					_ = os.Remove(keyPath + ".pub")
					if err := c.model.Manager.DeleteHost(alias); err != nil {
						c.model.ErrorText = "Failed to delete service host: " + err.Error()
					} else {
						c.model.AlertText = fmt.Sprintf("Service host %q and SSH key files deleted", alias)
					}
					c.model.Reload()
					return nil
				},
				OnCancel: func() tea.Cmd {
					if err := c.model.Manager.DeleteHost(alias); err != nil {
						c.model.ErrorText = "Failed to delete service host: " + err.Error()
					} else {
						c.model.AlertText = fmt.Sprintf("Service host %q deleted (keys preserved)", alias)
					}
					c.model.Reload()
					return nil
				},
			}
			return nil
		},
	}
}

// OpenServices opens the services overlay and triggers background auth checks.
func (c *cmdContext) OpenServices() {
	var serviceHosts []*config.Host
	for _, h := range c.model.Hosts {
		if h.IsService {
			serviceHosts = append(serviceHosts, h)
		}
	}
	c.model.ActiveComponent = &components.Services{
		Hosts:   serviceHosts,
		Results: make(map[string]*components.ServiceStatus),
		Theme:   theme.Global,
	}
	c.cmd = c.model.CheckAllServices()
}

// executeServiceFormSubmit processes service form submission for both add and edit actions.
func (m *Model) executeServiceFormSubmit(s *ServiceFormState) {
	s.ApplyPreset()

	resolved := s.ResolvedKeyPath()

	if s.KeySource == keySourceGenerate {
		if err := ssh.GenerateKey(resolved, s.KeyType, s.KeyComment); err != nil {
			m.ErrorText = "Key generation failed: " + err.Error()
			return
		}
	}

	h := &config.Host{
		Alias:        s.HostAlias,
		Name:         s.HostName,
		User:         s.HostUser,
		IdentityFile: resolved,
		IsService:    true,
		Properties:   make(map[string]string),
	}

	var err error
	if s.Action == actionEdit {
		err = m.Manager.UpdateHost(s.OriginalAlias, h)
	} else {
		err = m.Manager.AddServiceHost(h)
	}

	if err != nil {
		m.ErrorText = "Failed to save service host: " + err.Error()
		return
	}

	m.Reload()

	if s.KeySource == keySourceGenerate {
		pubKey, err := ssh.ReadPublicKey(resolved)
		if err != nil {
			m.AlertText = fmt.Sprintf("Key created at %s — could not read public key: %s", resolved, err)
			return
		}

		_ = clipboard.WriteAll(pubKey)

		m.ActiveComponent = &components.Alert{
			Title:   "SSH Key Created — Add to " + s.HostName,
			Message: fmt.Sprintf("Public key copied to your clipboard!\n\nPaste this key into your %s account SSH settings:\n\n%s", s.HostName, pubKey),
			Theme:   theme.Global,
		}
		return
	}

	m.AlertText = fmt.Sprintf("Service host %q configured", s.HostAlias)
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
	case utils.MatchesMultipleStringOptions(cmd, quitCmd):
		action = commands.Quit()

	case utils.MatchesMultipleStringOptions(cmd, newCmd):
		action = commands.New()

	case utils.MatchesMultipleStringOptions(cmd, editCmd):
		action = commands.Edit(len(m.Filtered) > 0)

	case utils.MatchesMultipleStringOptions(cmd, deleteCmd):
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
		}
		return m, nil

	case utils.MatchesMultipleStringOptions(cmd, moveCmd):
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			action = commands.Move(m.Manager, selected, parts)
		} else {
			return m, nil
		}

	case utils.MatchesMultipleStringOptions(cmd, helpCmd):
		action = commands.Help()

	case utils.MatchesMultipleStringOptions(cmd, pingAllCmd):
		return m, m.PingAll()

	case utils.MatchesMultipleStringOptions(cmd, pingCmd):
		if len(m.Filtered) > 0 {
			selected := m.Filtered[m.SelectedIndex]
			return m, m.PingHost(selected)
		}
		return m, nil

	case utils.MatchesMultipleStringOptions(cmd, addConfigCmd):
		action = commands.AddConfig(m.Manager, parts)

	case utils.MatchesMultipleStringOptions(cmd, renameConfigCmd):
		action = commands.RenameConfig(m.Manager, parts)

	case utils.MatchesMultipleStringOptions(cmd, deleteConfigCmd):
		action = commands.DeleteConfig(m.Manager, parts)

	case utils.MatchesMultipleStringOptions(cmd, tagCmd):
		var selected *config.Host
		if len(m.Filtered) > 0 {
			selected = m.Filtered[m.SelectedIndex]
		}
		action = commands.Tag(m.Manager, selected, parts)

	case utils.MatchesMultipleStringOptions(cmd, untagCmd):
		var selected *config.Host
		if len(m.Filtered) > 0 {
			selected = m.Filtered[m.SelectedIndex]
		}
		action = commands.Untag(m.Manager, selected, parts)

	case utils.MatchesMultipleStringOptions(cmd, serviceCmd):
		subcmd := ""
		alias := ""
		if len(parts) > 1 {
			subcmd = parts[1]
		}
		if len(parts) > 2 {
			alias = parts[2]
		}
		action = commands.Service(subcmd, alias)

	default:
		m.ErrorText = "Unknown command: " + cmd
		return m, nil
	}

	ctx := &cmdContext{model: m}
	action(ctx)

	return m, ctx.cmd
}
