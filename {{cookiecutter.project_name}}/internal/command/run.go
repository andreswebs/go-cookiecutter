// Package command owns the CLI surface. This file is the framework-free
// contract: main calls Run, tests drive Run with buffers, and no CLI
// framework type appears in any exported identifier of the package. The
// framework interior lives in root.go, commands.go, and usage.go; replacing
// the framework replaces those files and keeps this one and the tests.
package command

import (
	"context"
	"errors"
	"io"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/output"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"
)

// Deps carries the injected process environment. It starts lean; a field is
// added only when a test must fake it (Getenv, a clock, terminal detection).
// Arguments are a parameter of Run, not a dependency.
type Deps struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// Run is the delegate entry point: it executes the CLI with the given
// arguments (without the program name) and injected streams, and returns the
// process exit code. It owns the exit boundary: typed errors choose the exit
// code, the error envelope is emitted to deps.Err, and nothing else inspects
// errors. A tool that catches SIGINT/SIGTERM for graceful shutdown overrides
// the return value here with 128 plus the signal number.
func Run(args []string, deps Deps) int {
	err := runRoot(context.Background(), args, deps)
	if err == nil {
		return 0
	}

	var coded terr.Coded
	if !errors.As(err, &coded) {
		if usageErr := classifyUsage(err); usageErr != nil {
			err = usageErr
		}
	}

	output.EmitError(deps.Err, err)
	return output.ExitCodeFor(err)
}
