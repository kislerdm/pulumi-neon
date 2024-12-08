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
	"maps"
	"reflect"
	"slices"
	"strings"

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
	ID                        string `pulumi:"identifier"`
	DefaultBranchName         string `pulumi:"default_branch_name"`
	DefaultRoleName           string `pulumi:"default_role_name"`
	DefaultRolePassword       string `pulumi:"default_role_password"`
	DefaultDatabaseName       string `pulumi:"default_database_name"`
	ConnectionURI             string `pulumi:"connection_uri"`
	ConnectionURIPooler       string `pulumi:"connection_uri_pooler"`
	DefaultEndpointHost       string `pulumi:"default_endpoint_host"`
	DefaultEndpointHostPooler string `pulumi:"default_endpoint_host_pooler"`
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
		output.DefaultDatabaseName = resp.DatabasesResponse.Databases[0].Name
		output.DefaultRoleName = resp.DatabasesResponse.Databases[0].OwnerName
		output.DefaultBranchName = resp.BranchResponse.Branch.Name

		for _, role := range resp.RolesResponse.Roles {
			if role.Name == output.DefaultRoleName {
				if role.Password != nil {
					output.DefaultRolePassword = *role.Password
				}
				break
			}
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

			normalizedState = ProjectState{
				ProjectArgs: normalizedInputs,
				ID:          canonicalID,
			}

			var defaultBranchID string
			for _, br := range respBranches.BranchesResponse.Branches {
				if br.Default {
					normalizedState.DefaultBranchName = br.Name
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
			normalizedState.DefaultDatabaseName = earliestCreatedDatabase.Name
			normalizedState.DefaultRoleName = earliestCreatedDatabase.OwnerName

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
			pooled := false
			respURI, err = c.GetConnectionURI(canonicalID, &defaultBranchID, &earliestCreatedEndpoint.ID,
				earliestCreatedDatabase.Name, earliestCreatedDatabase.OwnerName, &pooled)
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

func (pr Project) Diff(ctx context.Context, id string, olds ProjectState, news ProjectArgs) (diff p.DiffResponse,
	err error) {

	_, _, stateCloud, err := pr.Read(ctx, id, news, olds)
	if err != nil {
		return diff, fmt.Errorf("could not read the current real state: %w", err)
	}

	// define the deviation of the SaaS state from the pulumi state
	drift := projectDrift(stateCloud, olds)

	// define the change between the old and the new inputs
	inputChange := projectInputChange(olds, news)

	o := p.DiffResponse{
		DeleteBeforeReplace: drift.DeleteBeforeReplace || inputChange.DeleteBeforeReplace,
		HasChanges:          drift.HasChanges || inputChange.HasChanges,
		DetailedDiff:        drift.DetailedDiff,
	}
	maps.Copy(o.DetailedDiff, inputChange.DetailedDiff)

	return o, err
}

func projectInputChange(olds ProjectState, news ProjectArgs) p.DiffResponse {
	var o = p.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          false,
		DetailedDiff:        make(map[string]p.PropertyDiff),
	}

	// the project's name should be changed if the new input differs from the old pulumi input
	var changedName = news.Name != nil && !reflect.DeepEqual(news.Name, olds.inputState.Name)

	// the cloud state will not be changed if the name was removed from the manifest, or set to the empty string
	if news.Name == nil || (news.Name != nil && *news.Name == "") {
		changedName = false
	}

	if changedName {
		o.HasChanges = true
		o.DetailedDiff["name"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: true,
		}
	}

	// the project should be moved to the org id the new input differs from the old pulumi input
	var changedOrgID = !reflect.DeepEqual(news.OrgID, olds.inputState.OrgID)
	if changedOrgID {
		o.HasChanges = true
		o.DetailedDiff["org_id"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: true,
		}
	}

	return o
}

func projectDrift(read ProjectState, pulumi ProjectState) p.DiffResponse {

	var o = p.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          false,
		DetailedDiff:        make(map[string]p.PropertyDiff),
	}

	if !reflect.DeepEqual(read.Name, pulumi.Name) {
		o.HasChanges = true
		o.DetailedDiff["name"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	if !reflect.DeepEqual(read.OrgID, pulumi.OrgID) {
		o.HasChanges = true
		o.DetailedDiff["org_id"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	if read.DefaultBranchName != pulumi.DefaultBranchName {
		o.HasChanges = true
		o.DetailedDiff["default_branch_name"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	if read.DefaultRoleName != pulumi.DefaultRoleName {
		o.HasChanges = true
		o.DetailedDiff["default_role_name"] = p.PropertyDiff{
			Kind:      p.DeleteReplace,
			InputDiff: false,
		}
	}

	if read.DefaultDatabaseName != pulumi.DefaultDatabaseName {
		o.HasChanges = true
		o.DetailedDiff["default_database_name"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	if read.ConnectionURI != pulumi.ConnectionURI {
		o.HasChanges = true
		o.DetailedDiff["connection_uri"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	if read.ConnectionURIPooler != pulumi.ConnectionURIPooler {
		o.HasChanges = true
		o.DetailedDiff["connection_uri_pooler"] = p.PropertyDiff{
			Kind:      p.Update,
			InputDiff: false,
		}
	}

	return o
}
