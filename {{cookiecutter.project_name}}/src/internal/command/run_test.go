package command_test

import (
	"bytes"
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/command"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"
)

// failWriter simulates a broken stdout sink.
type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("sink failed") }

// TestUnclassifiedErrorIsInternal drives a write failure through the exit
// boundary: the resulting error carries no code, so it must surface as
// internal_error with exit 70 rather than being silently mislabeled. This is
// also the scenario that exercises exit 70 for the registry-coverage check.
func TestUnclassifiedErrorIsInternal(t *testing.T) {
	var stderr bytes.Buffer
	exit := command.Run([]string{"greet"}, command.Deps{
		In:  strings.NewReader(""),
		Out: failWriter{},
		Err: &stderr,
	})

	observedMu.Lock()
	observedExits[exit] = true
	observedMu.Unlock()

	if exit != 70 {
		t.Errorf("exit code = %d, want 70", exit)
	}
	if !strings.Contains(stderr.String(), `"code":"internal_error"`) {
		t.Errorf("stderr does not carry the internal_error envelope: %s", stderr.String())
	}
}

// TestRegistryMatchesDeclarations asserts the exit-code registry equals the
// set derivable from the code: 0 for success, 70 for the unclassified
// fallback, plus the exit code of every registered sentinel. A sentinel with
// a new exit code fails this test until the registry declares it.
func TestRegistryMatchesDeclarations(t *testing.T) {
	want := map[int]bool{0: true, 70: true}
	for _, e := range terr.All() {
		want[e.ExitCode()] = true
	}

	var wantCodes []int
	for c := range want {
		wantCodes = append(wantCodes, c)
	}
	sort.Ints(wantCodes)

	got := append([]int(nil), command.ExitCodes...)
	sort.Ints(got)

	if len(got) != len(wantCodes) {
		t.Fatalf("registry = %v, want %v", got, wantCodes)
	}
	for i := range got {
		if got[i] != wantCodes[i] {
			t.Fatalf("registry = %v, want %v", got, wantCodes)
		}
	}
}
