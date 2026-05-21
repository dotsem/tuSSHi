package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// Manager coordinates the loading, parsing, editing, and writing of one or more
// OpenSSH configuration files. It supports Include glob directives natively.
type Manager struct {
	// PrimaryPath is the absolute path to the main SSH config file (usually ~/.ssh/config).
	PrimaryPath string

	// Configs holds the parsed AST of each individual config file, keyed by its absolute path.
	Configs map[string]*ssh_config.Config

	// FileOrder tracks the order in which files are discovered to preserve hierarchy during display.
	FileOrder []string
}

// NewManager creates and initializes a Manager with the specified primary config file.
// It performs basic path expansion on the provided primary path.
func NewManager(primaryPath string) *Manager {
	absPath := expandTilde(primaryPath)
	if abs, err := filepath.Abs(absPath); err == nil {
		absPath = abs
	}
	return &Manager{
		PrimaryPath: absPath,
		Configs:     make(map[string]*ssh_config.Config),
		FileOrder:   []string{absPath},
	}
}

// Load reads and parses the primary configuration file and recursively parses all
// included configuration files discovered via the Include directives.
func (m *Manager) Load() error {
	m.Configs = make(map[string]*ssh_config.Config)
	m.FileOrder = []string{m.PrimaryPath}
	return m.loadPath(m.PrimaryPath, 0)
}

// GetHosts converts the parsed ASTs into clean, high-level Host models.
// It automatically resolves wildcard settings for each specific host.
func (m *Manager) GetHosts() []*Host {
	var hosts []*Host
	// We gather global wildcard configs to perform inheritance resolution later.
	globalConfig := m.buildGlobalConfig()

	for _, filePath := range m.FileOrder {
		cfg, exists := m.Configs[filePath]
		if !exists {
			continue
		}

		for _, astHost := range cfg.Hosts {
			// Skip the implicit default "Host *" block added by the parser
			val := reflect.ValueOf(astHost)
			if val.Kind() == reflect.Ptr && !val.IsNil() {
				elem := val.Elem()
				implicitField := elem.FieldByName("implicit")
				if implicitField.IsValid() && implicitField.Kind() == reflect.Bool && implicitField.Bool() {
					continue
				}
			}

			// Extract all aliases defined in this Host block.
			for _, pat := range astHost.Patterns {
				alias := pat.String()
				if alias == "" {
					continue
				}

				isWildcard := strings.ContainsAny(alias, "*?")
				h := &Host{
					Alias:      alias,
					SourceFile: filePath,
					IsWildcard: isWildcard,
					Properties: make(map[string]string),
				}

				// Extract explicit key-value properties from the host block's nodes.
				for _, node := range astHost.Nodes {
					if kv, ok := node.(*ssh_config.KV); ok {
						h.Properties[kv.Key] = kv.Value
					}
				}

				// Map critical properties to top-level fields for convenience.
				h.Name = h.Properties["HostName"]
				h.User = h.Properties["User"]
				h.Port = h.Properties["Port"]
				h.IdentityFile = h.Properties["IdentityFile"]

				// Resolve final values using global OpenSSH inheritance.
				h.ResolvedProperties = make(map[string]string)
				for k, v := range h.Properties {
					h.ResolvedProperties[k] = v
				}

				// Inject inherited properties from matching wildcard blocks.
				if !isWildcard && globalConfig != nil {
					for _, key := range []string{"HostName", "User", "Port", "IdentityFile", "ForwardAgent", "ProxyJump"} {
						if _, explicit := h.Properties[key]; !explicit {
							if resolvedVal, err := globalConfig.Get(alias, key); err == nil && resolvedVal != "" {
								h.ResolvedProperties[key] = resolvedVal
							}
						}
					}
				}

				// Update resolved shortcuts.
				if h.Name == "" && h.ResolvedProperties["HostName"] != "" {
					h.Name = h.ResolvedProperties["HostName"]
				}
				if h.User == "" && h.ResolvedProperties["User"] != "" {
					h.User = h.ResolvedProperties["User"]
				}
				if h.Port == "" && h.ResolvedProperties["Port"] != "" {
					h.Port = h.ResolvedProperties["Port"]
				}
				if h.IdentityFile == "" && h.ResolvedProperties["IdentityFile"] != "" {
					h.IdentityFile = h.ResolvedProperties["IdentityFile"]
				}

				hosts = append(hosts, h)
			}
		}
	}
	return hosts
}

// loadPath handles the recursive parsing and include tracking.
func (m *Manager) loadPath(path string, depth int) error {
	if depth > 5 {
		return fmt.Errorf("exceeded max include recursion depth")
	}

	if _, exists := m.Configs[path]; exists {
		return nil
	}

	// #nosec G304 - path is a config file path loaded from the SSH configuration
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) && path == m.PrimaryPath {
			// If primary file doesn't exist, we start with a clean empty config
			m.Configs[path] = &ssh_config.Config{Hosts: []*ssh_config.Host{}}
			return nil
		}
		return err
	}
	defer func() {
		_ = f.Close() // TODO: should this be handled?
	}()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return err
	}

	m.Configs[path] = cfg

	// Scan AST nodes to discover and parse any Include directives.
	for _, astHost := range cfg.Hosts {
		for _, node := range astHost.Nodes {
			if incl, ok := node.(*ssh_config.Include); ok {
				inclStr := strings.TrimSpace(incl.String())
				if idx := strings.Index(inclStr, "#"); idx != -1 {
					inclStr = strings.TrimSpace(inclStr[:idx])
				}
				inclStr = strings.TrimPrefix(inclStr, "Include")
				inclStr = strings.TrimSpace(inclStr)
				inclStr = strings.TrimPrefix(inclStr, "=")
				inclStr = strings.TrimSpace(inclStr)

				for _, pattern := range strings.Fields(inclStr) {
					m.resolveAndLoadIncludes(pattern, depth+1)
				}
			}
		}
	}
	return nil
}

// resolveAndLoadIncludes matches globs and recursively loads matched files.
func (m *Manager) resolveAndLoadIncludes(pattern string, depth int) {
	expanded := expandTilde(pattern)
	if !filepath.IsAbs(expanded) {
		expanded = filepath.Join(filepath.Dir(m.PrimaryPath), expanded)
	}

	matches, err := filepath.Glob(expanded)
	if err != nil {
		return
	}

	for _, match := range matches {
		absMatch, err := filepath.Abs(match)
		if err != nil {
			continue
		}

		if err := m.loadPath(absMatch, depth); err == nil {
			// Track order of newly discovered files.
			found := slices.Contains(m.FileOrder, absMatch)
			if !found {
				m.FileOrder = append(m.FileOrder, absMatch)
			}
		}
	}
}

// buildGlobalConfig merges all parsed configs into a single config object for queries.
func (m *Manager) buildGlobalConfig() *ssh_config.Config {
	var mergedHosts []*ssh_config.Host
	for _, path := range m.FileOrder {
		if cfg, exists := m.Configs[path]; exists {
			mergedHosts = append(mergedHosts, cfg.Hosts...)
		}
	}
	return &ssh_config.Config{Hosts: mergedHosts}
}

// expandTilde replaces ~/ prefix with the user home directory path.
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
