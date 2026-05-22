package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"tusshi/pkg/config"
)

// AddConfig triggers creation of a new configuration file.
func AddConfig(mgr *config.Manager, parts []string) func(Context) {
	return func(ctx Context) {
		if len(parts) < 2 {
			ctx.SetError("Usage: :add-config <filename>")
			return
		}

		arg := parts[1]
		var targetPath string
		if filepath.IsAbs(arg) || strings.HasPrefix(arg, "~/") {
			targetPath = arg
		} else {
			targetPath = filepath.Join(filepath.Dir(mgr.PrimaryPath), arg)
		}

		if err := mgr.AddConfigFile(targetPath); err != nil {
			ctx.SetError("Add config error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Created config file %q.", filepath.Base(targetPath)))
			ctx.SetActiveTab(targetPath)
		}
		ctx.Reload()
	}
}

// RenameConfig triggers renaming of an existing configuration file.
func RenameConfig(mgr *config.Manager, parts []string) func(Context) {
	return func(ctx Context) {
		activeTab := ctx.GetActiveTab()

		var oldName, newName string
		if len(parts) == 2 {
			if activeTab == "All" {
				ctx.SetError("Cannot rename from 'All' tab. Usage: :rename-config <old-name> <new-name>")
				return
			}
			oldName = activeTab
			newName = parts[1]
		} else if len(parts) >= 3 {
			oldName = parts[1]
			newName = parts[2]
		} else {
			if activeTab == "All" {
				ctx.SetError("Usage: :rename-config <old-name> <new-name>")
			} else {
				ctx.SetError("Usage: :rename-config <new-name>")
			}
			return
		}

		var oldPath string
		for _, file := range mgr.FileOrder {
			if file == oldName || filepath.Base(file) == oldName || strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) == oldName {
				oldPath = file
				break
			}
		}

		if oldPath == "" {
			ctx.SetError(fmt.Sprintf("Config file %q not found", oldName))
			return
		}

		var newPath string
		if filepath.IsAbs(newName) || strings.HasPrefix(newName, "~/") {
			newPath = newName
		} else {
			newPath = filepath.Join(filepath.Dir(mgr.PrimaryPath), newName)
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
}

// DeleteConfig triggers deletion of a configuration file when it contains no connections.
func DeleteConfig(mgr *config.Manager, parts []string) func(Context) {
	return func(ctx Context) {
		activeTab := ctx.GetActiveTab()

		var targetName string
		if len(parts) >= 2 {
			targetName = parts[1]
		} else {
			if activeTab == "All" {
				ctx.SetError("Usage: :delete-config <filename> (or switch to a tab and run :delete-config)")
				return
			}
			targetName = activeTab
		}

		var targetPath string
		for _, file := range mgr.FileOrder {
			if file == targetName || filepath.Base(file) == targetName || strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) == targetName {
				targetPath = file
				break
			}
		}

		if targetPath == "" {
			ctx.SetError(fmt.Sprintf("Config file %q not found", targetName))
			return
		}

		if err := mgr.DeleteConfigFile(targetPath); err != nil {
			ctx.SetError("Delete config error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Deleted config file %q.", filepath.Base(targetPath)))
			if activeTab == targetPath {
				ctx.SetActiveTab("All")
			}
		}
		ctx.Reload()
	}
}
