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
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
)

// TestConvert tests the "convert" subcommand against a generated .tfplan file.
func TestConvert(t *testing.T) {
	_, cfg := setup(t)

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

	jsonFixtures := make(map[string][]byte)

	matches, _ := filepath.Glob(filepath.Join(jsonGenerateDir, "*.json"))
	for _, fixturePath := range matches {
		fixtureFileName := strings.TrimPrefix(fixturePath, jsonGenerateDir+"/")
		fixtureName := strings.TrimSuffix(fixtureFileName, ".json")

		fixtureData, err := ioutil.ReadFile(fixturePath)
		if err != nil {
			t.Fatalf("Reading: %v", err)
		}

		jsonFixtures[fixtureName] = fixtureData
	}

	t.Run("Disk", func(t *testing.T) {
		requireEqualJSON(t,
			jsonFixtures["disk"],
			assetsByType["compute.googleapis.com/Disk"][0].Resource.Data,
		)
	})

	t.Run("Project", func(t *testing.T) {
		requireEqualJSON(t,
			jsonFixtures["project"],
			assetsByType["cloudresourcemanager.googleapis.com/Project"][0].Resource.Data,
		)
	})

	t.Run("ProjectBillingInfo", func(t *testing.T) {
		requireEqualJSON(t,
			jsonFixtures["project_billing_info"],
			assetsByType["cloudbilling.googleapis.com/ProjectBillingInfo"][0].Resource.Data,
		)
	})

	t.Run("Firewall", func(t *testing.T) {
		requireEqualJSON(t,
			jsonFixtures["firewall"],
			assetsByType["compute.googleapis.com/Firewall"][0].Resource.Data,
		)
	})
}
