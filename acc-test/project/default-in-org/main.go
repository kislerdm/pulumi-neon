package main

import (
	"os"

	"github.com/kislerdm/pulumi-neon/sdk/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		p, err := provider.NewProject(ctx, "this", &provider.ProjectArgs{
			Name:   pulumi.String("pulumi-project-test-in-org"),
			Org_id: pulumi.String(os.Getenv("ORG_ID")),
		})
		ctx.Export("identifier", p.Identifier)
		ctx.Export("name", p.Name)
		ctx.Export("org_id", p.Org_id)
		ctx.Export("connection_uri", p.Connection_uri)
		ctx.Export("connection_uri_pooler", p.Connection_uri_pooler)
		ctx.Export("default_branch_name", p.Default_branch_name)
		ctx.Export("default_role_name", p.Default_role_name)
		ctx.Export("default_role_password", p.Default_role_password)
		ctx.Export("default_database_name", p.Default_database_name)
		ctx.Export("default_endpoint_host", p.Default_endpoint_host)
		ctx.Export("default_endpoint_host_pooler", p.Default_endpoint_host_pooler)
		return err
	})
}
