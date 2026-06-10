package tui

import (
	"path/filepath"

	"tusshi/internal/config"
	"tusshi/internal/validation"

	"github.com/charmbracelet/huh"
)

// BuildHostForm creates a beautiful multi-step interactive form using Huh
// for adding or editing an SSH connection. It accommodates standard fields
// and common advanced settings cleanly.
func (m *Model) BuildHostForm(defaultFile string) *huh.Form {
	m.FormHost = &config.Host{
		Properties: make(map[string]string),
	}

	m.FormDestFile = defaultFile
	m.FormProxyJump = ""
	m.FormForwardAgent = "no"

	if m.FormAction == actionEdit && m.SelectedIndex < len(m.Filtered) {
		selected := m.Filtered[m.SelectedIndex]
		m.FormOriginalAlias = selected.Alias
		m.FormHost.Alias = selected.Alias
		m.FormHost.Name = selected.Name
		m.FormHost.User = selected.User
		m.FormHost.Port = selected.Port
		m.FormHost.IdentityFile = selected.IdentityFile
		m.FormDestFile = selected.SourceFile

		// Pre-populate advanced fields
		m.FormProxyJump = selected.Properties["ProxyJump"]
		if agent, ok := selected.Properties["ForwardAgent"]; ok {
			m.FormForwardAgent = agent
		}
	}

	// Build destination options for file selector
	var fileOptions []huh.Option[string]
	for _, f := range m.Manager.FileOrder {
		fileOptions = append(fileOptions, huh.NewOption(filepath.Base(f), f))
	}

	// Create dynamic steps
	var groups []*huh.Group

	// Step 1: If creating and "All" tab is active, show file selection
	if m.FormAction == actionAdd && (defaultFile == tabAll || defaultFile == "") {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Destination SSH Config File").
				Description("Choose the file where this connection block will be saved").
				Options(fileOptions...).
				Value(&m.FormDestFile),
		))
	}

	// Step 2: Core Connection Properties
	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Alias / Connection Name").
			Description("What you will type to connect (e.g. prod-web-01)").
			Placeholder("my-server").
			Value(&m.FormHost.Alias).
			Validate(validation.ValidateAlias),
		huh.NewInput().
			Title("Server Address / HostName").
			Description("Domain or IP address of the target server").
			Placeholder("10.200.1.45").
			Value(&m.FormHost.Name),
		huh.NewInput().
			Title("Username").
			Description("SSH login user").
			Placeholder("deploy").
			Value(&m.FormHost.User),
		huh.NewInput().
			Title("Port").
			Description("SSH target port").
			Placeholder("22").
			Value(&m.FormHost.Port),
	))

	// Step 3: Advanced Options (IdentityFile, ProxyJump, ForwardAgent)
	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Identity File Path").
			Description("Private key path (e.g. ~/.ssh/id_rsa)").
			Placeholder("~/.ssh/id_rsa").
			Value(&m.FormHost.IdentityFile),
		huh.NewInput().
			Title("Proxy Jump Gateway").
			Description("Intermediate jump host alias if routing through a bastion").
			Placeholder("bastion-host").
			Value(&m.FormProxyJump),
		huh.NewSelect[string]().
			Title("Forward Agent").
			Description("Allow agent forwarding for chained authentication").
			Options(
				huh.NewOption("No", "no"),
				huh.NewOption("Yes", "yes"),
			).
			Value(&m.FormForwardAgent),
	))

	// Construct the final beautiful form
	form := huh.NewForm(groups...).
		WithTheme(huh.ThemeCharm()).
		WithWidth(60)

	return form
}
