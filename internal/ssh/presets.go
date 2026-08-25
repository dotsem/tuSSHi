package ssh

// PresetCustom represents the custom service preset identifier.
const PresetCustom = "custom"

// ServicePreset describes a known Git/SSH service that uses key-based auth.
type ServicePreset struct {
	Name     string // display name, e.g. "GitHub"
	Alias    string // SSH config Host alias, e.g. "github"
	HostName string // actual destination, e.g. "github.com"
	User     string // remote user, always "git" for hosting services
}

// Presets contains built-in service definitions for the key wizard.
// Custom and Bitbucket entries can be added in future iterations.
var Presets = []ServicePreset{
	{Name: "GitHub", Alias: "github", HostName: "github.com", User: "git"},
	{Name: "GitLab", Alias: "gitlab", HostName: "gitlab.com", User: "git"},
	// TODO: add more
	{Name: "Custom", Alias: "service", HostName: "", User: ""},
}

// FindPreset returns the preset matching the given alias, or false if not found.
func FindPreset(alias string) (ServicePreset, bool) {
	for _, p := range Presets {
		if p.Alias == alias {
			return p, true
		}
	}
	return ServicePreset{}, false
}
