// Package commands defines the interactive command execution handlers for the TUSSHI TUI.
package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"tusshi/internal/validation"
)

const tabAll = "All"

type addConfigCmd struct{}

func (c *addConfigCmd) Keys() []string      { return []string{"addconf", "add-config"} }
func (c *addConfigCmd) Description() string { return "Add a new config file" }
func (c *addConfigCmd) Execute(ctx Context, parts []string) {
	if len(parts) < 2 {
		ctx.SetError("Usage: :add-config <filename>")
		return
	}

	mgr := ctx.GetManager()
	arg := parts[1]
	var targetPath string
	if filepath.IsAbs(arg) || strings.HasPrefix(arg, "~/") {
		targetPath = arg
	} else {
		targetPath = filepath.Join(filepath.Dir(mgr.PrimaryPath), arg)
	}

	if err := validation.ValidateConfigName(filepath.Base(targetPath)); err != nil {
		ctx.SetError("Invalid config name: " + err.Error())
		return
	}

	if err := mgr.AddConfigFile(targetPath); err != nil {
		ctx.SetError("Add config error: " + err.Error())
	} else {
		ctx.SetAlert(fmt.Sprintf("Created config file %q.", filepath.Base(targetPath)))
		ctx.SetActiveTab(targetPath)
	}
	ctx.Reload()
}

type renameConfigCmd struct{}

func (c *renameConfigCmd) Keys() []string      { return []string{"mvconf", "rename-config"} }
func (c *renameConfigCmd) Description() string { return "Rename a config file" }
func (c *renameConfigCmd) Execute(ctx Context, parts []string) {
	mgr := ctx.GetManager()
	activeTab := ctx.GetActiveTab()

	var oldName, newName string
	switch {
	case len(parts) == 2:
		if activeTab == tabAll {
			ctx.SetError("Cannot rename from 'All' tab. Usage: :rename-config <old-name> <new-name>")
			return
		}
		oldName = activeTab
		newName = parts[1]
	case len(parts) >= 3:
		oldName = parts[1]
		newName = parts[2]
	default:
		if activeTab == tabAll {
			ctx.SetError("Usage: :rename-config <old-name> <new-name>")
		} else {
			ctx.SetError("Usage: :rename-config <new-name>")
		}
		return
	}

	oldPath, found := mgr.FindConfigFile(oldName)
	if !found {
		ctx.SetError(fmt.Sprintf("Config file %q not found", oldName))
		return
	}

	var newPath string
	if filepath.IsAbs(newName) || strings.HasPrefix(newName, "~/") {
		newPath = newName
	} else {
		newPath = filepath.Join(filepath.Dir(mgr.PrimaryPath), newName)
	}

	if err := validation.ValidateConfigName(filepath.Base(newPath)); err != nil {
		ctx.SetError("Invalid config name: " + err.Error())
		return
	}

	if err := mgr.RenameConfigFile(oldPath, newPath); err != nil {
		ctx.SetError("Rename config error: " + err.Error())
	} else {
		ctx.SetAlert(fmt.Sprintf("Renamed config file to %q.", filepath.Base(newPath)))
		if activeTab == oldPath {
			ctx.SetActiveTab(newPath)
		}
	}
	ctx.Reload()
}

type deleteConfigCmd struct{}

func (c *deleteConfigCmd) Keys() []string      { return []string{"rmconf", "delete-config"} }
func (c *deleteConfigCmd) Description() string { return "Delete empty config file" }
func (c *deleteConfigCmd) Execute(ctx Context, parts []string) {
	mgr := ctx.GetManager()
	activeTab := ctx.GetActiveTab()

	var targetName string
	if len(parts) >= 2 {
		targetName = parts[1]
	} else {
		if activeTab == tabAll {
			ctx.SetError("Usage: :delete-config <filename> (or switch to a tab and run :delete-config)")
			return
		}
		targetName = activeTab
	}

	targetPath, found := mgr.FindConfigFile(targetName)
	if !found {
		ctx.SetError(fmt.Sprintf("Config file %q not found", targetName))
		return
	}

	if err := mgr.DeleteConfigFile(targetPath); err != nil {
		ctx.SetError("Delete config error: " + err.Error())
	} else {
		ctx.SetAlert(fmt.Sprintf("Deleted config file %q.", filepath.Base(targetPath)))
		if activeTab == targetPath {
			ctx.SetActiveTab(tabAll)
		}
	}
	ctx.Reload()
}

func init() {
	Register(&addConfigCmd{})
	Register(&renameConfigCmd{})
	Register(&deleteConfigCmd{})
}
