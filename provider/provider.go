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
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

var Version string

const (
	Name = "neon"
)

func Provider() p.Provider {
	o := infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			Description:       "Pulumi Neon Provider",
			DisplayName:       Name,
			Keywords:          []string{"pulumi", Name, "category/cloud"},
			Homepage:          "https://github.com/kislerdm/pulumi-neon",
			Repository:        "https://github.com/kislerdm/pulumi-neon",
			Publisher:         "kislerdm",
			PluginDownloadURL: "https://github.com/kislerdm/pulumi-neon/releases/download/v${VERSION}",
			LogoURL:           "https://raw.githubusercontent.com/kislerdm/pulumi-neon/refs/heads/main/fig/logo.svg",
			License:           "Apache-2.0",
		},
		Resources: []infer.InferredResource{
			infer.Resource[*Project, ProjectArgs, ProjectState](),
		},
		Config: infer.Config[*Config](),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"neon": "index",
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
