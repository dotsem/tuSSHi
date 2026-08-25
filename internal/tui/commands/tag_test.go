package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"tusshi/internal/config"
	"tusshi/internal/tui/commands"

	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	alertText string
	errorText string
	reloaded  bool
}

func (m *mockContext) Quit()             {}
func (m *mockContext) OpenHelp()         {}
func (m *mockContext) OpenForm(_ string) {}
func (m *mockContext) SetAlert(text string) {
	m.alertText = text
}
func (m *mockContext) SetError(text string) {
	m.errorText = text
}
func (m *mockContext) Reload()                                  { m.reloaded = true }
func (m *mockContext) GetActiveTab() string                     { return "All" }
func (m *mockContext) SetActiveTab(_ string)                    {}
func (m *mockContext) OpenServiceForm(_ string, _ *config.Host) {}
func (m *mockContext) OpenServiceEdit(_ string)                 {}
func (m *mockContext) DeleteService(_ string)                   {}
func (m *mockContext) OpenServices()                            {}

func TestTagCommand(t *testing.T) {
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

	hosts := mgr.GetHosts()
	assert.Len(t, hosts, 1)

	t.Run("adds new tags to selected host", func(t *testing.T) {
		ctx := &mockContext{}
		action := commands.Tag(mgr, hosts[0], []string{":tag", "aws", "k8s"})
		action(ctx)

		assert.True(t, ctx.reloaded)
		assert.Contains(t, ctx.alertText, "Tagged")

		updatedHosts := mgr.GetHosts()
		assert.Equal(t, []string{"production", "aws", "k8s"}, updatedHosts[0].Tags)
	})

	t.Run("removes tags from selected host", func(t *testing.T) {
		ctx := &mockContext{}
		action := commands.Untag(mgr, hosts[0], []string{":untag", "aws"})
		action(ctx)

		assert.True(t, ctx.reloaded)
		assert.Contains(t, ctx.alertText, "Removed tags")

		updatedHosts := mgr.GetHosts()
		assert.Equal(t, []string{"production", "k8s"}, updatedHosts[0].Tags)
	})

	t.Run("error on missing arguments", func(t *testing.T) {
		ctx := &mockContext{}
		action := commands.Tag(mgr, hosts[0], []string{":tag"})
		action(ctx)

		assert.Contains(t, ctx.errorText, "Usage")
	})
}
