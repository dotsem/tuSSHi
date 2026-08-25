// Package utils provides general helper functions for string matching and formatting.
package utils

import "strings"

// MatchesMultipleStringOptions checks if input matches any of the comma-separated options.
func MatchesMultipleStringOptions(input string, options string) bool {
	cmds := strings.SplitSeq(options, ",")
	for cmd := range cmds {
		if strings.TrimSpace(input) == strings.TrimSpace(cmd) {
			return true
		}
	}
	return false
}
