package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupFileMap holds path association of a backed up config file.
type BackupFileMap struct {
	OriginalPath string `json:"original_path"`
	BackupName   string `json:"backup_name"`
}

// BackupMetadata preserves information about the original source configurations.
type BackupMetadata struct {
	Timestamp string          `json:"timestamp"`
	Files     []BackupFileMap `json:"files"`
}

// EnsureFirstRunBackup checks if a pre-tusshi backup already exists.
// If not, it creates a pristine backup of all configuration files in m.FileOrder
// inside a pre-tusshi directory next to the primary configuration.
// It checks existence based on the metadata file.
func (m *Manager) EnsureFirstRunBackup() error {
	backupDir := filepath.Join(filepath.Dir(m.PrimaryPath), "pre-tusshi")
	metaPath := filepath.Join(backupDir, "metadata.json")

	if _, err := os.Stat(metaPath); err == nil {
		return nil // backup already exists
	}

	var existingFiles []string
	for _, path := range m.FileOrder {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			existingFiles = append(existingFiles, path)
		}
	}

	// enforce secure, user-only read/write access to the backup folder
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	meta := BackupMetadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Files:     make([]BackupFileMap, 0),
	}

	for i, srcPath := range existingFiles {
		baseName := filepath.Base(srcPath)
		backupName := fmt.Sprintf("%s_%d", baseName, i)
		dstPath := filepath.Join(backupDir, backupName)

		if err := backupFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to backup file %q: %w", srcPath, err)
		}

		meta.Files = append(meta.Files, BackupFileMap{
			OriginalPath: srcPath,
			BackupName:   backupName,
		})
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup metadata: %w", err)
	}

	// enforce secure user-only access on the metadata file
	if err := os.WriteFile(metaPath, metaBytes, 0600); err != nil {
		return fmt.Errorf("failed to write backup metadata: %w", err)
	}

	return nil
}

// backupFile copies the content of a file to the destination with 0600 permissions.
func backupFile(src, dst string) error {
	srcFile, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := os.OpenFile(filepath.Clean(dst), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}
