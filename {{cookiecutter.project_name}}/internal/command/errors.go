package command

import "github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"

// Sentinel errors of the CLI surface. Each declaration documents why it
// carries its exit code, so the classification survives review and
// refactoring. Sentinels are immutable; attach per-invocation context with
// Wrap or WithDetails, which copy.

// ErrUsage classifies the CLI framework's own parse failures (unknown flags,
// bad flag values): exit 64 (EX_USAGE) because the invocation, not the data
// or the environment, is wrong.
var ErrUsage = terr.New("usage_error", 64,
	"run '{{ cookiecutter.project_name }} --help' for usage",
	"invalid usage")

// ErrNoCommand rejects a bare invocation with no command: exit 64 (EX_USAGE)
// because the CLI surface requires a command.
var ErrNoCommand = terr.New("no_command", 64,
	"run '{{ cookiecutter.project_name }} --help' for available commands",
	"no command specified")

// ErrUnknownCommand rejects an argument that names no command: exit 64
// (EX_USAGE) because the invocation is wrong.
var ErrUnknownCommand = terr.New("unknown_command", 64,
	"run '{{ cookiecutter.project_name }} --help' for available commands",
	"unknown command")

// ErrEmptyName rejects an explicitly empty name argument to greet: exit 65
// (EX_DATAERR) because the CLI surface was used correctly but the input
// payload is invalid. It is a placeholder demonstrating a data-class
// sentinel; replace it with the tool's real data errors.
var ErrEmptyName = terr.New("empty_name", 65,
	"pass a non-empty name, or no argument for the default",
	"name must not be empty")
