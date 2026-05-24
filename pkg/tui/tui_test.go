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

	t.Run("initialization", func(t *testing.T) {
		assert.Equal(t, ModeNormal, m.Mode)
		assert.Len(t, m.Filtered, 2)
		assert.Equal(t, 0, m.SelectedIndex)
	})

	t.Run("down navigation", func(t *testing.T) {
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		assert.Nil(t, cmd)
		m = updatedModel.(*Model)
		assert.Equal(t, 1, m.SelectedIndex)
	})

	t.Run("up navigation", func(t *testing.T) {
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		assert.Nil(t, cmd)
		m = updatedModel.(*Model)
		assert.Equal(t, 0, m.SelectedIndex)
	})

	t.Run("command mode toggle", func(t *testing.T) {
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
		assert.NotNil(t, cmd)
		m = updatedModel.(*Model)
		assert.Equal(t, ModeCommand, m.Mode)
		assert.True(t, m.CommandInput.Focused())
	})
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

	t.Run("execute quit", func(t *testing.T) {
		_, cmd := m.executeCommand("quit")
		assert.NotNil(t, cmd)
		msg := cmd()
		_, isQuit := msg.(tea.QuitMsg)
		assert.True(t, isQuit)
	})
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

	subPath := filepath.Join(tmpDir, "sub-config-tui")
	renamedPath := filepath.Join(tmpDir, "renamed-config-tui")

	t.Run("add-config", func(t *testing.T) {
		_, _ = m.executeCommand("add-config " + subPath)
		assert.FileExists(t, subPath)
		assert.Contains(t, m.Manager.FileOrder, subPath)
		assert.Equal(t, subPath, m.ActiveTab)
	})

	t.Run("rename-config", func(t *testing.T) {
		_, _ = m.executeCommand("rename-config " + renamedPath)
		assert.FileExists(t, renamedPath)
		assert.NoFileExists(t, subPath)
		assert.Contains(t, m.Manager.FileOrder, renamedPath)
		assert.Equal(t, renamedPath, m.ActiveTab)
	})

	t.Run("delete-config", func(t *testing.T) {
		_, _ = m.executeCommand("delete-config")
		assert.NoFileExists(t, renamedPath)
		assert.NotContains(t, m.Manager.FileOrder, renamedPath)
		assert.Equal(t, "All", m.ActiveTab)
	})
}

// TestTUIPrimaryConfigHiding verifies the primary config tab visibility based on connections.
func TestTUIPrimaryConfigHiding(t *testing.T) {
	t.Run("config file is not empty - show tab", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryPath := filepath.Join(tmpDir, "config")

		content := "Host myhost\n    HostName 10.0.0.1\n"
		err := os.WriteFile(primaryPath, []byte(content), 0600)
		assert.NoError(t, err)

		mgr := config.NewManager(primaryPath)
		err = mgr.Load()
		assert.NoError(t, err)

		m := NewModel(mgr)
		assert.Contains(t, m.Tabs, primaryPath)
	})

	t.Run("config file is empty - hide tab", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryPath := filepath.Join(tmpDir, "config")

		err := os.WriteFile(primaryPath, []byte("# Empty config\n"), 0600)
		assert.NoError(t, err)

		mgr := config.NewManager(primaryPath)
		err = mgr.Load()
		assert.NoError(t, err)

		m := NewModel(mgr)
		assert.NotContains(t, m.Tabs, primaryPath)
	})
}

// TestTUIPing verifies that the background ping structures work correctly.
func TestTUIPing(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	content := `
Host srv-1
    HostName localhost
    Port 9999
`
	err := os.WriteFile(primaryPath, []byte(content), 0600)
	assert.NoError(t, err)

	mgr := config.NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	m := NewModel(mgr)

	t.Run("PingAll initialization", func(t *testing.T) {
		cmd := m.PingAll()
		assert.NotNil(t, cmd)
		assert.Contains(t, m.PingResults, "srv-1")
		assert.True(t, m.PingResults["srv-1"].Pending)
	})

	t.Run("PingResultMsg processing", func(t *testing.T) {
		updatedModel, cmd := m.Update(PingResultMsg{
			Alias:   "srv-1",
			Online:  true,
			Latency: 12.34,
		})
		assert.Nil(t, cmd)
		m = updatedModel.(*Model)
		assert.False(t, m.PingResults["srv-1"].Pending)
		assert.True(t, m.PingResults["srv-1"].Online)
		assert.Equal(t, 12.34, m.PingResults["srv-1"].Latency)
	})

	t.Run("render status cell", func(t *testing.T) {
		cell := m.renderStatusCell("srv-1", false, 15)
		assert.Contains(t, cell, "Online")
		assert.Contains(t, cell, "12ms")
	})
}
