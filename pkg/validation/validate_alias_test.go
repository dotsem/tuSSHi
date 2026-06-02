// Package validation contains testing suites for input validation functions.
package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateAlias checks the SSH host connection alias validation rules.
func TestValidateAlias(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty alias",
			alias:   "",
			wantErr: true,
			errMsg:  "alias is required",
		},
		{
			name:    "alias with spaces",
			alias:   "my server",
			wantErr: true,
			errMsg:  "alias cannot contain spaces",
		},
		{
			name:    "alias with multiple spaces",
			alias:   "my  server",
			wantErr: true,
			errMsg:  "alias cannot contain spaces",
		},
		{
			name:    "valid alias",
			alias:   "my-server-01",
			wantErr: false,
		},
		{
			name:    "forbidden less than",
			alias:   "server<01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden greater than",
			alias:   "server>01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden colon",
			alias:   "server:01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden double quote",
			alias:   `server"01`,
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden forward slash",
			alias:   "server/01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden backslash",
			alias:   `server\01`,
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden pipe",
			alias:   "server|01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden question mark",
			alias:   "server?01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden asterisk",
			alias:   "server*01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
		{
			name:    "forbidden dollar",
			alias:   "server$01",
			wantErr: true,
			errMsg:  "alias cannot contain forbidden characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAlias(tt.alias)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
