import os
import subprocess
import re

{% if cookiecutter.git_init -%}
try:
    subprocess.run(['git', 'init'], check=True)
except subprocess.CalledProcessError as e:
    print(f"Error: Failed to initialize git repository. {e}")
    exit(1)
{%- endif %}

os.rename('.gitignore.tmp', '.gitignore')

author_id = re.sub(r"^@", '', '{{ cookiecutter.author_handle }}')
project_name = '{{ cookiecutter.project_name }}'
subprocess.run(['go', 'mod', 'init', f"github.com/{author_id}/{project_name}"], check=True, cwd="src")
