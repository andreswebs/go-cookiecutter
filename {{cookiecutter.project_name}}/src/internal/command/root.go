// The framework interior: everything in this file (and commands.go and
// usage.go) is urfave/cli specific and replaceable. The no-leak rule applies:
// the framework is imported only inside this package, no framework type
// appears in an exported identifier, and the framework's own error printing
// and exit handling are neutralized here so run.go maps errors after the
// framework returns.

package command

import (
	"context"
	"fmt"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/version"

	"github.com/urfave/cli/v3"
)

// runRoot builds the framework command tree bound to the injected streams and
// executes it. It is the only seam between the contract (run.go) and the
// framework.
func runRoot(ctx context.Context, args []string, deps Deps) error {
	root := newRoot(deps)
	return root.Run(ctx, append([]string{root.Name}, args...))
}

func newRoot(deps Deps) *cli.Command {
	root := &cli.Command{
		Name:      "{{ cookiecutter.project_name }}",
		Usage:     "{{ cookiecutter.project_short_description }}",
		Version:   version.Current(),
		Reader:    deps.In,
		Writer:    deps.Out,
		ErrWriter: deps.Err,
		Commands:  subcommands(),
		Action:    rootAction,
	}
	neutralize(root)
	return root
}

// subcommands declares the command tree once; the framework wiring and the
// schema command's self-description are both projections of this declaration.
func subcommands() []*cli.Command {
	return []*cli.Command{
		greetCommand(),
		schemaCommand(),
	}
}

// neutralize disables the framework's own error printing, help-on-error, and
// exit handling on cmd and every command below it: parse errors must return
// to the boundary in run.go untouched, so stderr carries exactly one error
// envelope and the framework never calls os.Exit.
func neutralize(cmd *cli.Command) {
	cmd.ExitErrHandler = func(context.Context, *cli.Command, error) {}
	cmd.OnUsageError = func(_ context.Context, _ *cli.Command, err error, _ bool) error {
		return err
	}
	for _, sub := range cmd.Commands {
		neutralize(sub)
	}
}

// rootAction handles invocations that reach the root itself: a bare call or a
// first argument that names no command. Both are usage errors.
func rootAction(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Present() {
		return ErrUnknownCommand.Wrap(fmt.Errorf("%q", cmd.Args().First()))
	}
	return ErrNoCommand
}
