package project_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/jackc/pgx/v5"
	sdk "github.com/kislerdm/neon-sdk-go"
	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/fsutil"
	"github.com/stretchr/testify/assert"
)

var (
	cwd     string
	sdkPath string
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	sdkPath = path.Join(cwd, "..", "..", "sdk")
	_, err = os.ReadDir(sdkPath)
	if err != nil {
		panic(fmt.Sprintf("SDK not found in %s: %v", sdkPath, err))
	}
}

func verifyProjectOutputs(t *testing.T, stack integration.RuntimeValidationStackInfo, client *sdk.Client,
	wantProjectID string) {
	t.Helper()

	assert.Equal(t, wantProjectID, stack.Outputs["identifier"].(string),
		"project ID should match resource identifier")

	assert.NotEmpty(t, stack.Outputs["default_branch_name"].(string),
		"project default_branch_name should be not empty")

	assert.NoError(t, testQuery(stack.Outputs["connection_uri"].(string)),
		"project should include default database with valid connection URI")
	assert.NoError(t, testQuery(stack.Outputs["connection_uri_pooler"].(string)),
		"project should include default database with valid connection URI, pooling active")

	uri := newURI(stack.Outputs["default_role_name"].(string),
		stack.Outputs["default_role_password"].(string),
		stack.Outputs["default_database_name"].(string),
		stack.Outputs["default_endpoint_host"].(string),
	)
	assert.NoError(t, testQuery(uri),
		"project should include default database, role and endpoint")

	uriPooling := newURI(stack.Outputs["default_role_name"].(string),
		stack.Outputs["default_role_password"].(string),
		stack.Outputs["default_database_name"].(string),
		stack.Outputs["default_endpoint_host"].(string),
	)
	assert.NoError(t, testQuery(uriPooling),
		"project should include default database, role and endpoint, pooling active")
}

func newURI(roleName string, rolePassword string, dbName string, host string) string {
	const sslMode = "?sslmode=require"
	return "postgres://" + roleName + ":" + rolePassword + "@" + host + "/" + dbName + sslMode
}

func testQuery(uri string) error {
	conn, err := pgx.Connect(context.TODO(), uri)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close(context.Background()) }()

	wantVal := 1
	r, _ := conn.Query(context.TODO(), fmt.Sprintf("select %d as val;", wantVal))
	defer func() { r.Close() }()

	vals, err := pgx.CollectRows(r, func(row pgx.CollectableRow) (int, error) {
		var val int
		err := row.Scan(&val)
		return val, err
	})

	if err == nil {
		if len(vals) != 1 || vals[0] != wantVal {
			err = fmt.Errorf("expected to return %d as the query result", wantVal)
		}
	}

	return err
}

func TestProject(t *testing.T) {
	token := os.Getenv("NEON_API_KEY")
	if token == "" {
		t.Fatal("neon API key must be set as env variable NEON_API_KEY for integration tests")
	}

	client, err := sdk.NewClient(sdk.Config{Key: token})
	assert.NoError(t, err)

	t.Run("default config", func(t *testing.T) {
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Quick:       true,
			SkipRefresh: true,
			PrepareProject: func(projinfo *engine.Projinfo) error {
				return fsutil.CopyFile(projinfo.Root, sdkPath, nil)
			},
			Dir: path.Join(cwd, "default"),
			Secrets: map[string]string{
				"neon:api_key": token,
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				gotName := stack.Outputs["name"].(string)
				assert.NotEmpty(t, gotName, "project name should be not empty")

				resp, err := client.ListProjects(nil, nil, &gotName, nil)
				assert.NoError(t, err)
				assert.Len(t, resp.Projects, 1, "defined project should be created")

				verifyProjectOutputs(t, stack, client, resp.ProjectsResponse.Projects[0].ID)

				assert.Nil(t, stack.Outputs["org_id"], "project org_id should be empty")
			},
		})
	})

	t.Run("custom project name pulumi-project-test-custom-name", func(t *testing.T) {
		wantName := "pulumi-project-test-custom-name"
		integration.ProgramTest(t, &integration.ProgramTestOptions{
			Quick:       true,
			SkipRefresh: true,
			PrepareProject: func(projinfo *engine.Projinfo) error {
				return fsutil.CopyFile(projinfo.Root, sdkPath, nil)
			},
			Dir: path.Join(cwd, "custom-name"),
			Secrets: map[string]string{
				"neon:api_key": token,
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				resp, err := client.ListProjects(nil, nil, &wantName, nil)
				assert.NoError(t, err)
				assert.Len(t, resp.Projects, 1)

				assert.Equal(t, wantName, stack.Outputs["name"].(string),
					"project should have configured name")

				verifyProjectOutputs(t, stack, client, resp.ProjectsResponse.Projects[0].ID)

				assert.Nil(t, stack.Outputs["org_id"], "project org_id should be empty")
			},
		})
	})
}

func TestProjectInOrganization(t *testing.T) {
	orgID := os.Getenv("ORG_ID")
	if orgID == "" {
		t.Skip("ORG_ID is not set")
	}

	token := os.Getenv("NEON_API_KEY")
	if token == "" {
		t.Fatal("neon API key must be set as env variable NEON_API_KEY for integration tests")
	}

	client, err := sdk.NewClient(sdk.Config{Key: token})
	assert.NoError(t, err)

	wantName := "pulumi-project-test-in-org"

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:       true,
		SkipRefresh: true,
		PrepareProject: func(projinfo *engine.Projinfo) error {
			return fsutil.CopyFile(projinfo.Root, sdkPath, nil)
		},
		Dir: path.Join(cwd, "default-in-org"),
		Env: []string{"ORG_ID=" + orgID},
		Secrets: map[string]string{
			"neon:api_key": token,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			resp, err := client.ListProjects(nil, nil, &wantName, &orgID)
			assert.NoError(t, err)
			assert.Len(t, resp.Projects, 1)

			verifyProjectOutputs(t, stack, client, resp.ProjectsResponse.Projects[0].ID)

			assert.Equal(t, orgID, stack.Outputs["org_id"])
		},
	})
}
