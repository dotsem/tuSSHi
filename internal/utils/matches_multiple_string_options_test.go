package utils_test

import (
	"testing"
	"tusshi/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestMatchesMultipleStringOptions(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		options  string
		expected bool
	}{
		{"a", "a,b,c", true},
		{"b", "a,b,c", true},
		{"c", "a,b,c", true},
		{"d", "a,b,c", false},
		{"", "a,b,c", false},
		{"a", "", false},
		{"a", " a , b, c ", true},
		{"A", "a,b,c", false},
		{"a,b,c", "a,b,c", false},
		{" a ", " a , b, c ", true},
	}

	for _, test := range tests {
		actual := utils.MatchesMultipleStringOptions(test.input, test.options)
		assert.Equal(test.expected, actual, "Input: %q, Options: %q", test.input, test.options)
	}
}
