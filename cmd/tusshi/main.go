// Package main is the entry point of the tusshi application, initializing
// the configuration manager and launching the Bubble Tea interactive TUI.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"tusshi/pkg/config"
	"tusshi/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	mgr := config.NewManager(configPath)
	model := tui.NewModel(mgr)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR running TUSSHI: %v\n", err)
		os.Exit(1)
	}
}
