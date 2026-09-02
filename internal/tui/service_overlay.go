package tui

import (
	"fmt"
	"os"

	"tusshi/internal/config"
	"tusshi/internal/ssh"
	"tusshi/internal/tui/components"
	"tusshi/internal/tui/theme"
	"tusshi/internal/utils"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// OpenServiceForm opens the SSH service form for adding or editing a service host.
func (c *cmdContext) OpenServiceForm(action string, targetHost *config.Host) {
	state := &ServiceFormState{
		Action:      action,
		KeySource:   keySourceGenerate,
		KeyType:     ssh.KeyTypeED25519,
		PresetAlias: "github.com",
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
		} else if preset, ok := ssh.FindPreset(targetHost.Name); ok {
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

	keyPath := utils.ExpandTilde(found.IdentityFile)
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

	keyPath := s.KeyPath
	if keyPath == "" {
		keyPath = resolved
	}

	h := &config.Host{
		Alias:        s.HostAlias,
		Name:         s.HostName,
		User:         s.HostUser,
		IdentityFile: keyPath,
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
