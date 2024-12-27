# Python SDK for Pulumi Neon Provider

-----

<div align="center">
    ⭐ The project needs your support! Please leave a star and become a GitHub sponsor! ⭐
</div>

-----

The SDK to provision [Neon](https://neon.tech/) Postgres projects using the [Pulumi Neon provider](https://github.com/kislerdm/pulumi-neon).

## How to use

### Prerequisites

- Pulumi v3.134.1+
- Python 3.9+
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`

### Steps

1. Create a Pulumi project by running `pulumi new python`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Active venv:

```shell
source venv/bin/activate
```

4. Add the SDK as dependency:

```shell
pip install pulumi_neon
```

5. Edit the file `__main__.py`:

```python
from pulumi_neon.resource.project import Project, ProjectArgs

Project("myproject", ProjectArgs(name="myProjectProvisionedWithPulumiPythonSDK"))
```

6. Run `pulumi up -f`
7. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiPythonSDK".

## How to contribute

Please raise [GitHub issue](https://github.com/kislerdm/pulumi-neon/issues/new) in case of proposals, or found bugs.
