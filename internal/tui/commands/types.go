package commands

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
}
