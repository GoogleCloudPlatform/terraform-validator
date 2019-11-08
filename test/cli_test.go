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
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

const (
	samplePolicyPath = "../testdata/sample_policies"
	defaultAncestry  = "organization/12345/folder/67890"
	defaultProject   = "foobar"
)

var (
	data      *testData
	tfvBinary string
)

// testData represents the full dataset that is used for templating terraform
// configs. It contains Google API resources that are expected to be returned
// after converting the terraform plan.
type testData struct {
	// is not nil - Terraform 12 version used
	TFVersion string
	// provider "google"
	Provider map[string]string
	Project  map[string]string
	Ancestry string
}

// init initializes the variables used for testing. As tests rely on
// environment variables, the parsing of those are only done once.
func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current directory: %v", err)
	}
	tfvBinary = filepath.Join(cwd, "..", "bin", "terraform-validator")
	project, ok := os.LookupEnv("TEST_PROJECT")
	if !ok {
		log.Printf("Missing required env var TEST_PROJECT. Default (%s) will be used.", defaultProject)
		project = defaultProject
	}
	credentials, ok := os.LookupEnv("TEST_CREDENTIALS")
	if ok {
		// Make credentials path relative to repo root rather than
		// test/ dir if it is a relative path.
		if !filepath.IsAbs(credentials) {
			credentials = filepath.Join(cwd, "..", credentials)
		}
	} else {
		log.Printf("missing env var TEST_CREDENTIALS, will try to use Application Default Credentials")
	}
	ancestry, ok := os.LookupEnv("TEST_ANCESTRY")
	if !ok {
		log.Printf("Missing required env var TEST_ANCESTRY. Default (%s) will be used.", defaultAncestry)
		ancestry = defaultAncestry
	}
	providerVersion := "1.20"
	if version.TF12 == version.LeastSupportedVersion() {
		providerVersion = "2.12.0"
	}
	data = &testData{
		TFVersion: version.LeastSupportedVersion(),
		Provider: map[string]string{
			"version":     providerVersion,
			"project":     project,
			"credentials": credentials,
		},
		Project: map[string]string{
			"Name":               "My Project Name",
			"ProjectId":          "my-project-id",
			"BillingAccountName": "012345-567890-ABCDEF",
		},
		Ancestry: ancestry,
	}
}

// TestCLI tests the "convert" and "validate" subcommand against a generated .tfplan file.
func TestCLI(t *testing.T) {
	// Define the reusable rules to be use for the test cases.
	type rule struct {
		name            string
		wantError       bool
		wantOutputRegex string
	}
	// Currently, we only test one rule. Moving forward, resource specific rules
	// should be added to increase the coverage.
	alwaysViolate := rule{name: "always_violate", wantError: true, wantOutputRegex: "Constraint always_violates_all on resource"}

	// Test cases for each type of resource is defined here.
	cases := []struct {
		name    string
		offline bool
		rules   []rule
	}{
		{name: "bucket", offline: true, rules: []rule{alwaysViolate}},
		{name: "bucket", offline: false, rules: []rule{alwaysViolate}},
		{name: "disk", offline: true, rules: []rule{alwaysViolate}},
		{name: "disk", offline: false, rules: []rule{alwaysViolate}},
		{name: "firewall", offline: true, rules: []rule{alwaysViolate}},
		{name: "firewall", offline: false, rules: []rule{alwaysViolate}},
		{name: "instance", offline: true, rules: []rule{alwaysViolate}},
		{name: "instance", offline: false, rules: []rule{alwaysViolate}},
		{name: "project", offline: true, rules: []rule{alwaysViolate}},
		{name: "project", offline: false, rules: []rule{alwaysViolate}},
		{name: "sql", offline: true, rules: []rule{alwaysViolate}},
		{name: "sql", offline: false, rules: []rule{alwaysViolate}},
	}
	for _, c := range cases {
		// As tests are run in parallel, the test case need to be cloned to local
		// scope.
		c := c
		t.Run(fmt.Sprintf("v=%s/tf=%s/offline=%t", version.LeastSupportedVersion(), c.name, c.offline), func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for running terraform.
			dir, err := ioutil.TempDir(os.TempDir(), "terraform")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			// Generate the <name>.tf and <name>_assets.json files into the temporary directory.
			generateTestFiles(t, "../testdata/templates", dir, c.name+".tf")
			generateTestFiles(t, "../testdata/templates", dir, c.name+".json")

			terraform(t, dir, c.name)

			t.Run("cmd=convert", func(t *testing.T) {
				testConvertCommand(t, dir, c.name, c.offline)
			})

			for _, r := range c.rules {
				t.Run(fmt.Sprintf("cmd=validate/rule=%s", r.name), func(t *testing.T) {
					testValidateCommand(t, r.wantError, r.wantOutputRegex, dir, c.name, c.offline, r.name)
				})
			}
		})
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

	if !offline {
		// remove the ancestry as the value of that is dependent on project,
		// and is not important for the test.
		for i := range got {
			got[i].Ancestry = ""
		}
		for i := range want {
			want[i].Ancestry = ""
		}
	}
	// compare assets
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}

func testValidateCommand(t *testing.T, wantError bool, want, dir, name string, offline bool, ruleName string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get current directory: %v", err)
	}
	policyPath := filepath.Join(cwd, samplePolicyPath, ruleName)
	var got []byte
	switch version.LeastSupportedVersion() {
	case version.TF11:
		got = tfvValidate(t, wantError, dir, name+".tfplan", policyPath, offline)
	default:
		got = tfvValidate(t, wantError, dir, name+".tfplan.json", policyPath, offline)
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
	wantError := false
	cmd := exec.Command(executable, "init", "-input=false", dir)
	cmd.Env = []string{"HOME=" + filepath.Join(dir, "fakehome")}
	cmd.Dir = dir
	run(t, cmd, wantError)
}

func terraformPlan(t *testing.T, executable, dir, tfplan string) {
	wantError := false
	cmd := exec.Command(executable, "plan", "-input=false", "--out", tfplan, dir)
	cmd.Env = []string{"HOME=" + filepath.Join(dir, "fakehome")}
	cmd.Dir = dir
	run(t, cmd, wantError)
}

func terraformShow(t *testing.T, executable, dir, tfplan string) []byte {
	wantError := false
	cmd := exec.Command(executable, "show", "--json", tfplan)
	cmd.Env = []string{"HOME=" + filepath.Join(dir, "fakehome")}
	cmd.Dir = dir
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
	cmd.Env = []string{"GOOGLE_APPLICATION_CREDENTIALS=" + data.Provider["credentials"]}
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
		t.Fatalf("running %s: \nerror=%v \nstderr=%s \nstdout=%s", cmd, err, stderr.String(), stdout.String())
	}
	// Print env, stdout and stderr if verbose flag is used.
	if len(cmd.Env) != 0 {
		t.Logf("=== Environment Variable of %s ===", cmd)
		t.Log(strings.Join(cmd.Env, "\n"))
	}
	if stdout.String() != "" {
		t.Logf("=== STDOUT of %s ===", cmd)
		t.Log(stdout.String())
	}
	if stderr.String() != "" {
		t.Logf("=== STDERR of %s ===", cmd)
		t.Log(stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

func generateTestFiles(t *testing.T, sourceDir string, targetDir string, selector string) {
	funcMap := template.FuncMap{
		"pastLastSlash": func(s string) string {
			split := strings.Split(s, "/")
			return split[len(split)-1]
		},
	}
	tmpls, err := template.New("").Funcs(funcMap).
		ParseGlob(filepath.Join(sourceDir, selector))
	if err != nil {
		t.Fatalf("generateTestFiles: %v", err)
	}
	for _, tmpl := range tmpls.Templates() {
		if tmpl.Name() == "" {
			continue // Skip base template.
		}
		path := filepath.Join(targetDir, tmpl.Name())
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("creating terraform file %v: %v", path, err)
		}
		if err := tmpl.Execute(f, data); err != nil {
			t.Fatalf("templating terraform file %v: %v", path, err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("closing file %v: %v", path, err)
		}
	}
}
