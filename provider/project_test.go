package provider

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"

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

	_, err := sdk.NewClient(sdk.Config{Key: token})
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
		})
	})
}
