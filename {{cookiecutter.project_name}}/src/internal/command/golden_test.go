// The golden-triple harness: every scenario drives Run with buffer streams
// and compares stdout, stderr, and the exit code against golden files under
// testdata/<scenario>/. Regenerate with `go test ./internal/command -update`
// after a contract change; the golden diff is the reviewable evidence.
// Because scenarios touch only the Run contract, the suite is framework-blind
// and survives a CLI-framework replacement unchanged.

package command_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/command"
)

var update = flag.Bool("update", false, "rewrite golden files")

// observedExits records every exit code produced through runScenario, so
// TestMain can enforce that the declared registry is fully exercised.
var (
	observedMu    sync.Mutex
	observedExits = make(map[int]bool)
)

func TestMain(m *testing.M) {
	flag.Parse()
	code := m.Run()
	// Registry coverage is checked only on a full, non-updating run: a
	// -run filter legitimately exercises a subset.
	if code == 0 && !*update && flag.Lookup("test.run").Value.String() == "" {
		var missing []int
		for _, c := range command.ExitCodes {
			if !observedExits[c] {
				missing = append(missing, c)
			}
		}
		if len(missing) > 0 {
			fmt.Fprintf(os.Stderr,
				"declared exit codes not exercised by any scenario: %v\n", missing)
			code = 1
		}
	}
	os.Exit(code)
}

func TestGolden(t *testing.T) {
	scenarios := []struct {
		name string
		args []string
	}{
		{"greet_default", []string{"greet"}},
		{"greet_name", []string{"greet", "gopher"}},
		{"greet_empty_name", []string{"greet", ""}},
		{"usage_unknown_flag", []string{"greet", "--bogus"}},
		{"unknown_command", []string{"bogus"}},
		{"no_command", nil},
		{"schema", []string{"schema"}},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			runGolden(t, sc.name, sc.args)
		})
	}
}

// runScenario drives Run with buffer streams, records the exit code, and
// asserts registry membership: every exit the tool can produce must be
// declared in command.ExitCodes.
func runScenario(t *testing.T, args []string) (exit int, stdout, stderr *bytes.Buffer) {
	t.Helper()

	stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	exit = command.Run(args, command.Deps{
		In:  strings.NewReader(""),
		Out: stdout,
		Err: stderr,
	})

	observedMu.Lock()
	observedExits[exit] = true
	observedMu.Unlock()

	if !slices.Contains(command.ExitCodes, exit) {
		t.Errorf("exit code %d is not declared in the registry %v",
			exit, command.ExitCodes)
	}
	return exit, stdout, stderr
}

func runGolden(t *testing.T, name string, args []string) {
	t.Helper()

	exit, stdout, stderr := runScenario(t, args)
	dir := filepath.Join("testdata", name)

	if *update {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
		writeGolden(t, filepath.Join(dir, "stdout"), stdout.Bytes())
		writeGolden(t, filepath.Join(dir, "stderr"), stderr.Bytes())
		writeGolden(t, filepath.Join(dir, "exit"), []byte(strconv.Itoa(exit)+"\n"))
		t.Logf("updated golden files in %s", dir)
		return
	}

	compareGolden(t, filepath.Join(dir, "stdout"), stdout.Bytes())
	compareGolden(t, filepath.Join(dir, "stderr"), stderr.Bytes())

	wantExitRaw, err := os.ReadFile(filepath.Join(dir, "exit"))
	if err != nil {
		t.Fatalf("read golden exit: %v", err)
	}
	wantExit, err := strconv.Atoi(strings.TrimSpace(string(wantExitRaw)))
	if err != nil {
		t.Fatalf("parse golden exit: %v", err)
	}
	if exit != wantExit {
		t.Errorf("exit code = %d, want %d", exit, wantExit)
	}
}

func writeGolden(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, normGolden(data), 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
}

func compareGolden(t *testing.T, path string, got []byte) {
	t.Helper()

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}

	want = normGolden(want)
	got = normGolden(got)

	if !bytes.Equal(got, want) {
		t.Errorf("mismatch %s:\ngot:\n%s\nwant:\n%s", path, got, want)
	}
}

// normGolden normalizes volatile values so golden files stay deterministic:
// CRLF line endings and the working directory. Extend it (ports, timestamps,
// hostnames) rather than letting a volatile value into a golden file, which
// is the harness's main failure mode.
func normGolden(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	if cwd, err := os.Getwd(); err == nil {
		b = bytes.ReplaceAll(b, []byte(cwd), []byte("$CWD"))
	}
	return b
}
