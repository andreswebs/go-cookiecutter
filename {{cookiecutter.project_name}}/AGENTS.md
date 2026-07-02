# AI agent instructions

`{{ cookiecutter.project_name }}` is {{ cookiecutter.project_short_description }}.
Go source lives under `src/`; all commands are run from the project root via
`make`.

## Build & validation

| Command          | Purpose                                                                    |
| ---------------- | -------------------------------------------------------------------------- |
| `make build`     | Full build: fmt-check, vet, lint, test, clean, then compile to `bin/`      |
| `make run`       | Run {% if cookiecutter.project_type == 'cli' %}the CLI{% else %}the server{% endif %} directly with `go run`                              |
| `make test`      | Run all tests (`go test ./...`)                                            |
| `make test-race` | Run tests with the race detector                                           |
| `make vet`       | Run `go vet ./...`                                                         |
| `make fmt`       | Format all Go source with `gofmt -w`                                       |
| `make fmt-check` | Fail if any files are not formatted                                        |
| `make lint`      | Run `golangci-lint`                                                        |
| `make dist`      | Cross-compile every platform and package archives + checksums into `dist/` |
| `make clean`     | Remove build artifacts                                                     |
{%- if cookiecutter.project_type == 'api' or cookiecutter.containerize == 'true' %}
| `make container` | Build the Docker image (runs the quality gate first)                       |
{%- endif %}

### Validating your work

After any code change, **always run `make build`** from the project root before
considering the task complete. It enforces the full quality gate:

1. `fmt-check` — code must be properly formatted
2. `vet` — no suspicious constructs
3. `lint` — no lint violations (golangci-lint)
4. `test` — all tests must pass
5. `clean` + compile — produces the binary

If `make build` fails at any step, fix the issue before proceeding. Do not
silence lint errors with `_ =` — handle them properly (return, log, or assert in
a test).

## Releases & supply-chain security

{% if cookiecutter.project_type == 'cli' -%}
Pushing a `v*.*.*` tag runs `.github/workflows/release.yml`, which cross-compiles
every platform, keyless-signs the checksums and SBOM with cosign, attaches SLSA
build provenance with `actions/attest-build-provenance`, and publishes the
GitHub release with verification instructions in the notes.
{%- else -%}
`.github/workflows/docker.yml` builds and pushes the image to Docker Hub, then
keyless-signs it with cosign, attaches SLSA build provenance with
`actions/attest-build-provenance`, and attaches a cosign-signed SPDX SBOM. It
needs the `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` repository secrets.
{%- endif %}

Third-party actions are pinned by commit SHA and kept current by Dependabot.
