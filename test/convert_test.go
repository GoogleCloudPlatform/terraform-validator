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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
)

const samplePolicyPath = "sample_policies"

var conversionTests = []struct {
	name      string
	assetType string
}{
	{"disk", "compute.googleapis.com/Disk"},
	{"project", "cloudresourcemanager.googleapis.com/Project"},
	{"project_billing_info", "cloudbilling.googleapis.com/ProjectBillingInfo"},
	{"firewall", "compute.googleapis.com/Firewall"},
	{"instance", "compute.googleapis.com/Instance"},
	{"bucket", "storage.googleapis.com/Bucket"},
	{"sql", "sqladmin.googleapis.com/Instance"},
}

// TestConvert tests the "convert" subcommand against a generated .tfplan file.
func TestConvert(t *testing.T) {
	for _, tfVersion := range []string{tfplan.TF11} {
		_, cfg := setup(tfVersion, t)

		err, stdOutput, errOutput := runWithCred(cfg.credentials,
			filepath.Join("..", "bin", "terraform-validator"),
			"convert",
			"--tf-version", tfVersion,
			"--project", cfg.project,
			planPath,
		)

		if err != nil {
			t.Fatalf("%v:\n%v", err, string(errOutput))
		}

		var assets []google.Asset
		if err := json.Unmarshal(stdOutput, &assets); err != nil {
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
			t.Run(tt.name+"/TF_"+tfVersion, func(t *testing.T) {
				if len(assetsByType[tt.assetType]) == 0 {
					t.Fatalf("asset type %q not found", tt.assetType)
				}
				if len(jsonFixtures[tt.name]) == 0 {
					t.Fatalf("json fixtures %q not found", tt.name)
				}
				requireEqualJSON(t,
					jsonFixtures[tt.name],
					assetsByType[tt.assetType][0].Resource.Data,
				)
			})
		}

		validationTests := []struct {
			name            string
			wantError       bool
			wantOutputRegex string
		}{
			{
				name:            "always_violate",
				wantError:       true,
				wantOutputRegex: "Constraint always_violates_all on resource",
			},
		}
		for _, tt := range validationTests {
			t.Run(fmt.Sprintf("validate/%s", tt.name), func(t *testing.T) {
				wantRe := regexp.MustCompile(tt.wantOutputRegex)
				err, stdOutput, errOutput := runWithCred(cfg.credentials,
					filepath.Join("..", "bin", "terraform-validator"),
					"validate",
					"--project", cfg.project,
					"--ancestry", "/organization/test",
					"--policy-path", filepath.Join(samplePolicyPath, tt.name),
					planPath,
				)
				if gotError := (err != nil); gotError != tt.wantError {
					t.Fatalf("binary return %v with stderr=%s, got %v, want %v", err, errOutput, gotError, tt.wantError)
				}
				if tt.wantOutputRegex != "" && !wantRe.Match(stdOutput) {
					t.Fatalf("binary did not return expect output, got=%s\nwant (regex)=%s", string(stdOutput), tt.wantOutputRegex)
				}
			})
		}
	}
}
