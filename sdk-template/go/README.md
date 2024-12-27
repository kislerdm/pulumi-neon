# Go SDK for Pulumi Neon Provider

-----

<div align="center">
    ⭐ The project needs your support! Please leave a star and become a GitHub sponsor! ⭐
</div>

-----

The SDK to provision [Neon](https://neon.tech/) Postgres projects using the [Pulumi Neon provider](https://github.com/kislerdm/pulumi-neon).

## How to use

### Prerequisites

- Pulumi v3.134.1+
- Go 1.21+
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`

### Steps

1. Create a Pulumi project by running `pulumi new go`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Add the SDK as dependency:

```shell
go get "github.com/kislerdm/pulumi-sdk-neon"
```

4. Edit the file `main.go`:

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
			Name: pulumi.String("myProjectProvisionedWithPulumiGoSDK"),
		})
		if err != nil {
			log.Println(err)
		}
		return err
	})
}
```

5. Run `pulumi up -f`
6. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiGoSDK".

## How to contribute

Please raise [GitHub issue](https://github.com/kislerdm/pulumi-neon/issues/new) in case of proposals, or found bugs.
