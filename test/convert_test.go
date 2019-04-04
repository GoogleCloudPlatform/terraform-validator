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
	cloudbillingv1 "google.golang.org/api/cloudbilling/v1"
	cloudresourcemanagerv1 "google.golang.org/api/cloudresourcemanager/v1"
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
		requireEqualJSONValues(t,
			// Expected:
			data.Disk,
			// Received:
			assetsByType["compute.googleapis.com/Disk"][0].Resource.Data,
			// Type of received data:
			&computev1.Disk{},
		)
	})

	t.Run("Project", func(t *testing.T) {
		requireEqualJSONValues(t,
			data.Project,
			assetsByType["cloudresourcemanager.googleapis.com/Project"][0].Resource.Data,
			&cloudresourcemanagerv1.Project{},
		)
	})

	t.Run("ProjectBillingInfo", func(t *testing.T) {
		requireEqualJSONValues(t,
			data.ProjectBillingInfo,
			assetsByType["cloudbilling.googleapis.com/ProjectBillingInfo"][0].Resource.Data,
			&cloudbillingv1.ProjectBillingInfo{},
		)
	})
}
