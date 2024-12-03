package main

import (
	"github.com/kislerdm/pulumi-neon/sdk/go/neon/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		p, err := provider.NewProject(ctx, "this", nil)
		if err != nil {
			return err
		}
		ctx.Export("connection_uri", p.Connection_uri)
		ctx.Export("connection_uri_pooler", p.Connection_uri_pooler)
		return nil
	})
}
