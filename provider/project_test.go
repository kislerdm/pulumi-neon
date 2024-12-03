package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
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

	err = exec.Command("cd", path.Join(cwd, ".."), "&&", "make", "go_sdk").Run()
	if err != nil {
		panic(err)
	}

	sdkPath = path.Join(cwd, "..", "sdk")
}

func TestProject(t *testing.T) {
	if v, _ := strconv.ParseBool(os.Getenv("ACC_TEST")); !v {
		t.Skip("ACC_TEST is not set")
	}

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
			Dir: path.Join(cwd, "acc-test", "project", "default"),
			Secrets: map[string]string{
				"neon:api_key": token,
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.NoError(t, testQuery(stack.Outputs["connection_uri"].(string)))
				assert.NoError(t, testQuery(stack.Outputs["connection_uri_pooler"].(string)))
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
			Dir: path.Join(cwd, "acc-test", "project", "custom-name"),
			Secrets: map[string]string{
				"neon:api_key": token,
			},
			ExtraRuntimeValidation: func(t *testing.T, _ integration.RuntimeValidationStackInfo) {
				resp, err := client.ListProjects(nil, nil, &wantName, nil)
				assert.NoError(t, err)
				assert.Len(t, resp.Projects, 1)
			},
		})
	})
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

func TestProjectInOrganization(t *testing.T) {
	if v, _ := strconv.ParseBool(os.Getenv("ACC_TEST")); !v {
		t.Skip("ACC_TEST is not set")
	}

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
		Dir: path.Join(cwd, "acc-test", "project", "default-in-org"),
		Env: []string{"ORG_ID=" + orgID},
		Secrets: map[string]string{
			"neon:api_key": token,
		},
		ExtraRuntimeValidation: func(t *testing.T, _ integration.RuntimeValidationStackInfo) {
			resp, err := client.ListProjects(nil, nil, &wantName, &orgID)
			assert.NoError(t, err)
			assert.Len(t, resp.Projects, 1)
		},
	})
}
