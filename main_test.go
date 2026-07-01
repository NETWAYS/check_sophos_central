package main

import (
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestMainCheck_matches(t *testing.T) {
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
			if tc.expected != actual {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestMainAlerts_GetPerfdata(t *testing.T) {
	testcases := map[string]struct {
		ao       AlertOverview
		expected string
	}{
		"simple-overview": {
			ao:       AlertOverview{},
			expected: "alerts=0 alerts_high=0 alerts_medium=0 alerts_low=0",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := tc.ao.GetPerfdata()
			if tc.expected != actual {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestMainAlerts_GetOutput(t *testing.T) {
	testcases := map[string]struct {
		ao       AlertOverview
		expected string
	}{
		"simple-overview": {
			ao: AlertOverview{
				Total:  1,
				Output: []string{"test"},
			},
			expected: "\n## Alerts\ntest\n",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := tc.ao.GetOutput()
			if tc.expected != actual {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestMainAlerts_GetStatus(t *testing.T) {
	testcases := map[string]struct {
		ao       AlertOverview
		expected check.Status
	}{
		"simple-overview": {
			ao:       AlertOverview{},
			expected: check.OK,
		},
		"simple-warning": {
			ao: AlertOverview{
				Medium: 1,
			},
			expected: check.Warning,
		},
		"simple-critical": {
			ao: AlertOverview{
				High: 1,
			},
			expected: check.Critical,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := tc.ao.GetStatus()
			if tc.expected != actual {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
