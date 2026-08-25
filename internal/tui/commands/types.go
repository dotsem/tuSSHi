package commands

import "tusshi/internal/config"

// Context defines the behavioral interface for executing TUI commands.
type Context interface {
	Quit()
	OpenHelp()
	OpenForm(action string)
	SetAlert(text string)
	SetError(text string)
	Reload()
	GetActiveTab() string
	SetActiveTab(tab string)
	OpenServiceForm(action string, targetHost *config.Host)
	OpenServiceEdit(alias string)
	DeleteService(alias string)
	OpenServices()
}
