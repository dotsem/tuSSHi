package tui

import (
	"os"
	"path/filepath"
	"testing"

	"tusshi/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTagCommands(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	content := `
Host web-server
    # tags: production
    HostName 10.0.0.1
`
	err := os.WriteFile(primaryPath, []byte(content), 0600)
	assert.NoError(t, err)

	mgr := config.NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	m := NewModel(mgr)

	t.Run("execute :tag command", func(t *testing.T) {
		m.SelectedIndex = 0
		m.executeCommand(":tag aws k8s")
		hosts := m.Manager.GetHosts()
		assert.Len(t, hosts, 1)
		assert.Equal(t, []string{"production", "aws", "k8s"}, hosts[0].Tags)
		assert.Contains(t, m.AlertText, "Tagged")
	})

	t.Run("execute :untag command", func(t *testing.T) {
		m.SelectedIndex = 0
		m.executeCommand(":untag aws")
		hosts := m.Manager.GetHosts()
		assert.Len(t, hosts, 1)
		assert.Equal(t, []string{"production", "k8s"}, hosts[0].Tags)
		assert.Contains(t, m.AlertText, "Removed tags")
	})

	t.Run("execute :tag with explicit alias target", func(t *testing.T) {
		m.executeCommand(":tag web-server staging")
		hosts := m.Manager.GetHosts()
		assert.Len(t, hosts, 1)
		assert.Contains(t, hosts[0].Tags, "staging")
	})
}
