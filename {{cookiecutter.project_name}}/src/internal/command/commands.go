package command

import (
	"context"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/output"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"

	"github.com/urfave/cli/v3"
)

// greetCommand is a placeholder: replace it with the tool's real commands.
// It demonstrates the output contract: one JSON result envelope on stdout,
// coded errors for failures.
func greetCommand() *cli.Command {
	return &cli.Command{
		Name:      "greet",
		Usage:     "print a greeting envelope",
		ArgsUsage: "[name]",
		Action:    greetAction,
	}
}

func greetAction(_ context.Context, cmd *cli.Command) error {
	name := "world"
	if cmd.Args().Present() {
		name = cmd.Args().First()
		if name == "" {
			return ErrEmptyName
		}
	}
	env := output.GreetEnvelope{Head: output.OKHead(), Greeting: "hello, " + name}
	return output.EmitJSON(cmd.Root().Writer, env)
}

// schemaCommand reports the tool's machine-readable contract: the command
// surface, the declared exit codes, and the error inventory. Everything it
// emits is a projection of declarations that drive behavior (the command
// tree, the exit-code registry, the terr registry), so it cannot drift.
func schemaCommand() *cli.Command {
	return &cli.Command{
		Name:   "schema",
		Usage:  "describe the tool's commands, exit codes, and errors",
		Action: schemaAction,
	}
}

func schemaAction(_ context.Context, cmd *cli.Command) error {
	env := output.SchemaEnvelope{
		Head:      output.OKHead(),
		Tool:      cmd.Root().Name,
		Version:   cmd.Root().Version,
		ExitCodes: ExitCodes,
	}
	for _, sub := range subcommands() {
		sc := output.SchemaCommand{
			Name:  sub.Name,
			Usage: sub.Usage,
			Args:  sub.ArgsUsage,
		}
		for _, f := range sub.Flags {
			names := f.Names()
			if len(names) == 0 {
				continue
			}
			sf := output.SchemaFlag{Name: names[0]}
			if doc, ok := f.(cli.DocGenerationFlag); ok {
				sf.Usage = doc.GetUsage()
			}
			sc.Flags = append(sc.Flags, sf)
		}
		env.Commands = append(env.Commands, sc)
	}
	for _, e := range terr.All() {
		env.Errors = append(env.Errors, output.SchemaError{
			Code:     e.Code(),
			ExitCode: e.ExitCode(),
			Hint:     e.Hint(),
		})
	}
	return output.EmitJSON(cmd.Root().Writer, env)
}
