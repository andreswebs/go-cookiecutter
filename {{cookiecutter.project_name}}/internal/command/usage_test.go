package command_test

import (
	"strings"
	"testing"
)

// TestFrameworkParseErrorsClassifyAsUsage drives real framework parse
// failures through Run and asserts they map to exit 64 with the usage_error
// code. If the framework's error wording ever changes, the classifier stops
// matching, the error falls through as internal_error (exit 70), and these
// cases fail loudly.
func TestFrameworkParseErrorsClassifyAsUsage(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"unknown flag", []string{"greet", "--bogus"}},
		{"unknown flag on root", []string{"--bogus"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit, _, stderr := runScenario(t, tt.args)

			if exit != 64 {
				t.Errorf("exit code = %d, want 64", exit)
			}
			if !strings.Contains(stderr.String(), `"code":"usage_error"`) {
				t.Errorf("stderr does not carry the usage_error envelope: %s",
					stderr.String())
			}
		})
	}
}
