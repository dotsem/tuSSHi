// Package config provides abstractions for parsing, resolving, updating,
// and serializing OpenSSH configuration files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// AddHost appends a new Host connection block to a specific target configuration file
// and serializes the modified AST back to disk.
func (m *Manager) AddHost(targetFile string, h *Host) error {
	absTarget := expandTilde(targetFile)
	if abs, err := filepath.Abs(absTarget); err == nil {
		absTarget = abs
	}

	cfg, exists := m.Configs[absTarget]
	if !exists {
		cfg = &ssh_config.Config{Hosts: []*ssh_config.Host{}}
		m.Configs[absTarget] = cfg
		m.FileOrder = append(m.FileOrder, absTarget)

		if absTarget != m.PrimaryPath {
			if err := m.registerInclude(absTarget); err != nil {
				return err
			}
		}
	}

	hostBlockStr := buildHostString(h)
	decoded, err := ssh_config.Decode(strings.NewReader(hostBlockStr))
	if err != nil {
		return err
	}

	var newASTHost *ssh_config.Host
	for _, astHost := range decoded.Hosts {
		val := reflect.ValueOf(astHost)
		isImplicit := false
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			elem := val.Elem()
			implicitField := elem.FieldByName("implicit")
			if implicitField.IsValid() && implicitField.Kind() == reflect.Bool && implicitField.Bool() {
				isImplicit = true
			}
		}
		if !isImplicit && len(astHost.Patterns) > 0 {
			newASTHost = astHost
			break
		}
	}

	if newASTHost != nil {
		cfg.Hosts = append(cfg.Hosts, newASTHost)
	}
	return m.SaveFile(absTarget)
}

// UpdateHost edits an existing Host block matched by its original alias.
// It preserves formatting, comments, and other unrecognized nodes in the block.
func (m *Manager) UpdateHost(originalAlias string, h *Host) error {
	targetIdx := -1
	var targetCfg *ssh_config.Config
	var targetFile string

	for filePath, cfg := range m.Configs {
		for i, astHost := range cfg.Hosts {
			for _, pat := range astHost.Patterns {
				if pat.String() == originalAlias {
					targetIdx = i
					targetCfg = cfg
					targetFile = filePath
					break
				}
			}
			if targetIdx != -1 {
				break
			}
		}
	}

	if targetIdx == -1 {
		return fmt.Errorf("host alias %q not found", originalAlias)
	}

	hostBlockStr := buildHostString(h)
	decoded, err := ssh_config.Decode(strings.NewReader(hostBlockStr))
	if err != nil {
		return err
	}

	var newASTHost *ssh_config.Host
	for _, astHost := range decoded.Hosts {
		val := reflect.ValueOf(astHost)
		isImplicit := false
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			elem := val.Elem()
			implicitField := elem.FieldByName("implicit")
			if implicitField.IsValid() && implicitField.Kind() == reflect.Bool && implicitField.Bool() {
				isImplicit = true
			}
		}
		if !isImplicit && len(astHost.Patterns) > 0 {
			newASTHost = astHost
			break
		}
	}

	if newASTHost != nil {
		// why: user comments must be preserved inside the edited host block to prevent data loss
		for _, node := range targetCfg.Hosts[targetIdx].Nodes {
			if _, isComment := node.(*ssh_config.Empty); isComment {
				newASTHost.Nodes = append(newASTHost.Nodes, node)
			}
		}
		targetCfg.Hosts[targetIdx] = newASTHost
	}

	return m.SaveFile(targetFile)
}

// DeleteHost removes a Host entry matched by its pattern from its respective file.
func (m *Manager) DeleteHost(alias string) error {
	for filePath, cfg := range m.Configs {
		for i, astHost := range cfg.Hosts {
			for _, pat := range astHost.Patterns {
				if pat.String() == alias {
					cfg.Hosts = append(cfg.Hosts[:i], cfg.Hosts[i+1:]...)
					return m.SaveFile(filePath)
				}
			}
		}
	}
	return fmt.Errorf("host alias %q not found", alias)
}

// MoveHost transfers a host block from one configuration file to another.
func (m *Manager) MoveHost(alias string, targetFile string) error {
	absTarget := expandTilde(targetFile)
	if abs, err := filepath.Abs(absTarget); err == nil {
		absTarget = abs
	}

	var foundHost *ssh_config.Host
	var sourceFile string

	for file, cfg := range m.Configs {
		for i, astHost := range cfg.Hosts {
			for _, pat := range astHost.Patterns {
				if pat.String() == alias {
					foundHost = astHost
					sourceFile = file
					cfg.Hosts = append(cfg.Hosts[:i], cfg.Hosts[i+1:]...)
					break
				}
			}
			if foundHost != nil {
				break
			}
		}
	}

	if foundHost == nil {
		return fmt.Errorf("host alias %q not found for move", alias)
	}

	targetCfg, exists := m.Configs[absTarget]
	if !exists {
		targetCfg = &ssh_config.Config{Hosts: []*ssh_config.Host{}}
		m.Configs[absTarget] = targetCfg
		m.FileOrder = append(m.FileOrder, absTarget)

		if absTarget != m.PrimaryPath {
			if err := m.registerInclude(absTarget); err != nil {
				return err
			}
		}
	}

	targetCfg.Hosts = append(targetCfg.Hosts, foundHost)

	if err := m.SaveFile(sourceFile); err != nil {
		return err
	}
	return m.SaveFile(absTarget)
}

// SaveFile serializes the specified config AST and writes it back to disk losslessly.
func (m *Manager) SaveFile(filePath string) error {
	cfg, exists := m.Configs[filePath]
	if !exists {
		return fmt.Errorf("config for file %q not tracked", filePath)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	content := cfg.String()
	return os.WriteFile(filePath, []byte(content), 0600)
}

// registerInclude appends an Include statement to the primary file if not already present.
func (m *Manager) registerInclude(includePath string) error {
	primaryCfg, exists := m.Configs[m.PrimaryPath]
	if !exists {
		primaryCfg = &ssh_config.Config{Hosts: []*ssh_config.Host{}}
		m.Configs[m.PrimaryPath] = primaryCfg
	}

	relPath := includePath
	if primaryDir := filepath.Dir(m.PrimaryPath); strings.HasPrefix(includePath, primaryDir) {
		if rel, err := filepath.Rel(primaryDir, includePath); err == nil {
			relPath = rel
		}
	}

	decoded, err := ssh_config.Decode(strings.NewReader("Include " + relPath + "\n"))
	if err != nil || len(decoded.Hosts) == 0 {
		return fmt.Errorf("failed to generate include AST node: %w", err)
	}

	primaryCfg.Hosts = append([]*ssh_config.Host{decoded.Hosts[0]}, primaryCfg.Hosts...)
	return m.SaveFile(m.PrimaryPath)
}

// buildHostString generates the raw formatted host block text for decoding.
func buildHostString(h *Host) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Host %s\n", h.Alias)
	if h.Name != "" {
		fmt.Fprintf(&sb, "    %s %s\n", keyHostName, h.Name)
	}
	if h.User != "" {
		fmt.Fprintf(&sb, "    %s %s\n", keyUser, h.User)
	}
	if h.Port != "" {
		fmt.Fprintf(&sb, "    %s %s\n", keyPort, h.Port)
	}
	if h.IdentityFile != "" {
		fmt.Fprintf(&sb, "    %s %s\n", keyIdentityFile, h.IdentityFile)
	}

	for k, v := range h.Properties {
		if k != keyHostName && k != keyUser && k != keyPort && k != keyIdentityFile && v != "" {
			fmt.Fprintf(&sb, "    %s %s\n", k, v)
		}
	}
	return sb.String()
}
