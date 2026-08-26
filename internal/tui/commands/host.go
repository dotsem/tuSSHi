package commands

import (
	"fmt"
	"path/filepath"
)

type newHostCmd struct{}

func (c *newHostCmd) Keys() []string      { return []string{"n", "new"} }
func (c *newHostCmd) Description() string { return "Create a new connection" }
func (c *newHostCmd) Execute(ctx Context, _ []string) {
	ctx.OpenForm("add")
}

type editHostCmd struct{}

func (c *editHostCmd) Keys() []string      { return []string{"e", "edit"} }
func (c *editHostCmd) Description() string { return "Edit selected connection" }
func (c *editHostCmd) Execute(ctx Context, _ []string) {
	if ctx.GetSelectedHost() != nil {
		ctx.OpenForm("edit")
	}
}

type deleteHostCmd struct{}

func (c *deleteHostCmd) Keys() []string      { return []string{"d", "rm"} }
func (c *deleteHostCmd) Description() string { return "Delete selected connection" }
func (c *deleteHostCmd) Execute(ctx Context, _ []string) {
	selected := ctx.GetSelectedHost()
	if selected == nil {
		return
	}

	title := "Delete Connection?"
	message := fmt.Sprintf("Are you sure you want to delete host '%s'?", selected.Alias)

	ctx.Confirm(title, message, true, func() {
		mgr := ctx.GetManager()
		if err := mgr.DeleteHost(selected.Alias); err != nil {
			ctx.SetError("Delete error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Deleted connection %q.", selected.Alias))
		}
		ctx.Reload()
	})
}

type moveHostCmd struct{}

func (c *moveHostCmd) Keys() []string      { return []string{"m", "mv"} }
func (c *moveHostCmd) Description() string { return "Move connection to another config file" }
func (c *moveHostCmd) Execute(ctx Context, args []string) {
	selected := ctx.GetSelectedHost()
	if selected == nil {
		return
	}

	if len(args) < 2 {
		ctx.SetError("Usage: :move <target-file-nickname>")
		return
	}

	mgr := ctx.GetManager()
	targetNickname := args[1]
	matchedFile, found := mgr.FindConfigFile(targetNickname)
	if !found {
		matchedFile = filepath.Join(filepath.Dir(mgr.PrimaryPath), targetNickname)
	}

	if err := mgr.MoveHost(selected.Alias, matchedFile); err != nil {
		ctx.SetError("Move error: " + err.Error())
	} else {
		ctx.SetAlert(fmt.Sprintf("Moved %q to %s.", selected.Alias, filepath.Base(matchedFile)))
	}
	ctx.Reload()
}

type pingHostCmd struct{}

func (c *pingHostCmd) Keys() []string      { return []string{"p", "ping"} }
func (c *pingHostCmd) Description() string { return "Ping selected connection" }
func (c *pingHostCmd) Execute(ctx Context, _ []string) {
	selected := ctx.GetSelectedHost()
	if selected != nil {
		ctx.PingHost(selected)
	}
}

type pingAllCmd struct{}

func (c *pingAllCmd) Keys() []string      { return []string{"P", "pingall"} }
func (c *pingAllCmd) Description() string { return "Ping all connections" }
func (c *pingAllCmd) Execute(ctx Context, _ []string) {
	ctx.PingAll()
}

func init() {
	Register(&newHostCmd{})
	Register(&editHostCmd{})
	Register(&deleteHostCmd{})
	Register(&moveHostCmd{})
	Register(&pingHostCmd{})
	Register(&pingAllCmd{})
}
