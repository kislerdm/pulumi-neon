---
title: Neon Provider
meta_desc: Overview of the Pulumi Neon provider.
layout: package
---

# Neon provider

![logo](https://raw.githubusercontent.com/kislerdm/pulumi-neon/refs/heads/master/fig/logo.png)

The Pulumi provider to manage [Neon Platform](https://neon.tech/home) resources.

> Ship faster with Postgres
> The database you love, on a serverless platform designed to help you build reliable and scalable applications faster.

Find more about Neon [here](https://neon.tech/docs/introduction).

## Table of Contents

* [How to configure the provider](#how-to-configure-the-provider)
* [Example: how to provision a Neon Project](#example-how-to-provision-a-neon-project)
   + [Go](#go)
   + [Typescript](#typescript)
   + [Python](#python)
   + [C#](#c)
   + [YAML](#yaml)

## How to configure the provider

1. Sign up for Neon and [create an API token](https://api-docs.neon.tech/reference/authentication#neon-api-keys).
2. Export the token as the environment variable `NEON_API_KEY`.
3. Initiate a Pulumi project by running `pulumi new`.
4. Select the `template` in the dropdown.
5. Select one of the technologies supported by the provider:
    - `go`
    - `typescript`
    - `python`
    - `csharp`
    - `yaml`
6. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
7. Install the plugin by running ``

## Example: how to provision a Neon Project

**Prerequisite:** the [configuration steps](#how-to-configure-the-provider) above are completed.

### Go

1. Add the SDK as dependency:

```shell
go get "github.com/kislerdm/pulumi-sdk-neon"
```

2. Edit the file `main.go`:

```go
package main

import (
	"log"

	"github.com/kislerdm/pulumi-sdk-neon/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := resource.NewProject(ctx, "myproject", &resource.ProjectArgs{
			Name: pulumi.String("myproject"),
		}, pulumi.Protect(true))
		if err != nil {
			log.Println(err)
		}
		return err
	})
}
```

3. Run `pulumi up -f`
4. Examine the Neon console: it's expected to see a new project there.

### Typescript

1. Add the SDK as dependency:

```shell
npm install "@pulumi/neon"
```

2. Edit the file `index.ts`:

```typescript
import {Project, ProjectArgs} from "@pulumi/neon/provider/project";

new Project("myproject", {name: "myproject"} as ProjectArgs);
```

3. Run `pulumi up -f`
4. Examine the Neon console: it's expected to see a new project there.

### Python

1. Active venv:

```shell
source venv/bin/activate
```

2. Add the SDK as dependency:

```shell
pip install pulumi_neon
```

3. Edit the file `__main__.py`:

```python
from pulumi_neon.provider.project import Project, ProjectArgs

Project("myproject", ProjectArgs(name="myproject"))
```

4. Run `pulumi up -f`
5. Examine the Neon console: it's expected to see a new project there.

### C#

1. Add the SDK as dependency:

```commandline
dotnet add package PulumiSdk.Neon
```

2. Edit the file `Program.cs`:

```csharp
using Pulumi;
using PulumiSdk.Neon.Resource;

return await Deployment.RunAsync(() =>
{
    var project = new Project("myproject", new ProjectArgs
    {
        Name = "myproject",
    });
});
```

3. Run `pulumi up -f`
4. Examine the Neon console: it's expected to see a new project there.

### YAML

1. Edit the file `Pulumi.yaml` so it looks like this snippet:

```yaml
name: ##your project name##
runtime: yaml
description: ##your project description##
resources:
  project:
    type: neon:resource:Project
    properties:
      name: myproject
```

2. Run `pulumi up -f`
3. Examine the Neon console: it's expected to see a new project there.
