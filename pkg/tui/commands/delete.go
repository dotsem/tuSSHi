package commands

import (
	"fmt"

	"tusshi/pkg/config"
)

// Delete returns a function that executes the deletion of a host connection.
func Delete(mgr *config.Manager, selectedHost *config.Host) func(Context) {
	return func(ctx Context) {
		if err := mgr.DeleteHost(selectedHost.Alias); err != nil {
			ctx.SetError("Delete error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Deleted connection %q.", selectedHost.Alias))
		}
		ctx.Reload()
	}
}
