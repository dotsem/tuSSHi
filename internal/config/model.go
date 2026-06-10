package config

// Host represents a high-level representation of an SSH connection definition.
// It maps the raw AST structure of a Host block in an SSH config file to a clean
// Go struct suitable for search, display, and connection execution.
type Host struct {
	// Alias is the identifier pattern used to match the SSH command (e.g., "prod-web").
	// This represents the primary match pattern of the Host block.
	Alias string

	// Name is the actual destination HostName (e.g., "10.200.1.45" or "prod.example.com").
	// If the HostName is not explicitly defined, this remains empty or inherits.
	Name string

	// User is the remote login user (e.g., "deploy" or "root").
	User string

	// Port is the SSH port number.
	Port string

	// IdentityFile is the path to the private key used for authentication.
	IdentityFile string

	// SourceFile is the absolute path to the configuration file on disk
	// where this connection is explicitly defined.
	SourceFile string

	// IsWildcard indicates whether this host block is a global/wildcard block
	// (e.g., "Host *") rather than a specific destination connection.
	IsWildcard bool

	// Properties stores all configuration parameters explicitly defined under this host
	// block as key-value pairs (e.g., "ForwardAgent": "yes").
	Properties map[string]string

	// ResolvedProperties stores all parameters including those inherited from wildcards
	// or system defaults, providing a complete view of the resolved connection state.
	ResolvedProperties map[string]string
}
