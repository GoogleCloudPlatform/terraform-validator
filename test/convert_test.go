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

var conversionTests = []struct {
	name      string
	assetType string
}{
	{"disk", "compute.googleapis.com/Disk"},
	{"project", "cloudresourcemanager.googleapis.com/Project"},
	{"project_billing_info", "cloudbilling.googleapis.com/ProjectBillingInfo"},
	{"firewall", "compute.googleapis.com/Firewall"},
}

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
			t.Fatalf("Error reading %v: %v", fixturePath, err)
		}

		jsonFixtures[fixtureName] = fixtureData
	}

	for _, tt := range conversionTests {
		// actual := assetsByType[tt.assetType][0].Resource.Data
		t.Run(tt.name, func(t *testing.T) {
			requireEqualJSON(t,
				jsonFixtures[tt.name],
				assetsByType[tt.assetType][0].Resource.Data,
			)
		})
	}
}
