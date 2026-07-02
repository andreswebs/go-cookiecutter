# go-cookiecutter

[Cookiecutter](https://www.cookiecutter.io/) template to generate a Go project
in one of two shapes, a command-line tool or an HTTP API server, with a modern,
signed release pipeline.

## What it generates

Both shapes share:

- A `src/` module with an `internal/version` package whose value is stamped in at
  build time (`-ldflags -X`), a `Makefile` quality gate (`fmt-check`, `vet`,
  `lint`, `test`), and a `golangci-lint` v2 config.
- A `CI` workflow that runs the full `make build` gate on push and pull request.
- Third-party GitHub Actions pinned by commit SHA, kept current by Dependabot.
- An `AGENTS.md` (symlinked to `CLAUDE.md`) describing the build and release flow.

The `cli` shape adds:

- A `urfave/cli/v3` entrypoint at `src/cmd/<project_name>/`.
- A `Release` workflow triggered by `v*.*.*` tags that cross-compiles every
  platform, keyless-signs the checksums and an SPDX SBOM with
  [cosign](https://docs.sigstore.dev/), attaches SLSA build provenance with
  `actions/attest-build-provenance`, and publishes the release with verification
  instructions in the notes.

The `api` shape adds:

- An HTTP server at `src/cmd/server/` with liveness and readiness probes,
  graceful shutdown, and an `api-spec/openapi.yaml`.
- A container-publish workflow that pushes the image to Docker Hub, then
  keyless-signs it with cosign, attaches SLSA build provenance, and attaches a
  cosign-signed SPDX SBOM.

Containers are always generated for the `api` shape and are optional for the
`cli` shape (see `containerize`).

## Variables

| Variable                    | Description                                                                 |
| --------------------------- | --------------------------------------------------------------------------- |
| `project_name_full`         | Human-readable project name (e.g. `My Tool`).                               |
| `project_name`              | Project slug; used for the repo, binary, and module path. Derived from the full name. |
| `project_short_description` | One-line description.                                                       |
| `project_type`              | `cli` or `api`. Selects the project shape.                                  |
| `containerize`              | `true`/`false`. Generate a Dockerfile and signed container workflow. Forced on for `api`. |
| `go_version`                | Go minor version for the Docker build image (e.g. `1.26`).                  |
| `author_name`               | Author full name.                                                           |
| `author_handle`             | Author handle (e.g. `@you`).                                                |
| `author_id`                 | Handle without the `@`; used in the module path and image name. Derived.    |
| `author_link`               | Author profile URL.                                                         |
| `git_init`                  | Run `git init` after generating.                                            |

## Pre-requisites

Install [cookiecutter](https://cookiecutter.readthedocs.io/en/stable) and a Go
toolchain (the post-generation hook runs `go mod init` and `go mod tidy`).

We recommend [uv](https://docs.astral.sh/uv/) to manage the Python environment:

```sh
brew install uv
uv venv
source .venv/bin/activate
uv pip install cookiecutter
```

## Run

```sh
cookiecutter gh:andreswebs/go-cookiecutter
```

Render non-interactively with defaults (useful for testing):

```sh
cookiecutter gh:andreswebs/go-cookiecutter --no-input \
  project_name_full="My Tool" project_type=cli
```

## Container publishing

The `api` shape (and a containerized `cli`) publishes to Docker Hub. Set these
repository secrets before the workflow runs:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

## Authors

**Andre Silva** - [@andreswebs](https://github.com/andreswebs)

## License

This project is licensed under the [Unlicense](UNLICENSE).
