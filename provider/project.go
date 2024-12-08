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
	"slices"
	"strings"

	"github.com/davecgh/go-spew/spew"
	sdk "github.com/kislerdm/neon-sdk-go"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Project struct{}

type ProjectArgs struct {
	Name  *string `pulumi:"name,optional"`
	OrgID *string `pulumi:"org_id,optional"`
}

func (pr *ProjectArgs) Annotate(a infer.Annotator) {
	a.Describe(&pr.Name, "Neon project name.")
	a.Describe(&pr.OrgID, "Neon Org ID.")
}

type ProjectState struct {
	ProjectArgs
	inputState                ProjectArgs
	ID                        string  `pulumi:"identifier"`
	DefaultBranchName         *string `pulumi:"default_branch_name,optional"`
	DefaultRoleName           *string `pulumi:"default_role_name,optional"`
	DefaultRolePassword       string  `pulumi:"default_role_password,secret"`
	DefaultDatabaseName       *string `pulumi:"default_database_name,optional"`
	ConnectionURI             string  `pulumi:"connection_uri,secret"`
	ConnectionURIPooler       string  `pulumi:"connection_uri_pooler,secret"`
	DefaultEndpointHost       string  `pulumi:"default_endpoint_host"`
	DefaultEndpointHostPooler string  `pulumi:"default_endpoint_host_pooler"`
}

func (pr *ProjectState) Annotate(a infer.Annotator) {
	a.Describe(&pr.ID, "Project ID.")
	a.Describe(&pr.Name, "Neon project name.")
	a.Describe(&pr.OrgID, "Neon Org ID.")
	a.Describe(&pr.DefaultBranchName, "Neon default branch's name.")
	a.Describe(&pr.DefaultDatabaseName, "Neon default database's name.")
	a.Describe(&pr.DefaultRoleName, "Neon default role's name.")
	a.Describe(&pr.DefaultRolePassword, "Neon default role's password.")
	a.Describe(&pr.ConnectionURI, "URI to connect to the default database using the default endpoint.")
	a.Describe(&pr.ConnectionURIPooler,
		"URI to connect to the default database using the default endpoint in the pooler mode.")
	a.Describe(&pr.DefaultEndpointHost, "The default endpoint's host.")
	a.Describe(&pr.DefaultEndpointHostPooler, "The default endpoint's host with the pooler mode active.")
}

func (pr Project) Create(ctx context.Context, _ string, inputs ProjectArgs, preview bool) (
	id string, output ProjectState, err error) {
	c, err := NewSDKClient(ctx)

	if !preview && err == nil {
		var resp sdk.CreatedProject
		resp, err = c.CreateProject(sdk.ProjectCreateRequest{
			Project: sdk.ProjectCreateRequestProject{
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

		if pass != nil {
			output.DefaultRolePassword = *pass
		}
		host := resp.EndpointsResponse.Endpoints[0].Host
		output.DefaultEndpointHost = host
		output.DefaultEndpointHostPooler = newHostPooler(host)
		output.ConnectionURI = resp.ConnectionURIs[0].ConnectionURI
		output.ConnectionURIPooler = newURIPooler(output.ConnectionURI)

		// preserve the inputs
		output.inputState = inputs
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

func (pr Project) Update(ctx context.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (
	output ProjectState, err error) {
	c, err := NewSDKClient(ctx)
	if err != nil {
		return output, err
	}

	if !preview {
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

	} else {
		_, _, output, err = pr.Read(ctx, id, news, olds)
	}

	return output, err
}

func (pr Project) Read(ctx context.Context, id string, _ ProjectArgs, _ ProjectState) (
	canonicalID string, normalizedInputs ProjectArgs, normalizedState ProjectState, err error) {
	c, err := NewSDKClient(ctx)
	if err == nil {
		var resp sdk.ProjectResponse
		resp, err = c.GetProject(id)
		if err == nil {
			canonicalID = resp.Project.ID
			normalizedInputs.Name = &resp.Project.Name
			normalizedInputs.OrgID = resp.Project.OrgID

			var respBranches sdk.ListProjectBranchesRespObj
			respBranches, err = c.ListProjectBranches(canonicalID, nil)
			if err != nil {
				return canonicalID, normalizedInputs, normalizedState, err
			}

			var defaultBranchID string
			for _, br := range respBranches.BranchesResponse.Branches {
				if br.Default {
					normalizedState.DefaultBranchName = &br.Name
					defaultBranchID = br.ID
					break
				}
			}

			var respDB sdk.DatabasesResponse
			respDB, err = c.ListProjectBranchDatabases(canonicalID, defaultBranchID)
			if err != nil {
				return "", ProjectArgs{}, ProjectState{}, err
			}

			slices.SortStableFunc(respDB.Databases, func(a, b sdk.Database) int {
				return a.CreatedAt.Compare(b.CreatedAt)
			})

			// the earliest created database is assumed default
			earliestCreatedDatabase := respDB.Databases[0]
			normalizedState.DefaultDatabaseName = &earliestCreatedDatabase.Name
			normalizedState.DefaultRoleName = &earliestCreatedDatabase.OwnerName

			normalizedState = ProjectState{
				ProjectArgs: normalizedInputs,
				ID:          canonicalID,
			}

			var respPass sdk.RolePasswordResponse
			respPass, err = c.GetProjectBranchRolePassword(canonicalID, defaultBranchID, earliestCreatedDatabase.OwnerName)
			if err != nil {
				return "", ProjectArgs{}, ProjectState{}, err
			}

			normalizedState.DefaultRolePassword = respPass.Password

			var respEndpoints sdk.EndpointsResponse
			respEndpoints, err = c.ListProjectBranchEndpoints(canonicalID, defaultBranchID)
			if err != nil {
				return "", ProjectArgs{}, ProjectState{}, err
			}
			// the earliest created endpoint is assumed default
			slices.SortStableFunc(respEndpoints.Endpoints, func(a, b sdk.Endpoint) int {
				return a.CreatedAt.Compare(b.CreatedAt)
			})

			earliestCreatedEndpoint := respEndpoints.Endpoints[0]
			normalizedState.DefaultEndpointHost = earliestCreatedEndpoint.Host
			normalizedState.DefaultEndpointHostPooler = newHostPooler(earliestCreatedEndpoint.Host)

			var respURI sdk.ConnectionURIResponse
			respURI, err = c.GetConnectionURI(canonicalID, &defaultBranchID, &earliestCreatedEndpoint.ID,
				earliestCreatedDatabase.Name, earliestCreatedDatabase.OwnerName, nil)
			if err != nil {
				return "", ProjectArgs{}, ProjectState{}, err
			}
			normalizedState.ConnectionURI = respURI.URI
			normalizedState.ConnectionURIPooler = newURIPooler(respURI.URI)
		}
	}

	return canonicalID, normalizedInputs, normalizedState, err
}

func (pr Project) Delete(ctx context.Context, id string, _ ProjectState) error {
	c, err := NewSDKClient(ctx)
	if err == nil {
		_, err = c.DeleteProject(id)
	}
	return err
}

func (pr Project) Diff(_ context.Context, _ string, olds ProjectState, news ProjectArgs) (diff p.DiffResponse,
	err error) {
	var (
		isDiff    bool
		isDestroy bool
	)

	spew.Dump("olds", olds)
	spew.Dump("news", news)

	isDiff = olds.Name != news.Name
	isDiff = isDiff || isDestroy

	diff = p.DiffResponse{
		DeleteBeforeReplace: isDestroy,
		HasChanges:          isDiff,
	}

	return diff, err
}
