package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tussh "tusshi/internal/ssh"

	"github.com/charmbracelet/huh"
)

const (
	keySourceGenerate = "generate"
	keySourceExisting = "existing"
)

// ServiceFormState holds all mutable form values for adding or editing an SSH service connection.
type ServiceFormState struct {
	Action        string // "add" or "edit"
	OriginalAlias string
	KeySource     string
	KeyType       string
	KeyPath       string
	KeyComment    string
	PresetAlias   string
	HostAlias     string
	HostName      string
	HostUser      string

	lastPreset  string
	lastKeyType string
}

// BuildServiceForm constructs the multi-step form for adding or editing a service host.
func BuildServiceForm(s *ServiceFormState) *huh.Form {
	if s.Action == actionAdd {
		s.ApplyPreset()
		s.lastPreset = s.PresetAlias
		s.lastKeyType = s.KeyType
	}

	presetOptions := make([]huh.Option[string], len(tussh.Presets))
	for i, p := range tussh.Presets {
		presetOptions[i] = huh.NewOption(p.Name, p.Alias)
	}

	keyTypeOptions := []huh.Option[string]{
		huh.NewOption("ed25519 (recommended)", tussh.KeyTypeED25519),
		huh.NewOption("RSA 4096", tussh.KeyTypeRSA),
		huh.NewOption("ECDSA", tussh.KeyTypeECDSA),
	}

	keySourceOptions := []huh.Option[string]{
		huh.NewOption("Generate a new key", keySourceGenerate),
		huh.NewOption("Use an existing key file", keySourceExisting),
	}

	inputKeyPath := huh.NewInput().
		Title("Output Key Path").
		Description("Where to save the private key").
		Value(&s.KeyPath)

	inputExistingKeyPath := huh.NewInput().
		Title("Private Key Path").
		Description("Absolute or ~ path to your existing private key").
		Placeholder("~/.ssh/id_ed25519").
		Value(&s.KeyPath).
		Validate(func(v string) error {
			expanded := expandTildePath(v)
			if _, err := os.Stat(expanded); err != nil {
				return fmt.Errorf("file not found: %s", expanded)
			}
			return nil
		})

	inputHostAlias := huh.NewInput().
		Title("Host Alias").
		Description("Name used in SSH config (e.g. github or github-work)").
		Value(&s.HostAlias)

	inputHostName := huh.NewInput().
		Title("HostName").
		Description("Actual destination hostname").
		Value(&s.HostName)

	inputHostUser := huh.NewInput().
		Title("User").
		Description("Remote user for this service").
		Value(&s.HostUser)

	syncFields := func() {
		if s.Action != actionAdd {
			return
		}
		if s.PresetAlias != s.lastPreset || s.KeyType != s.lastKeyType {
			if preset, ok := tussh.FindPreset(s.PresetAlias); ok {
				s.HostAlias = preset.Alias
				s.HostName = preset.HostName
				s.HostUser = preset.User
			} else if s.PresetAlias == tussh.PresetCustom {
				if s.lastPreset != "" {
					s.HostAlias = ""
					s.HostName = ""
					s.HostUser = "git"
				}
			}
			s.KeyPath = s.ProvideDefaultKeyPath()

			inputKeyPath.Value(&s.KeyPath)
			inputHostAlias.Value(&s.HostAlias)
			inputHostName.Value(&s.HostName)
			inputHostUser.Value(&s.HostUser)

			s.lastPreset = s.PresetAlias
			s.lastKeyType = s.KeyType
		}
	}

	step1 := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Service Preset").
			Description("Which service are you configuring?").
			Options(presetOptions...).
			Value(&s.PresetAlias),
		huh.NewSelect[string]().
			Title("Key Source").
			Description("How do you want to provide the SSH key?").
			Options(keySourceOptions...).
			Value(&s.KeySource),
	)

	step2generate := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Key Type").
			Options(keyTypeOptions...).
			Value(&s.KeyType),
		inputKeyPath,
		huh.NewInput().
			Title("Comment").
			Description("Identifies the key (e.g. email)").
			Placeholder("you@example.com").
			Value(&s.KeyComment),
	).WithHideFunc(func() bool {
		syncFields()
		return s.KeySource != keySourceGenerate
	})

	step2existing := huh.NewGroup(
		inputExistingKeyPath,
	).WithHideFunc(func() bool {
		syncFields()
		return s.KeySource != keySourceExisting
	})

	step3 := huh.NewGroup(
		inputHostAlias,
		inputHostName,
		inputHostUser,
	).WithHideFunc(func() bool {
		syncFields()
		return false
	})

	form := huh.NewForm(step1, step2generate, step2existing, step3).
		WithTheme(huh.ThemeCharm()).
		WithWidth(60).
		WithShowHelp(false)

	return form
}

// ApplyPreset fills HostAlias, HostName, and HostUser from the selected preset if still empty on submit.
func (s *ServiceFormState) ApplyPreset() {
	if preset, ok := tussh.FindPreset(s.PresetAlias); ok {
		s.HostAlias = preset.Alias
		s.HostName = preset.HostName
		s.HostUser = preset.User
	} else if s.PresetAlias == tussh.PresetCustom {
		if s.HostUser == "" {
			s.HostUser = "git"
		}
	}

	s.KeyPath = s.ProvideDefaultKeyPath()
}

// ProvideDefaultKeyPath generates a non-colliding default SSH key path.
func (s *ServiceFormState) ProvideDefaultKeyPath() string {
	alias := s.HostAlias
	if alias == "" {
		alias = s.PresetAlias
	}
	if alias == "" || alias == tussh.PresetCustom {
		alias = "service"
	}
	kType := s.KeyType
	if kType == "" {
		kType = tussh.KeyTypeED25519
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	base := filepath.Join(home, ".ssh", fmt.Sprintf("id_%s_%s", kType, alias))
	target := base
	counter := 1
	for {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			break
		}
		target = fmt.Sprintf("%s_%d", base, counter)
		counter++
	}

	if strings.HasPrefix(target, home) {
		return "~" + target[len(home):]
	}
	return target
}

// ResolvedKeyPath returns the expanded absolute path for the configured key.
func (s *ServiceFormState) ResolvedKeyPath() string {
	return expandTildePath(s.KeyPath)
}

func expandTildePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
