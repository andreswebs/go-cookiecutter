import os
import shutil
import subprocess

PROJECT_TYPE = "{{ cookiecutter.project_type }}"
CONTAINERIZE = "{{ cookiecutter.containerize }}"
CONTAINERIZED = PROJECT_TYPE == "api" or CONTAINERIZE == "true"
AUTHOR_ID = "{{ cookiecutter.author_id }}"
PROJECT_NAME = "{{ cookiecutter.project_name }}"


def rm(path):
    try:
        os.remove(path)
    except FileNotFoundError:
        pass


def rmtree(path):
    shutil.rmtree(path, ignore_errors=True)


# Prune the shape that was not selected.
if PROJECT_TYPE == "cli":
    rmtree("cmd/server")
    rmtree("internal/config")
    rmtree("internal/server")
    rmtree("api-spec")
else:  # api
    rmtree(os.path.join("cmd", PROJECT_NAME))
    # The CLI machinery (delegate, typed errors, output contract) is CLI-only;
    # the API keeps internal/secret, which is shape-agnostic.
    rmtree("internal/command")
    rmtree("internal/output")
    rmtree("internal/terr")
    # The signed binary-release pipeline is CLI-only.
    rm(".github/workflows/release.yml")
    rmtree(".github/actions")
    # The SBOM exclude file pairs with the release workflow, so it is dead here.
    rm(".syft.yaml")

# Drop container assets unless the project is containerized.
if not CONTAINERIZED:
    rm("Dockerfile")
    rm(".dockerignore")
    rm(".github/workflows/docker.yml")

{% if not cookiecutter.include_specs -%}
rmtree("docs")
{%- endif %}

os.rename(".gitignore.tmp", ".gitignore")

# Initialize the module, then tidy so go.mod/go.sum reflect the selected shape
# (the CLI pulls urfave/cli/v3; the API stays on the standard library).
subprocess.run(
    ["go", "mod", "init", f"github.com/{AUTHOR_ID}/{PROJECT_NAME}"],
    check=True,
)
subprocess.run(["go", "mod", "tidy"], check=True)

os.symlink("AGENTS.md", "CLAUDE.md")

{% if cookiecutter.git_init -%}
try:
    subprocess.run(["git", "init"], check=True)
except subprocess.CalledProcessError as e:
    print(f"Error: Failed to initialize git repository. {e}")
    raise SystemExit(1)
{%- endif %}
