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
	"strings"

	sdk "github.com/kislerdm/neon-sdk-go"
)

type Project struct{}

type ProjectArgs struct {
	Name                *string `pulumi:"name,optional"`
	OrgID               *string `pulumi:"org_id,optional"`
	DefaultBranchName   *string `pulumi:"default_branch_name,optional"`
	DefaultRoleName     *string `pulumi:"default_role_name,optional"`
	DefaultDatabaseName *string `pulumi:"default_database_name,optional"`
}

type ProjectState struct {
	ProjectArgs
	ID                        string  `pulumi:"identifier"`
	DefaultRolePassword       *string `pulumi:"default_role_password"`
	ConnectionURI             string  `pulumi:"connection_uri"`
	ConnectionURIPooler       string  `pulumi:"connection_uri_pooler"`
	DefaultEndpointHost       string  `pulumi:"default_endpoint_host"`
	DefaultEndpointHostPooler string  `pulumi:"default_endpoint_host_pooler"`
}

func (p Project) Create(ctx context.Context, _ string, inputs ProjectArgs, preview bool) (
	id string, output ProjectState, err error) {
	c, err := NewSDKClient(ctx)

	if !preview && err == nil {
		var resp sdk.CreatedProject
		resp, err = c.CreateProject(sdk.ProjectCreateRequest{
			Project: sdk.ProjectCreateRequestProject{
				Branch: &sdk.ProjectCreateRequestProjectBranch{
					DatabaseName: inputs.DefaultDatabaseName,
					Name:         inputs.DefaultBranchName,
					RoleName:     inputs.DefaultRoleName,
				},
				Name:  inputs.Name,
				OrgID: inputs.OrgID,
			},
		})

		if err != nil {
			return id, output, err
		}

		id = resp.ProjectResponse.Project.ID
		output.ID = resp.ProjectResponse.Project.ID
		output.OrgID = resp.ProjectResponse.Project.OrgID
		output.Name = &resp.ProjectResponse.Project.Name
		output.DefaultDatabaseName = &resp.DatabasesResponse.Databases[0].Name
		output.DefaultRoleName = &resp.DatabasesResponse.Databases[0].OwnerName
		output.DefaultBranchName = &resp.BranchResponse.Branch.Name

		var pass *string
		for _, role := range resp.RolesResponse.Roles {
			if output.DefaultRoleName != nil && role.Name == *output.DefaultRoleName {
				pass = role.Password
				break
			}
		}

		output.DefaultRolePassword = pass
		host := resp.EndpointsResponse.Endpoints[0].Host
		output.DefaultEndpointHost = host
		output.DefaultEndpointHostPooler = newHostPooler(host)
		output.ConnectionURI = resp.ConnectionURIs[0].ConnectionURI
		output.ConnectionURIPooler = newURIPooler(output.ConnectionURI)
	}

	return id, output, err
}

func newHostPooler(host string) string {
	const poolerSuffix = "-pooler"
	els := strings.SplitN(host, ".", 2)
	return fmt.Sprintf("%s.%s", els[0]+poolerSuffix, els[1])
}

func newURIPooler(uri string) string {
	els := strings.SplitN(uri, "@", 2)
	suffParts := strings.SplitN(els[1], "/", 2)
	return fmt.Sprintf("%s@%s/%s", els[0], newHostPooler(suffParts[0]), suffParts[1])
}

func (p Project) Update(ctx context.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (
	output ProjectState, err error) {
	c, err := NewSDKClient(ctx)

	if err == nil && !preview && isProjectStateUpdated(olds, news) {
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

func (p Project) Read(ctx context.Context, id string, _ ProjectArgs, _ ProjectState) (
	canonicalID string, normalizedInputs ProjectArgs, normalizedState ProjectState, err error) {
	c, err := NewSDKClient(ctx)
	if err == nil {
		var resp sdk.ProjectResponse
		resp, err = c.GetProject(id)
		if err == nil {
			canonicalID = resp.Project.ID
			normalizedInputs.Name = &resp.Project.Name
			normalizedInputs.OrgID = resp.Project.OrgID
			normalizedState = ProjectState{
				ProjectArgs: normalizedInputs,
				ID:          resp.Project.ID,
			}
		}
	}

	return canonicalID, normalizedInputs, normalizedState, err
}

func (p Project) Delete(ctx context.Context, id string, _ ProjectState) error {
	c, err := NewSDKClient(ctx)
	if err == nil {
		_, err = c.DeleteProject(id)
	}
	return err
}
