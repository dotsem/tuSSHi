package tui

import (
	"os"
	"path/filepath"
	"testing"

	"tusshi/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestTUIBasicSetup tests TUI model initialization and basic key navigations.
func TestTUIBasicSetup(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	content := `
Host srv-1
    HostName 10.0.0.1
    User deploy
Host srv-2
    HostName 10.0.0.2
    User root
`
	err := os.WriteFile(primaryPath, []byte(content), 0600)
	assert.NoError(t, err)

	mgr := config.NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	m := NewModel(mgr)
	assert.Equal(t, ModeNormal, m.Mode)
	assert.Len(t, m.Filtered, 2)
	assert.Equal(t, 0, m.SelectedIndex)

	// Test j key (down navigation)
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Nil(t, cmd)
	m = updatedModel.(*Model)
	assert.Equal(t, 1, m.SelectedIndex)

	// Test k key (up navigation)
	updatedModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Nil(t, cmd)
	m = updatedModel.(*Model)
	assert.Equal(t, 0, m.SelectedIndex)

	// Test : key (command mode toggle)
	updatedModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	assert.NotNil(t, cmd)
	m = updatedModel.(*Model)
	assert.Equal(t, ModeCommand, m.Mode)
	assert.True(t, m.CommandInput.Focused())
}

// TestTUICommandExec tests command mode execution parsing.
func TestTUICommandExec(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	err := os.WriteFile(primaryPath, []byte("Host myhost\n"), 0600)
	assert.NoError(t, err)

	mgr := config.NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	m := NewModel(mgr)

	// Test :q execution
	_, cmd := m.executeCommand("quit")
	// Verification: A quit command returns a tea.Quit cmd
	assert.NotNil(t, cmd)
	// Execute the cmd function and verify it returns a tea.QuitMsg
	msg := cmd()
	_, isQuit := msg.(tea.QuitMsg)
	assert.True(t, isQuit)
}

// TestTUIConfigCommands tests executing config file management commands in the TUI.
func TestTUIConfigCommands(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	err := os.WriteFile(primaryPath, []byte("# Primary\n"), 0600)
	assert.NoError(t, err)

	mgr := config.NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	m := NewModel(mgr)

	// 1. Test add-config command execution
	subPath := filepath.Join(tmpDir, "sub-config-tui")
	_, _ = m.executeCommand("add-config " + subPath)
	assert.FileExists(t, subPath)
	assert.Contains(t, m.Manager.FileOrder, subPath)
	assert.Equal(t, subPath, m.ActiveTab)

	// 2. Test rename-config command execution
	renamedPath := filepath.Join(tmpDir, "renamed-config-tui")
	_, _ = m.executeCommand("rename-config " + renamedPath)
	assert.FileExists(t, renamedPath)
	assert.NoFileExists(t, subPath)
	assert.Contains(t, m.Manager.FileOrder, renamedPath)
	assert.Equal(t, renamedPath, m.ActiveTab)

	// 3. Test delete-config command execution
	_, _ = m.executeCommand("delete-config")
	assert.NoFileExists(t, renamedPath)
	assert.NotContains(t, m.Manager.FileOrder, renamedPath)
	assert.Equal(t, "All", m.ActiveTab)
}

