package main

import (
	"log"

	"github.com/kislerdm/pulumi-neon/sdk/go/neon/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := provider.NewProject(ctx, "bar", &provider.ProjectArgs{
			Name: pulumi.String("foo"),
		}, pulumi.Protect(true))
		if err != nil {
			log.Println(err)
		}
		return err
	})
}
