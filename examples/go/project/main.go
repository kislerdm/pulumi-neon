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
		}, pulumi.Protect(true))
		if err != nil {
			log.Println(err)
		}
		return err
	})
}
