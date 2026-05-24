package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestManagerBasicLoad tests parsing, wildcard inheritance, and host mapping.
func TestManagerBasicLoad(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	configContent := `
# Global Wildcard definition
Host *
    User default_user
    Port 2222
    ForwardAgent yes

# Production web host
Host prod-web-01
    HostName 10.200.1.45
    User deploy
    IdentityFile ~/.ssh/keys/work_rsa

# Database host without separate User/Port (inherits from *)
Host 10.200.1.46
`
	err := os.WriteFile(primaryPath, []byte(configContent), 0600)
	assert.NoError(t, err)

	mgr := NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	hosts := mgr.GetHosts()
	// Total hosts: Host *, prod-web-01, 10.200.1.46
	assert.Len(t, hosts, 3)

	var wildcardHost, prodHost, dbHost *Host
	for _, h := range hosts {
		switch h.Alias {
		case "*":
			wildcardHost = h
		case "prod-web-01":
			prodHost = h
		case "10.200.1.46":
			dbHost = h
		}
	}

	assert.NotNil(t, wildcardHost)
	assert.True(t, wildcardHost.IsWildcard)
	assert.Equal(t, "default_user", wildcardHost.User)

	// Verify prodHost details and explicit key mapping
	assert.NotNil(t, prodHost)
	assert.False(t, prodHost.IsWildcard)
	assert.Equal(t, "10.200.1.45", prodHost.Name)
	assert.Equal(t, "deploy", prodHost.User)
	assert.Equal(t, "~/.ssh/keys/work_rsa", prodHost.IdentityFile)
	// Verify prodHost inherited port 2222 from wildcard Host *
	assert.Equal(t, "2222", prodHost.ResolvedProperties["Port"])
	assert.Equal(t, "yes", prodHost.ResolvedProperties["ForwardAgent"])

	// Verify dbHost does not have alias (its alias is the IP itself)
	assert.NotNil(t, dbHost)
	assert.Equal(t, "10.200.1.46", dbHost.Alias)
	// dbHost should inherit User, Port, and ForwardAgent from wildcard
	assert.Equal(t, "default_user", dbHost.User)
	assert.Equal(t, "2222", dbHost.Port)
	assert.Equal(t, "yes", dbHost.ResolvedProperties["ForwardAgent"])
}

// TestManagerIncludes tests glob inclusion and recursive parsing.
func TestManagerIncludes(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")
	includeDir := filepath.Join(tmpDir, "config.d")
	err := os.Mkdir(includeDir, 0700)
	assert.NoError(t, err)

	mainContent := "Include config.d/*\n"
	err = os.WriteFile(primaryPath, []byte(mainContent), 0600)
	assert.NoError(t, err)

	workContent := `
Host work-host
    HostName 10.0.0.10
    User work_user
`
	err = os.WriteFile(filepath.Join(includeDir, "work"), []byte(workContent), 0600)
	assert.NoError(t, err)

	mgr := NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	hosts := mgr.GetHosts()
	assert.Len(t, hosts, 1)
	assert.Equal(t, "work-host", hosts[0].Alias)
	assert.Equal(t, "10.0.0.10", hosts[0].Name)
	assert.Equal(t, "work_user", hosts[0].User)
}

// TestManagerCRUD tests adding, updating, and deleting hosts.
func TestManagerCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")

	initialContent := `
# My connections
Host my-host
    HostName 127.0.0.1
`
	err := os.WriteFile(primaryPath, []byte(initialContent), 0600)
	assert.NoError(t, err)

	mgr := NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	// 1. Add host
	newHost := &Host{
		Alias: "added-host",
		Name:  "192.168.1.10",
		User:  "admin",
		Port:  "22",
		Properties: map[string]string{
			"ForwardAgent": "yes",
		},
	}
	err = mgr.AddHost(primaryPath, newHost)
	assert.NoError(t, err)

	// Reload to verify write
	mgr2 := NewManager(primaryPath)
	err = mgr2.Load()
	assert.NoError(t, err)
	hosts := mgr2.GetHosts()
	assert.Len(t, hosts, 2)

	var addedHost *Host
	for _, h := range hosts {
		if h.Alias == "added-host" {
			addedHost = h
		}
	}
	assert.NotNil(t, addedHost)
	assert.Equal(t, "192.168.1.10", addedHost.Name)
	assert.Equal(t, "admin", addedHost.User)
	assert.Equal(t, "22", addedHost.Port)
	assert.Equal(t, "yes", addedHost.Properties["ForwardAgent"])

	// 2. Update host
	updatedHost := &Host{
		Alias: "added-host-new",
		Name:  "192.168.1.15",
		User:  "root",
		Port:  "222",
		Properties: map[string]string{
			"ForwardAgent": "no",
			"ProxyJump":    "jump-box",
		},
	}
	err = mgr2.UpdateHost("added-host", updatedHost)
	assert.NoError(t, err)

	// Reload to verify update
	mgr3 := NewManager(primaryPath)
	err = mgr3.Load()
	assert.NoError(t, err)
	hosts3 := mgr3.GetHosts()
	assert.Len(t, hosts3, 2)

	var foundUpdated *Host
	for _, h := range hosts3 {
		if h.Alias == "added-host-new" {
			foundUpdated = h
		}
	}
	assert.NotNil(t, foundUpdated)
	assert.Equal(t, "192.168.1.15", foundUpdated.Name)
	assert.Equal(t, "root", foundUpdated.User)
	assert.Equal(t, "222", foundUpdated.Port)
	assert.Equal(t, "no", foundUpdated.Properties["ForwardAgent"])
	assert.Equal(t, "jump-box", foundUpdated.Properties["ProxyJump"])

	// 3. Delete host
	err = mgr3.DeleteHost("added-host-new")
	assert.NoError(t, err)

	// Reload to verify delete
	mgr4 := NewManager(primaryPath)
	err = mgr4.Load()
	assert.NoError(t, err)
	assert.Len(t, mgr4.GetHosts(), 1)
}

// TestManagerConfigFileCRUD tests creating, renaming, and deleting config files.
func TestManagerConfigFileCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	primaryPath := filepath.Join(tmpDir, "config")
	err := os.WriteFile(primaryPath, []byte("# Primary\n"), 0600)
	assert.NoError(t, err)

	mgr := NewManager(primaryPath)
	err = mgr.Load()
	assert.NoError(t, err)

	subPath := filepath.Join(tmpDir, "sub-config")

	// 1. Add Config File
	err = mgr.AddConfigFile(subPath)
	assert.NoError(t, err)
	assert.FileExists(t, subPath)
	assert.Contains(t, mgr.FileOrder, subPath)

	// Check if Include directive is added in primary
	// #nosec G304
	primaryContent, err := os.ReadFile(primaryPath)
	assert.NoError(t, err)
	relSub, err := filepath.Rel(filepath.Dir(primaryPath), subPath)
	assert.NoError(t, err)
	assert.Contains(t, string(primaryContent), "Include "+relSub)

	// 2. Rename Config File
	renamedPath := filepath.Join(tmpDir, "renamed-config")
	err = mgr.RenameConfigFile(subPath, renamedPath)
	assert.NoError(t, err)
	assert.FileExists(t, renamedPath)
	assert.NoFileExists(t, subPath)
	assert.Contains(t, mgr.FileOrder, renamedPath)
	assert.NotContains(t, mgr.FileOrder, subPath)

	// Check if Include is updated in primary
	// #nosec G304
	primaryContent2, err := os.ReadFile(primaryPath)
	assert.NoError(t, err)
	relRenamed, err := filepath.Rel(filepath.Dir(primaryPath), renamedPath)
	assert.NoError(t, err)
	assert.Contains(t, string(primaryContent2), "Include "+relRenamed)
	assert.NotContains(t, string(primaryContent2), "Include "+relSub)

	// 3. Prevent deleting if connections are present
	h := &Host{
		Alias: "test-host",
		Name:  "127.0.0.1",
	}
	err = mgr.AddHost(renamedPath, h)
	assert.NoError(t, err)

	err = mgr.DeleteConfigFile(renamedPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connections are still present")

	// Delete host first
	err = mgr.DeleteHost("test-host")
	assert.NoError(t, err)

	// 4. Delete Config File successfully when no connections are present
	err = mgr.DeleteConfigFile(renamedPath)
	assert.NoError(t, err)
	assert.NoFileExists(t, renamedPath)
	assert.NotContains(t, mgr.FileOrder, renamedPath)

	// Check if Include is removed from primary
	// #nosec G304
	primaryContent3, err := os.ReadFile(primaryPath)
	assert.NoError(t, err)
	assert.NotContains(t, string(primaryContent3), "Include "+relRenamed)
}
