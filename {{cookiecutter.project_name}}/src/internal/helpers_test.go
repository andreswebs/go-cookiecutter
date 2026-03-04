package config

import (
	"testing"
)

func TestRequiredString_Missing(t *testing.T) {
	t.Setenv("REQUIRED_TEST_VAR", "")

	_, err := requiredString("REQUIRED_TEST_VAR")
	if err == nil {
		t.Fatal("requiredString(\"REQUIRED_TEST_VAR\") error = nil, want error")
	}
}

func TestRequiredString_Present(t *testing.T) {
	t.Setenv("REQUIRED_TEST_VAR", "hello")

	got, err := requiredString("REQUIRED_TEST_VAR")
	if err != nil {
		t.Fatalf("requiredString(\"REQUIRED_TEST_VAR\") error = %v, want nil", err)
	}

	if got != "hello" {
		t.Errorf("requiredString(\"REQUIRED_TEST_VAR\") = %q, want %q", got, "hello")
	}
}
