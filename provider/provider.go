// Copyright 2016-2023, Pulumi Corporation.
// Copyright 2024, Dmitry Kisler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"

	sdk "github.com/kislerdm/neon-sdk-go"
	"github.com/kislerdm/pulumi-neon/provider/telemetry"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	javaGen "github.com/pulumi/pulumi-java/pkg/codegen/java"
	csharpGen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	goGen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsGen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pythonGen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

var Version = "testAcc"

const (
	Name = "neon"
)

func Provider() p.Provider {
	const repository = "github.com/kislerdm/pulumi-neon"
	o := infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			Description:       "Pulumi Neon Provider",
			DisplayName:       "Neon Provider",
			Keywords:          []string{"pulumi", Name, "category/database"},
			Homepage:          "https://" + repository,
			Repository:        "https://" + repository,
			Publisher:         "Dmitry Kisler",
			PluginDownloadURL: "https://" + repository + "/releases/download/v${VERSION}",
			LogoURL:           "https://raw.githubusercontent.com/kislerdm/pulumi-neon/refs/heads/master/fig/logo.png",
			License:           "Apache-2.0",
			LanguageMap: map[string]any{
				"go": goGen.GoPackageInfo{
					ImportBasePath: "github.com/kislerdm/pulumi-sdk-neon",
					PackageImportAliases: map[string]string{
						"github.com/kislerdm/pulumi-sdk-neon/provider": "pulumiNeon",
					},
					PulumiSDKVersion:               3,
					GenerateResourceContainerTypes: true,
					LiftSingleValueMethodReturns:   true,
					RespectSchemaVersion:           true,
				},
				"nodejs": nodejsGen.NodePackageInfo{
					PackageName:                  "@dkisler/pulumi-neon",
					LiftSingleValueMethodReturns: true,
					RespectSchemaVersion:         true,
				},
				"python": pythonGen.PackageInfo{
					Requires: map[string]string{
						"pulumi": ">=3.0.0,<4.0.0",
					},
					LiftSingleValueMethodReturns: true,
					RespectSchemaVersion:         true,
				},
				"csharp": csharpGen.CSharpPackageInfo{
					PackageReferences: map[string]string{
						"Pulumi": "3.*",
					},
					Compatibility:                "",
					DictionaryConstructors:       false,
					LiftSingleValueMethodReturns: true,
					RootNamespace:                "PulumiSdk",
					RespectSchemaVersion:         true,
				},
				"java": javaGen.PackageInfo{
					BuildFiles: "gradle",
					Dependencies: map[string]string{
						"com.pulumi:pulumi": "(0,1.0]",
					},
				},
			},
		},
		Resources: []infer.InferredResource{
			infer.Resource[Project, ProjectArgs, ProjectState](),
		},
		Config: infer.Config[*Config](),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "resource",
		},
	})

	return o
}

type Config struct {
	APIKey string `pulumi:"api_key"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.APIKey, "Neon API token.")
	a.SetDefault(&c.APIKey, nil, "NEON_API_KEY")
}

func NewSDKClient(ctx context.Context) (*sdk.Client, error) {
	c, err := sdk.NewClient(sdk.Config{
		Key:        infer.GetConfig[*Config](ctx).APIKey,
		HTTPClient: telemetry.NewHTTPClient("kislerdm/"+Name, Version),
	})
	if err != nil {
		err = fmt.Errorf("could not init Neon Client: %w", err)
	}
	return c, err
}
