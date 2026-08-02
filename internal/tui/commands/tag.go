package commands

import (
	"fmt"
	"slices"
	"strings"

	"tusshi/internal/config"
)

// Tag appends metadata tags to a target host or the selected host.
func Tag(mgr *config.Manager, selectedHost *config.Host, parts []string) func(Context) {
	return func(ctx Context) {
		if len(parts) < 2 {
			ctx.SetError("Usage: :tag [alias] <tag1> [tag2...]")
			return
		}

		targetHost, tagArgs := resolveTargetHostAndTags(mgr, selectedHost, parts[1:])
		if targetHost == nil {
			ctx.SetError("No connection selected")
			return
		}

		if len(tagArgs) == 0 {
			ctx.SetError("Usage: :tag [alias] <tag1> [tag2...]")
			return
		}

		newTags := targetHost.Tags
		for _, t := range tagArgs {
			clean := config.ExtractTagsFromComment("# tags: " + t)
			for _, ct := range clean {
				if !slices.Contains(newTags, ct) {
					newTags = append(newTags, ct)
				}
			}
		}

		targetHost.Tags = newTags
		if err := mgr.UpdateHost(targetHost.Alias, targetHost); err != nil {
			ctx.SetError("Tag error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Tagged %q with %s.", targetHost.Alias, strings.Join(tagArgs, ", ")))
		}
		ctx.Reload()
	}
}

// Untag removes metadata tags from a target host or the selected host.
func Untag(mgr *config.Manager, selectedHost *config.Host, parts []string) func(Context) {
	return func(ctx Context) {
		if len(parts) < 2 {
			ctx.SetError("Usage: :untag [alias] <tag1> [tag2...]")
			return
		}

		targetHost, tagArgs := resolveTargetHostAndTags(mgr, selectedHost, parts[1:])
		if targetHost == nil {
			ctx.SetError("No connection selected")
			return
		}

		if len(tagArgs) == 0 {
			ctx.SetError("Usage: :untag [alias] <tag1> [tag2...]")
			return
		}

		var tagsToRemove []string
		for _, t := range tagArgs {
			clean := config.ExtractTagsFromComment("# tags: " + t)
			tagsToRemove = append(tagsToRemove, clean...)
		}

		var remaining []string
		for _, t := range targetHost.Tags {
			if !slices.Contains(tagsToRemove, t) {
				remaining = append(remaining, t)
			}
		}

		targetHost.Tags = remaining
		if err := mgr.UpdateHost(targetHost.Alias, targetHost); err != nil {
			ctx.SetError("Untag error: " + err.Error())
		} else {
			ctx.SetAlert(fmt.Sprintf("Removed tags from %q.", targetHost.Alias))
		}
		ctx.Reload()
	}
}

func resolveTargetHostAndTags(mgr *config.Manager, selectedHost *config.Host, args []string) (*config.Host, []string) {
	if len(args) == 0 {
		return selectedHost, nil
	}

	hosts := mgr.GetHosts()
	firstArg := args[0]
	for _, h := range hosts {
		if h.Alias == firstArg {
			return h, args[1:]
		}
	}

	return selectedHost, args
}
