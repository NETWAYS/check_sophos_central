package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMainAlerts_matches(t *testing.T) {
	testcases := map[string]struct {
		input    string
		regex    []string
		expected bool
	}{
		"simple-t": {
			input:    "unittest",
			regex:    []string{"unit", "foobar"},
			expected: true,
		},
		"simple-f": {
			input:    "unittest",
			regex:    []string{"barfoo", "foobar"},
			expected: false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := matches(tc.input, tc.regex)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
