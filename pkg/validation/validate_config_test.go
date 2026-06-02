// Package validation contains testing suites for input validation functions.
package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateConfigName checks the SSH config file name validation rules.
func TestValidateConfigName(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty config name",
			config:  "",
			wantErr: true,
			errMsg:  "config name is required",
		},
		{
			name:    "config name with spaces",
			config:  "my config",
			wantErr: true,
			errMsg:  "config name cannot contain spaces",
		},
		{
			name:    "config name with multiple spaces",
			config:  "my  config",
			wantErr: true,
			errMsg:  "config name cannot contain spaces",
		},
		{
			name:    "valid config name",
			config:  "my-config",
			wantErr: false,
		},
		{
			name:    "valid config name with extension",
			config:  "config.txt",
			wantErr: false,
		},
		{
			name:    "forbidden less than",
			config:  "config<01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden greater than",
			config:  "config>01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden colon",
			config:  "config:01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden double quote",
			config:  `config"01`,
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden forward slash",
			config:  "config/01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden backslash",
			config:  `config\01`,
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden pipe",
			config:  "config|01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden question mark",
			config:  "config?01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden asterisk",
			config:  "config*01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
		{
			name:    "forbidden hash",
			config:  "config#01",
			wantErr: true,
			errMsg:  "config name cannot contain forbidden characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigName(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
