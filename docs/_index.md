---
title: Neon Provider
meta_desc: Overview of the Pulumi Neon provider.
layout: package
---

# Neon provider

![logo](https://raw.githubusercontent.com/kislerdm/pulumi-neon/refs/heads/main/fig/logo.svg)

The Pulumi provider to manage [Neon Platform](https://neon.tech/home) resources.

> Ship faster with Postgres
> The database you love, on a serverless platform designed to help you build reliable and scalable applications faster.

Find more about Neon [here](https://neon.tech/docs/introduction).

## How to configure the provider

1. Initiate a Pulumi project by running `pulumi new`.
2. Select the `template` in the dropdown.
3. Select one of the technologies supported by the provider:
    - `go`
    - `python`
    - `typescript`
4. Sign up for Neon and [create an API token](https://api-docs.neon.tech/reference/authentication#neon-api-keys).
5. Export the token as the environment variable `NEON_API_KEY`.
6. (Optionally) Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.

## How to provision a Neon Project

**Prerequisite:** the [configuration steps](#how-to-configure-the-provider) above are completed.

### Go

1. Add the SDK as dependency:

```shell
go get "github.com/kislerdm/pulumi-neon/sdk"
```

2. Edit the file `main.go`:

```go
package main

import (
	"log"

	"github.com/kislerdm/pulumi-neon/sdk/go/neon/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := provider.NewProject(ctx, "myproject", &provider.ProjectArgs{
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
