package commands

import (
	"fmt"
	"path/filepath"

	"tusshi/internal/config"
)

// NewHost triggers opening a form for creating a new host connection.
func NewHost() func(Context) {
	return func(ctx Context) {
		ctx.OpenForm("add")
	}
}

// EditHost triggers opening a form for editing the selected host connection.
func EditHost(hasFiltered bool) func(Context) {
	return func(ctx Context) {
		if hasFiltered {
			ctx.OpenForm("edit")
		}
	}
}

// MoveHost handles moving a connection from one configuration file/tab to another.
func MoveHost(mgr *config.Manager, selectedHost *config.Host, parts []string) func(Context) {
	return func(ctx Context) {
		if len(parts) < 2 {
			ctx.SetError("Usage: :move <target-file-nickname>")
			return
		}

		targetNickname := parts[1]
		matchedFile, found := mgr.FindConfigFile(targetNickname)
		if !found {
			matchedFile = filepath.Join(filepath.Dir(mgr.PrimaryPath), targetNickname)
		}

		if err := mgr.MoveHost(selectedHost.Alias, matchedFile); err != nil {
			ctx.SetError("Move error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Moved %q to %s.", selectedHost.Alias, filepath.Base(matchedFile)))
		}
		ctx.Reload()
	}
}

// DeleteHost returns a function that executes the deletion of a host connection.
func DeleteHost(mgr *config.Manager, selectedHost *config.Host) func(Context) {
	return func(ctx Context) {
		if err := mgr.DeleteHost(selectedHost.Alias); err != nil {
			ctx.SetError("Delete error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Deleted connection %q.", selectedHost.Alias))
		}
		ctx.Reload()
	}
}
