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
    + [Java](#java)
    + [YAML](#yaml)

## How to configure the provider

**Prerequisite:**

- Pulumi v3.134.1+

1. Sign up for Neon and [create an API token](https://api-docs.neon.tech/reference/authentication#neon-api-keys).
2. Export the token as the environment variable `NEON_API_KEY`.
3. Initiate a Pulumi project by running `pulumi new`.
4. Select the `template` in the dropdown.
5. Select one of the supported technologies:
    - `go`
    - `typescript`
    - `python`
    - `csharp`
    - `java-gradle` (note that the use `gradle` is a recommended way to build a java project)
    - `yaml`
6. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.

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

**Prerequisites**:

- dotnet 6+.
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`.

1. Create a pulumi project by running `pulumi new csharp`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Add the SDK as dependency:

```commandline
dotnet add package PulumiSdk.Neon
```

4. Edit the file `Program.cs`:

```csharp
using Pulumi;
using PulumiSdk.Neon.Resource;

return await Deployment.RunAsync(() =>
{
    new Project("myproject", new ProjectArgs
    {
        Name = "myProjectProvisionedWithPulumiDotnetSDK",
    });
});
```

5. (Optionally) Edit the project configuration to match the .Net version of your environment:

```xml

<Project Sdk="Microsoft.NET.Sdk">

    <PropertyGroup>
        <OutputType>Exe</OutputType>
        <!-- Configure your .Ner version here -->
        <TargetFramework>net8.0</TargetFramework>
        <Nullable>enable</Nullable>
    </PropertyGroup>

    <ItemGroup>
        <PackageReference Include="Pulumi" Version="3.*"/>
        <PackageReference Include="PulumiSdk.Neon" Version="0.*"/>
    </ItemGroup>

</Project>
```

6. Run `pulumi up -f`
7. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiDotnetSDK".

### Java

**Prerequisites**:

- Java 11+.
- Gradle >=8.12, <9.0.
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`.

1. Create a pulumi project by running `pulumi new java-gradle`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Add the SDK as dependency to the dependencies in the file `app/build.gradle`:

```gradle
dependencies {
    implementation ('com.pulumi:pulumi:(,1.0]')
    implementation ('com.pulumi:neon:(,1.0)')
}
```

4. Add the project configuration to the file `App.java`:

```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.neon.resource.Project;
import com.pulumi.neon.resource.ProjectArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(_ -> {
           new Project("myproject", ProjectArgs.builder()
                            .name("myProjectProvisionedWithPulumiJavaSDK")
                            .build());
        });
    }
}
```

5. Run `pulumi up -f`
6. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiJavaSDK".

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
      name: myProjectProvisionedWithPulumiYamlSDK
```

2. Run `pulumi up -f`
6. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiYamlSDK".
