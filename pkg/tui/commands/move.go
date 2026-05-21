package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"tusshi/pkg/config"
)

// Move handles moving a connection from one configuration file/tab to another.
func Move(mgr *config.Manager, selectedHost *config.Host, parts []string) func(Context) {
	return func(ctx Context) {
		if len(parts) < 2 {
			ctx.SetError("Usage: :move <target-file-nickname>")
			return
		}

		targetNickname := parts[1]
		var matchedFile string
		for _, file := range mgr.FileOrder {
			if filepath.Base(file) == targetNickname || strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) == targetNickname {
				matchedFile = file
				break
			}
		}

		if matchedFile == "" {
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
