// Package validation provides validation logic for connection aliases and config names.
package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var forbiddenAliasChars = regexp.MustCompile(`[<>:"/\\|?*$]`)
var forbiddenConfigNameChars = regexp.MustCompile(`[<>:"/\\|?*#]`)

// ValidateAlias checks if the alias matches the validation rules.
func ValidateAlias(str string) error {
	if str == "" {
		return errors.New("alias is required")
	}
	if strings.Contains(str, " ") {
		return errors.New("alias cannot contain spaces")
	}
	if forbiddenAliasChars.MatchString(str) {
		return fmt.Errorf("alias cannot contain forbidden characters (%v)", strings.Split(forbiddenAliasChars.String(), ""))
	}
	return nil
}

// ValidateConfigName checks if the configuration name matches the validation rules.
func ValidateConfigName(str string) error {
	if str == "" {
		return errors.New("config name is required")
	}
	if strings.Contains(str, " ") {
		return errors.New("config name cannot contain spaces")
	}
	if forbiddenConfigNameChars.MatchString(str) {
		return fmt.Errorf("config name cannot contain forbidden characters (%v)", strings.Split(forbiddenConfigNameChars.String(), ""))
	}
	return nil
}
