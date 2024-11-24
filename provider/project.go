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
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Project struct{}

type ProjectArgs struct {
	Name  *string `pulumi:"name,optional"`
	OrgID *string `pulumi:"org_id,optional"`
}

type ProjectState struct {
	ProjectArgs
	ID string `pulumi:"identifier"`
}

func (p Project) Create(ctx context.Context, _ string, inputs ProjectArgs, preview bool) (
	id string, output ProjectState, err error) {
	c, err := sdk.NewClient(sdk.Config{Key: infer.GetConfig[*Config](ctx).APIKey})
	if err != nil {
		err = fmt.Errorf("could not init Neon Client: %w", err)
	}

	if !preview && err == nil {
		var resp sdk.CreatedProject
		resp, err = c.CreateProject(sdk.ProjectCreateRequest{
			Project: sdk.ProjectCreateRequestProject{
				Name:  inputs.Name,
				OrgID: inputs.OrgID,
				// 	TODO: add more attributes
			},
		})

		if err != nil {
			return id, output, err
		}

		id = resp.ProjectResponse.Project.ID
		output.ID = resp.ProjectResponse.Project.ID
		output.OrgID = resp.ProjectResponse.Project.OrgID
		output.Name = &resp.ProjectResponse.Project.Name
	}

	return id, output, err
}

func (p Project) Update(ctx context.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (
	output ProjectState, err error) {
	c, err := sdk.NewClient(sdk.Config{Key: infer.GetConfig[*Config](ctx).APIKey})
	if err != nil {
		err = fmt.Errorf("could not init Neon Client: %w", err)
	}

	if !preview && isProjectStateUpdated(olds, news) {
		var resp sdk.UpdateProjectRespObj
		resp, err = c.UpdateProject(id, sdk.ProjectUpdateRequest{
			Project: sdk.ProjectUpdateRequestProject{
				Name: news.Name,
				// 	TODO: add more attributes
			},
		})

		if err == nil {
			output.Name = &resp.ProjectResponse.Project.Name
		}
	}

	return output, err
}

func isProjectStateUpdated(olds ProjectState, news ProjectArgs) bool {
	// TODO: extend
	return olds.Name != news.Name
}
