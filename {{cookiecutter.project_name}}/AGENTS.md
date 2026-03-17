<!--
  TODO: we'll add content here for AI agents
-->

## Ticket Workflow

Tickets are managed with the `tk` CLI. Files live in `.tickets/`.

```bash
tk ls                    # List all tickets
tk ready                 # List tickets with all deps resolved (ready to start)
tk blocked               # List tickets blocked by unresolved deps
tk show <id>             # Show full ticket details
tk dep tree <id>         # Show dependency tree
tk start <id>            # Mark as in_progress
tk close <id>            # Mark as closed
tk add-note <id> "..."   # Append a note
```

## Build & Validation

All commands are run from the project root via `make`. Go source lives under
`src/`.

| Command          | Purpose                                                                         |
| ---------------- | ------------------------------------------------------------------------------- |
| `make build`     | Full build — runs fmt-check, vet, lint, test, clean, then compiles to `bin/app` |
| `make run`       | Run the server directly with `go run`                                           |
| `make test`      | Run all tests (`go test ./...`)                                                 |
| `make test-race` | Run tests with the race detector                                                |
| `make vet`       | Run `go vet ./...`                                                              |
| `make fmt`       | Format all Go source with `gofmt -w`                                            |
| `make fmt-check` | Fail if any files are not formatted                                             |
| `make lint`      | Run `golangci-lint` (depends on `vet`)                                          |
| `make clean`     | Remove build artifacts from `bin/`                                              |
| `make container` | Build Docker image (runs fmt-check, vet, lint, test first)                      |

### Validating Your Work

After any code change, **always run `make build`** from the project root before
considering the task complete. This single command enforces the full quality
gate:

1. `fmt-check` — code must be properly formatted
2. `vet` — no suspicious constructs
3. `lint` — no lint violations (golangci-lint)
4. `test` — all tests must pass
5. `clean` + compile — produces the binary

If `make build` fails at any step, fix the issue before proceeding. Do not
silence lint errors with `_ =` — handle them properly (log, return, or report
in tests).

For a quicker feedback loop during development, use `make test` or `make lint`
individually, but always finish with a full `make build`.
