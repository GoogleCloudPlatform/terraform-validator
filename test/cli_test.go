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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

// TestCLI tests the "convert" and "validate" subcommand against a generated .tfplan file.
func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	// Define the reusable constraints to be use for the test cases.
	type constraint struct {
		name            string
		wantViolation   bool
		wantOutputRegex string
	}
	// Currently, we only test one rule. Moving forward, resource specific rules
	// should be added to increase the coverage.
	alwaysViolate := constraint{name: "always_violate", wantViolation: true, wantOutputRegex: "Constraint always_violates_all on resource"}

	// Test cases for each type of resource is defined here.
	cases := []struct {
		name        string
		constraints []constraint
	}{
		{name: "bucket"},
		{name: "bucket_iam"},
		{name: "disk"},
		{name: "firewall"},
		{name: "instance"},
		{name: "sql"},
		{name: "example_bigquery_dataset"},
		{name: "example_compute_disk"},
		{name: "example_compute_firewall"},
		{name: "example_compute_instance"},
		{name: "example_container_cluster"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_member"},
		{name: "example_organization_iam_policy"},
		{name: "example_project"},
		{name: "example_project_in_org"},
		{name: "example_project_in_folder"},
		{name: "example_project_iam"},
		{name: "example_project_iam_binding"},
		{name: "example_project_iam_member"},
		{name: "example_project_iam_policy"},
		{name: "example_sql_database_instance"},
		{name: "example_storage_bucket"},
		{name: "full_compute_firewall"},
		{name: "full_compute_instance"},
		{name: "full_container_cluster"},
		{name: "full_container_node_pool"},
		{name: "full_sql_database_instance"},
		{name: "full_storage_bucket"},
	}
	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		// Add default constraints if not set.
		if len(c.constraints) == 0 {
			c.constraints = []constraint{alwaysViolate}
		}

		// Test both offline and online mode.
		for _, offline := range []bool{true, false} {
			t.Run(fmt.Sprintf("v=%s/tf=%s/offline=%t", version.LeastSupportedVersion(), c.name, offline), func(t *testing.T) {
				t.Parallel()
				// Create a temporary directory for running terraform.
				dir, err := ioutil.TempDir(tmpDir, "terraform")
				if err != nil {
					log.Fatal(err)
				}
				defer os.RemoveAll(dir)

				// Generate the <name>.tf and <name>_assets.json files into the temporary directory.
				generateTestFiles(t, "../testdata/templates", dir, c.name+".tf")
				generateTestFiles(t, "../testdata/templates", dir, c.name+".json")

				terraform(t, dir, c.name)

				t.Run("cmd=convert", func(t *testing.T) {
					testConvertCommand(t, dir, c.name, offline)
				})

				for _, ct := range c.constraints {
					t.Run(fmt.Sprintf("cmd=validate/constraint=%s", ct.name), func(t *testing.T) {
						testValidateCommand(t, ct.wantViolation, ct.wantOutputRegex, dir, c.name, offline, ct.name)
					})
				}
			})
		}
	}
}

func testConvertCommand(t *testing.T, dir, name string, offline bool) {
	var payload []byte
	switch version.LeastSupportedVersion() {
	case version.TF11:
		payload = tfvConvert(t, dir, name+".tfplan", offline)
	default:
		payload = tfvConvert(t, dir, name+".tfplan.json", offline)
	}
	// Verify if the generated assets match the expected.
	var got []google.Asset
	err := json.Unmarshal(payload, &got)
	if err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}
	testfile := filepath.Join(dir, name+".json")
	payload, err = ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatalf("Error reading %v: %v", testfile, err)
	}
	var want []google.Asset
	if err := json.Unmarshal(payload, &want); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	gotJSON := normalizeAssets(t, got, offline)
	wantJSON := normalizeAssets(t, want, offline)
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}

func testValidateCommand(t *testing.T, wantViolation bool, want, dir, name string, offline bool, constraintName string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get current directory: %v", err)
	}
	policyPath := filepath.Join(cwd, samplePolicyPath, constraintName)
	var got []byte
	switch version.LeastSupportedVersion() {
	case version.TF11:
		got = tfvValidate(t, wantViolation, dir, name+".tfplan", policyPath, offline)
	default:
		got = tfvValidate(t, wantViolation, dir, name+".tfplan.json", policyPath, offline)
	}
	wantRe := regexp.MustCompile(want)
	if want != "" && !wantRe.Match(got) {
		t.Fatalf("binary did not return expect output, \ngot=%s \nwant (regex)=%s", string(got), want)
	}
}

func terraform(t *testing.T, dir, name string) {
	switch version.LeastSupportedVersion() {
	case version.TF11:
		terraformInit(t, "terraform", dir)
		terraformPlan(t, "terraform", dir, name+".tfplan")
	default:
		terraformInit(t, "terraform", dir)
		terraformPlan(t, "terraform", dir, name+".tfplan")
		payload := terraformShow(t, "terraform", dir, name+".tfplan")
		saveFile(t, dir, name+".tfplan.json", payload)
	}
}

func terraformInit(t *testing.T, executable, dir string) {
	terraformExec(t, executable, dir, "init", "-input=false")
}

func terraformPlan(t *testing.T, executable, dir, tfplan string) {
	terraformExec(t, executable, dir, "plan", "-input=false", "--out", tfplan)
}

func terraformShow(t *testing.T, executable, dir, tfplan string) []byte {
	return terraformExec(t, executable, dir, "show", "--json", tfplan)
}

func terraformExec(t *testing.T, executable, dir string, args ...string) []byte {
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{"HOME=" + filepath.Join(dir, "fakehome")}
	cmd.Dir = dir
	wantError := false
	payload, _ := run(t, cmd, wantError)
	return payload
}

func saveFile(t *testing.T, dir, filename string, payload []byte) {
	fullpath := filepath.Join(dir, filename)
	f, err := os.Create(fullpath)
	if err != nil {
		t.Fatalf("error while creating file %s, error %v", fullpath, err)
	}
	_, err = f.Write(payload)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", fullpath, err)
	}
}

func tfvConvert(t *testing.T, dir, tfplan string, offline bool) []byte {
	executable := tfvBinary
	wantError := false
	args := []string{"convert", "--project", data.Provider["project"]}
	if offline {
		args = append(args, "--ancestry", data.Ancestry)
	}
	args = append(args, tfplan)
	cmd := exec.Command(executable, args...)
	// Remove environment variables inherited from the test runtime.
	cmd.Env = []string{}
	// Add credentials back.
	if data.Provider["credentials"] != "" {
		cmd.Env = append(cmd.Env, "GOOGLE_APPLICATION_CREDENTIALS="+data.Provider["credentials"])
	}
	cmd.Dir = dir
	payload, _ := run(t, cmd, wantError)
	return payload
}

func tfvValidate(t *testing.T, wantError bool, dir, tfplan, policyPath string, offline bool) []byte {
	executable := tfvBinary
	args := []string{"validate", "--project", data.Provider["project"], "--policy-path", policyPath}
	if offline {
		args = append(args, "--offline", "--ancestry", data.Ancestry)
	}
	args = append(args, tfplan)
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{"GOOGLE_APPLICATION_CREDENTIALS=" + data.Provider["credentials"]}
	cmd.Dir = dir
	payload, _ := run(t, cmd, wantError)
	return payload
}

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, cmd *exec.Cmd, wantError bool) ([]byte, []byte) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	err := cmd.Run()
	if gotError := (err != nil); gotError != wantError {
		t.Fatalf("running %s: \nerror=%v \nstderr=%s \nstdout=%s", cmdToString(cmd), err, stderr.String(), stdout.String())
	}
	// Print env, stdout and stderr if verbose flag is used.
	if len(cmd.Env) != 0 {
		t.Logf("=== Environment Variable of %s ===", cmdToString(cmd))
		t.Log(strings.Join(cmd.Env, "\n"))
	}
	if stdout.String() != "" {
		t.Logf("=== STDOUT of %s ===", cmdToString(cmd))
		t.Log(stdout.String())
	}
	if stderr.String() != "" {
		t.Logf("=== STDERR of %s ===", cmdToString(cmd))
		t.Log(stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

// cmdToString clones the logic of https://golang.org/pkg/os/exec/#Cmd.String.
func cmdToString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := new(strings.Builder)
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
