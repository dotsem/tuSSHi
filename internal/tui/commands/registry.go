// Package commands defines the interactive command execution handlers for the TUSSHI TUI.
package commands

import (
	"fmt"
	"strings"

	"tusshi/internal/tui/components"
)

// Command defines the behavioral contract for executable TUI commands.
type Command interface {
	Keys() []string
	Description() string
	Execute(ctx Context, args []string)
}

// Registry manages the set of registered commands.
type Registry struct {
	commands []Command
}

// NewRegistry creates a new empty command registry.
func NewRegistry() *Registry {
	return &Registry{
		commands: make([]Command, 0),
	}
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd Command) {
	r.commands = append(r.commands, cmd)
}

// HelpOptions dynamically builds help options from all registered commands.
func (r *Registry) HelpOptions() []components.HelpOption {
	opts := make([]components.HelpOption, 0, len(r.commands))
	for _, cmd := range r.commands {
		shortcut := strings.Join(cmd.Keys(), ", ")
		opts = append(opts, components.HelpOption{
			Shortcut:    shortcut,
			Description: cmd.Description(),
		})
	}
	return opts
}

// Dispatch matches raw input against registered commands and executes the handler.
func (r *Registry) Dispatch(raw string, ctx Context) error {
	parts := strings.Fields(strings.TrimPrefix(raw, ":"))
	if len(parts) == 0 {
		return nil
	}

	key := parts[0]
	for _, cmd := range r.commands {
		for _, k := range cmd.Keys() {
			if key == strings.TrimSpace(k) {
				cmd.Execute(ctx, parts)
				return nil
			}
		}
	}

	return fmt.Errorf("unknown command: %s", key)
}

// DefaultRegistry is the central command registry instance.
var DefaultRegistry = NewRegistry()

// Register adds a command to the default package registry.
func Register(cmd Command) {
	DefaultRegistry.Register(cmd)
}

// Dispatch runs raw command string on the default package registry.
func Dispatch(raw string, ctx Context) error {
	return DefaultRegistry.Dispatch(raw, ctx)
}

// HelpOptions builds help options dynamically from the default registry.
func HelpOptions() []components.HelpOption {
	return DefaultRegistry.HelpOptions()
}
