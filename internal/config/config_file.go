package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// AddConfigFile creates a blank SSH configuration file on disk,
// registers it with the primary SSH config via an Include directive,
// and maps it internally for display and interaction.
func (m *Manager) AddConfigFile(targetPath string) error {
	absTarget := expandTilde(targetPath)
	if abs, err := filepath.Abs(absTarget); err == nil {
		absTarget = abs
	}

	if _, exists := m.Configs[absTarget]; exists {
		return fmt.Errorf("config file %q already tracked", targetPath)
	}

	dir := filepath.Dir(absTarget)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// write a simple marker comment to represent a blank config file
	if err := os.WriteFile(absTarget, []byte("# SSH config file created by tusshi\n"), 0600); err != nil {
		return err
	}

	f, err := os.Open(filepath.Clean(absTarget))
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return err
	}

	m.Configs[absTarget] = cfg
	m.FileOrder = append(m.FileOrder, absTarget)

	if absTarget != m.PrimaryPath {
		if err := m.registerInclude(absTarget); err != nil {
			return err
		}
	}

	return nil
}

// RenameConfigFile moves a configuration file to a new path on disk,
// updates all internal tracking indices, redirects child hosts, and updates
// the corresponding Include directive inside the primary config.
func (m *Manager) RenameConfigFile(oldPath, newPath string) error {
	absOld := expandTilde(oldPath)
	if abs, err := filepath.Abs(absOld); err == nil {
		absOld = abs
	}

	absNew := expandTilde(newPath)
	if abs, err := filepath.Abs(absNew); err == nil {
		absNew = abs
	}

	if absOld == m.PrimaryPath || absNew == m.PrimaryPath {
		return fmt.Errorf("cannot rename the primary configuration file")
	}

	cfg, exists := m.Configs[absOld]
	if !exists {
		return fmt.Errorf("config file %q not found", oldPath)
	}

	if _, exists := m.Configs[absNew]; exists {
		return fmt.Errorf("target config file %q already exists", newPath)
	}

	dir := filepath.Dir(absNew)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := os.Rename(absOld, absNew); err != nil {
		return err
	}

	delete(m.Configs, absOld)
	m.Configs[absNew] = cfg

	for i, f := range m.FileOrder {
		if f == absOld {
			m.FileOrder[i] = absNew
			break
		}
	}

	if err := m.updateInclude(absOld, absNew); err != nil {
		return err
	}

	return nil
}

// DeleteConfigFile removes a configuration file from disk and deletes
// its associated Include directive inside the primary configuration.
// It fails if any host connections are still defined inside the file.
func (m *Manager) DeleteConfigFile(targetPath string) error {
	absTarget := expandTilde(targetPath)
	if abs, err := filepath.Abs(absTarget); err == nil {
		absTarget = abs
	}

	if absTarget == m.PrimaryPath {
		return fmt.Errorf("cannot delete the primary configuration file")
	}

	cfg, exists := m.Configs[absTarget]
	if !exists {
		return fmt.Errorf("config file %q not tracked", targetPath)
	}

	for _, astHost := range cfg.Hosts {
		val := reflect.ValueOf(astHost)
		isImplicit := false
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			elem := val.Elem()
			implicitField := elem.FieldByName("implicit")
			if implicitField.IsValid() && implicitField.Kind() == reflect.Bool && implicitField.Bool() {
				isImplicit = true
			}
		}

		if isImplicit {
			continue
		}

		for _, pat := range astHost.Patterns {
			alias := pat.String()
			if alias != "" && !strings.ContainsAny(alias, "*?") {
				return fmt.Errorf("cannot delete config file %q: connections are still present", filepath.Base(targetPath))
			}
		}
	}

	if err := os.Remove(absTarget); err != nil && !os.IsNotExist(err) {
		return err
	}

	delete(m.Configs, absTarget)

	for i, f := range m.FileOrder {
		if f == absTarget {
			m.FileOrder = append(m.FileOrder[:i], m.FileOrder[i+1:]...)
			break
		}
	}

	if err := m.removeInclude(absTarget); err != nil {
		return err
	}

	return nil
}

// createIncludeNode constructs an AST node representing an Include directive.
func (m *Manager) createIncludeNode(path string) (ssh_config.Node, error) {
	relPath := path
	if primaryDir := filepath.Dir(m.PrimaryPath); strings.HasPrefix(path, primaryDir) {
		if rel, err := filepath.Rel(primaryDir, path); err == nil {
			relPath = rel
		}
	}

	decoded, err := ssh_config.Decode(strings.NewReader("Include " + relPath + "\n"))
	if err != nil || len(decoded.Hosts) == 0 || len(decoded.Hosts[0].Nodes) == 0 {
		return nil, fmt.Errorf("failed to generate include AST node")
	}
	return decoded.Hosts[0].Nodes[0], nil
}

// updateInclude updates an existing Include directive inside the primary config.
func (m *Manager) updateInclude(oldPath, newPath string) error {
	primaryCfg, exists := m.Configs[m.PrimaryPath]
	if !exists {
		return nil
	}

	oldRel := oldPath
	if primaryDir := filepath.Dir(m.PrimaryPath); strings.HasPrefix(oldPath, primaryDir) {
		if rel, err := filepath.Rel(primaryDir, oldPath); err == nil {
			oldRel = rel
		}
	}

	newNode, err := m.createIncludeNode(newPath)
	if err != nil {
		return err
	}

	updated := false
	for _, astHost := range primaryCfg.Hosts {
		for i, node := range astHost.Nodes {
			if incl, ok := node.(*ssh_config.Include); ok {
				inclStr := strings.TrimSpace(incl.String())
				if idx := strings.Index(inclStr, "#"); idx != -1 {
					inclStr = strings.TrimSpace(inclStr[:idx])
				}
				inclStr = strings.TrimPrefix(inclStr, "Include")
				inclStr = strings.TrimSpace(inclStr)
				inclStr = strings.TrimPrefix(inclStr, "=")
				inclStr = strings.TrimSpace(inclStr)

				if inclStr == oldRel || inclStr == oldPath {
					astHost.Nodes[i] = newNode
					updated = true
					break
				}
			}
		}
		if updated {
			break
		}
	}

	if !updated {
		return m.registerInclude(newPath)
	}

	return m.SaveFile(m.PrimaryPath)
}

// removeInclude locates and removes a tracked Include directive in the primary config.
func (m *Manager) removeInclude(targetPath string) error {
	primaryCfg, exists := m.Configs[m.PrimaryPath]
	if !exists {
		return nil
	}

	targetRel := targetPath
	if primaryDir := filepath.Dir(m.PrimaryPath); strings.HasPrefix(targetPath, primaryDir) {
		if rel, err := filepath.Rel(primaryDir, targetPath); err == nil {
			targetRel = rel
		}
	}

	updated := false
	for _, astHost := range primaryCfg.Hosts {
		for i, node := range astHost.Nodes {
			if incl, ok := node.(*ssh_config.Include); ok {
				inclStr := strings.TrimSpace(incl.String())
				if idx := strings.Index(inclStr, "#"); idx != -1 {
					inclStr = strings.TrimSpace(inclStr[:idx])
				}
				inclStr = strings.TrimPrefix(inclStr, "Include")
				inclStr = strings.TrimSpace(inclStr)
				inclStr = strings.TrimPrefix(inclStr, "=")
				inclStr = strings.TrimSpace(inclStr)

				if inclStr == targetRel || inclStr == targetPath {
					astHost.Nodes = append(astHost.Nodes[:i], astHost.Nodes[i+1:]...)
					updated = true
					break
				}
			}
		}
		if updated {
			break
		}
	}

	if updated {
		return m.SaveFile(m.PrimaryPath)
	}

	return nil
}
