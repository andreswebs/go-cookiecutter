package terr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"
)

var errSentinel = terr.New("test_sentinel", 65, "fix the input", "input rejected")

func TestNewRegistersForEnumeration(t *testing.T) {
	for _, e := range terr.All() {
		if e == errSentinel {
			return
		}
	}
	t.Fatal("sentinel created with New is missing from All()")
}

func TestNewPanicsOnDuplicateCode(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("New with a duplicate code did not panic")
		}
	}()
	terr.New("test_sentinel", 70, "", "duplicate")
}

func TestNewfDoesNotRegister(t *testing.T) {
	e := terr.Newf("test_oneoff", 65, "", "bad value %q", "x")
	if got, want := e.Error(), `bad value "x"`; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	for _, r := range terr.All() {
		if r.Code() == "test_oneoff" {
			t.Fatal("Newf error leaked into the registry")
		}
	}
}

func TestCodedAccessors(t *testing.T) {
	if errSentinel.Code() != "test_sentinel" {
		t.Errorf("Code() = %q", errSentinel.Code())
	}
	if errSentinel.ExitCode() != 65 {
		t.Errorf("ExitCode() = %d", errSentinel.ExitCode())
	}
	if errSentinel.Hint() != "fix the input" {
		t.Errorf("Hint() = %q", errSentinel.Hint())
	}
}

func TestWrapCopiesAndPreservesIdentity(t *testing.T) {
	cause := errors.New("underlying")
	wrapped := errSentinel.Wrap(cause)

	if errSentinel.Unwrap() != nil {
		t.Error("Wrap mutated the sentinel")
	}
	if !errors.Is(wrapped, errSentinel) {
		t.Error("wrapped copy does not match its sentinel under errors.Is")
	}
	if !errors.Is(wrapped, cause) {
		t.Error("wrapped copy does not expose its cause via Unwrap")
	}
	if got, want := wrapped.Error(), "input rejected: underlying"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrorsAsThroughFmtWrap(t *testing.T) {
	err := fmt.Errorf("loading profile: %w", errSentinel)

	var coded terr.Coded
	if !errors.As(err, &coded) {
		t.Fatal("errors.As did not find a Coded in the chain")
	}
	if coded.Code() != "test_sentinel" {
		t.Errorf("Code() = %q", coded.Code())
	}
}

func TestWithDetailsCopies(t *testing.T) {
	detailed := errSentinel.WithDetails(map[string]string{"field": "name"})

	if errSentinel.ErrorDetails() != nil {
		t.Error("WithDetails mutated the sentinel")
	}
	if detailed.ErrorDetails() == nil {
		t.Error("copy is missing its details")
	}
	if !errors.Is(detailed, errSentinel) {
		t.Error("detailed copy does not match its sentinel under errors.Is")
	}
}
