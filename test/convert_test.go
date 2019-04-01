// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// TestConvert tests the "convert" subcommand against a generated .tfplan file.
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

	assetsByType := make(map[string][]google.Asset)
	for _, a := range assets {
		assetsByType[a.Type] = append(assetsByType[a.Type], a)
	}

	t.Run("Disk", func(t *testing.T) {
		var converted computev1.Disk
		jsonify(t, assetsByType["compute.googleapis.com/Disk"][0].Resource.Data, &converted)
		require.Equal(t, data.Disk, converted)
	})
}
