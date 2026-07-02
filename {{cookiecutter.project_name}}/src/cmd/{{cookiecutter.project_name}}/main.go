// Command {{ cookiecutter.project_name }} is {{ cookiecutter.project_short_description }}.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/version"
	"github.com/urfave/cli/v3"
)

func newApp() *cli.Command {
	return &cli.Command{
		Name:    "{{ cookiecutter.project_name }}",
		Usage:   "{{ cookiecutter.project_short_description }}",
		Version: version.Current(),
		Commands: []*cli.Command{
			{
				Name:      "greet",
				Usage:     "print a friendly greeting",
				ArgsUsage: "[name]",
				Action:    greetAction,
			},
		},
	}
}

// greetAction is a placeholder command: replace it with the tool's real
// commands. It greets the first positional argument, defaulting to "world".
func greetAction(_ context.Context, cmd *cli.Command) error {
	name := cmd.Args().First()
	if name == "" {
		name = "world"
	}
	if _, err := fmt.Fprintf(cmd.Root().Writer, "hello, %s\n", name); err != nil {
		return err
	}
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	if err := newApp().Run(context.Background(), os.Args); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}
