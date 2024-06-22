package cli

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/DefangLabs/defang/src/pkg/cli/compose"
	"github.com/DefangLabs/defang/src/pkg/term"
	"github.com/sirupsen/logrus"
)

type ComposeLoader = compose.Loader

func TestLoadCompose(t *testing.T) {
	DoVerbose = true
	term.SetDebug(true)

	t.Run("no project name defaults to parent directory name", func(t *testing.T) {
		loader := ComposeLoader{"../../tests/noprojname/compose.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "noprojname" { // Use the parent directory name as project name
			t.Errorf("LoadCompose() failed: expected project name tenant-id, got %q", p.Name)
		}
	})

	t.Run("no project name defaults to fancy parent directory name", func(t *testing.T) {
		loader := ComposeLoader{"../../tests/Fancy-Proj_Dir/compose.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "fancy-proj_dir" { // Use the parent directory name as project name
			t.Errorf("LoadCompose() failed: expected project name tenant-id, got %q", p.Name)
		}
	})

	t.Run("use project name in compose file", func(t *testing.T) {
		loader := ComposeLoader{"../../tests/testproj/compose.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "tests" {
			t.Errorf("LoadCompose() failed: expected project name, got %q", p.Name)
		}
	})

	t.Run("COMPOSE_PROJECT_NAME env var should override project name", func(t *testing.T) {
		t.Setenv("COMPOSE_PROJECT_NAME", "overridename")
		loader := ComposeLoader{"../../tests/testproj/compose.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "overridename" {
			t.Errorf("LoadCompose() failed: expected project name to be overwritten by env var, got %q", p.Name)
		}
	})

	t.Run("use project name should not be overriden by tenantID", func(t *testing.T) {
		loader := ComposeLoader{"../../tests/testproj/compose.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "tests" {
			t.Errorf("LoadCompose() failed: expected project name tests, got %q", p.Name)
		}
	})

	t.Run("load starting from a sub directory", func(t *testing.T) {
		cwd, _ := os.Getwd()

		// setup
		setup := func() {
			os.MkdirAll("../../tests/alttestproj/subdir/subdir2", 0755)
			os.Chdir("../../tests/alttestproj/subdir/subdir2")
		}

		//teardown
		teardown := func() {
			os.Chdir(cwd)
			os.RemoveAll("../../tests/alttestproj/subdir")
		}

		setup()
		defer teardown()

		// execute test
		loader := ComposeLoader{}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "tests" {
			t.Errorf("LoadCompose() failed: expected project name, got %q", p.Name)
		}
	})

	t.Run("load alternative compose file", func(t *testing.T) {
		loader := ComposeLoader{"../../tests/alttestproj/altcomp.yaml"}
		p, err := loader.LoadCompose(context.Background())
		if err != nil {
			t.Fatalf("LoadCompose() failed: %v", err)
		}
		if p.Name != "tests" {
			t.Errorf("LoadCompose() failed: expected project name, got %q", p.Name)
		}
	})
}

func TestComposeGoNoDoubleWarningLog(t *testing.T) {
	var warnings bytes.Buffer
	logrus.SetOutput(&warnings)

	loader := compose.Loader{"../../tests/compose-go-warn/compose.yaml"}
	_, err := loader.LoadCompose(context.Background())
	if err != nil {
		t.Fatalf("LoadCompose() failed: %v", err)
	}

	if bytes.Count(warnings.Bytes(), []byte(`\"yes\" for boolean is not supported by YAML 1.2`)) != 1 {
		t.Errorf("Warning for using 'yes' for boolean from compose-go should appear exactly once")
	}
}
