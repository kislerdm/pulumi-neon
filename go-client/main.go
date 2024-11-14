package main

import (
	"log"

	"github.com/kislerdm/pulumi-neon/sdk/go/neon"
	"github.com/kislerdm/pulumi-neon/sdk/go/neon/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := neon.NewProvider(ctx, "neon", &neon.ProviderArgs{})
		if err != nil {
			log.Println(err)
		}
		p, err := provider.NewProject(ctx, "bar", &provider.ProjectArgs{})
		if err != nil {
			log.Println(err)
		}
		log.Println(p.Identifier)
		return err
	})
}
