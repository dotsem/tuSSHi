package ssh

const (
	// PresetCustom represents the custom service preset identifier.
	PresetCustom   = "custom"
	defaultUserGit = "git"
)

// ServicePreset describes a known Git/SSH service that uses key-based auth.
type ServicePreset struct {
	Name     string // display name, e.g. "GitHub"
	KeyName  string // base name for SSH key file, e.g. "github"
	HostName string // SSH Host pattern and destination hostname, e.g. "github.com"
	User     string // remote user, always "git" for hosting services
}

// Presets contains built-in service definitions for the key wizard.
var Presets = []ServicePreset{
	{Name: "GitHub", KeyName: "github", HostName: "github.com", User: defaultUserGit},
	{Name: "GitLab", KeyName: "gitlab", HostName: "gitlab.com", User: defaultUserGit},
	// TODO: add more
	{Name: "Custom", KeyName: "service", HostName: "service", User: defaultUserGit},
}

// FindPreset returns the preset matching the given name, hostname, or keyname.
func FindPreset(query string) (ServicePreset, bool) {
	for _, p := range Presets {
		if p.HostName == query || p.Name == query || p.KeyName == query {
			return p, true
		}
	}
	return ServicePreset{}, false
}
