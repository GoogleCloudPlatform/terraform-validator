package test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/stretchr/testify/require"
	computev1 "google.golang.org/api/compute/v1"
)

func TestConvert(t *testing.T) {
	data, cfg := setup(t)

	cmd := exec.Command(filepath.Join("..", "bin", "terraform-validator"),
		"convert",
		"--project", cfg.project,
		planPath,
	)
	cmd.Env = []string{"GOOGLE_APPLICATION_CREDENTIALS=" + cfg.credentials}
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout

	if err := cmd.Run(); err != nil {
		t.Fatalf("%v:\n%v", err, stderr.String())
	}

	var assets []google.Asset
	if err := json.Unmarshal(stdout.Bytes(), &assets); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	t.Run("Disk", func(t *testing.T) {
		var converted computev1.Disk
		jsonify(t, assets[0].Resource.Data, &converted)
		require.Equal(t, data.Disk, converted)
	})
}
