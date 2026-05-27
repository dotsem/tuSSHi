package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnsureFirstRunBackup tests the creation of a pre-tusshi backup on the first run.
func TestEnsureFirstRunBackup(t *testing.T) {
	t.Run("successful backup of primary and includes", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryPath := filepath.Join(tmpDir, "config")
		includePath := filepath.Join(tmpDir, "include_config")

		err := os.WriteFile(primaryPath, []byte("Host test\n  HostName 1.2.3.4\n"), 0600)
		assert.NoError(t, err)

		err = os.WriteFile(includePath, []byte("Host include\n  HostName 5.6.7.8\n"), 0600)
		assert.NoError(t, err)

		mgr := NewManager(primaryPath)
		mgr.FileOrder = []string{primaryPath, includePath}

		err = mgr.EnsureFirstRunBackup()
		assert.NoError(t, err)

		backupDir := filepath.Join(tmpDir, "pre-tusshi")
		assert.DirExists(t, backupDir)

		// check folder permissions (must be 0700 on Unix)
		info, err := os.Stat(backupDir)
		assert.NoError(t, err)
		assert.Equal(t, os.ModeDir|0700, info.Mode().Perm()|os.ModeDir)

		metaPath := filepath.Join(backupDir, "metadata.json")
		assert.FileExists(t, metaPath)

		// #nosec G304 - metaPath is a dynamic test path inside a temporary directory
		metaBytes, err := os.ReadFile(metaPath)
		assert.NoError(t, err)

		var meta BackupMetadata
		err = json.Unmarshal(metaBytes, &meta)
		assert.NoError(t, err)

		assert.NotEmpty(t, meta.Timestamp)
		assert.Len(t, meta.Files, 2)

		// verify metadata content and permissions of backed up files
		for _, f := range meta.Files {
			backedUpPath := filepath.Join(backupDir, f.BackupName)
			assert.FileExists(t, backedUpPath)

			fileInfo, err := os.Stat(backedUpPath)
			assert.NoError(t, err)
			assert.Equal(t, os.FileMode(0600), fileInfo.Mode().Perm())

			// #nosec G304 - testing with dynamic test file path
			originalContent, err := os.ReadFile(f.OriginalPath)
			assert.NoError(t, err)
			// #nosec G304 - testing with dynamic test file path
			backupContent, err := os.ReadFile(backedUpPath)
			assert.NoError(t, err)
			assert.Equal(t, originalContent, backupContent)
		}
	})

	t.Run("consecutive run prevents redundant backup", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryPath := filepath.Join(tmpDir, "config")

		err := os.WriteFile(primaryPath, []byte("Host first\n"), 0600)
		assert.NoError(t, err)

		mgr := NewManager(primaryPath)
		mgr.FileOrder = []string{primaryPath}

		err = mgr.EnsureFirstRunBackup()
		assert.NoError(t, err)

		backupDir := filepath.Join(tmpDir, "pre-tusshi")
		metaPath := filepath.Join(backupDir, "metadata.json")
		assert.FileExists(t, metaPath)

		// modify primary config to see if a second run overwrites it
		err = os.WriteFile(primaryPath, []byte("Host modified\n"), 0600)
		assert.NoError(t, err)

		err = mgr.EnsureFirstRunBackup()
		assert.NoError(t, err)

		// original backed up file should still have the original content
		backedUpPath := filepath.Join(backupDir, "config_0")
		// #nosec G304 - testing with dynamic test file path
		backupContent, err := os.ReadFile(backedUpPath)
		assert.NoError(t, err)
		assert.Equal(t, []byte("Host first\n"), backupContent)
	})

	t.Run("missing files are gracefully skipped", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryPath := filepath.Join(tmpDir, "config")
		missingPath := filepath.Join(tmpDir, "missing")

		err := os.WriteFile(primaryPath, []byte("Host primary\n"), 0600)
		assert.NoError(t, err)

		mgr := NewManager(primaryPath)
		mgr.FileOrder = []string{primaryPath, missingPath}

		err = mgr.EnsureFirstRunBackup()
		assert.NoError(t, err)

		backupDir := filepath.Join(tmpDir, "pre-tusshi")
		metaPath := filepath.Join(backupDir, "metadata.json")
		assert.FileExists(t, metaPath)

		// #nosec G304 - testing with dynamic test file path
		metaBytes, err := os.ReadFile(metaPath)
		assert.NoError(t, err)

		var meta BackupMetadata
		err = json.Unmarshal(metaBytes, &meta)
		assert.NoError(t, err)

		assert.Len(t, meta.Files, 1)
		assert.Equal(t, primaryPath, meta.Files[0].OriginalPath)
	})
}
