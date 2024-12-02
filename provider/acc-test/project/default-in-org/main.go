package main

import (
	"os"

	"github.com/kislerdm/pulumi-neon/sdk/go/neon/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := provider.NewProject(ctx, "this", &provider.ProjectArgs{
			Name:   pulumi.String("pulumi-project-test-in-org"),
			Org_id: pulumi.String(os.Getenv("ORG_ID")),
		})
		return err
	})
}
