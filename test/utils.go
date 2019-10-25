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
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

const (
	templateDir     = "tf_templates"
	generateDir     = "tf_generated"
	jsonTemplateDir = "json_templates"
	jsonGenerateDir = "json_generated"
)

func getGenerateDir(tfVersion string) string {
	return filepath.Join(generateDir, tfVersion)
}

func getJSONGenerateDir(tfVersion string) string {
	return filepath.Join(jsonGenerateDir, tfVersion)
}

var planPath string

// setup an end-to-end test.
// Pull in env vars (config).
// Build Google API resource values (data).
// Generate terraform files from templates using above data.
// Run terraform fmt/init/plan to create a .tfplan file.
// Return data and config.
func setup(tfVersion string, t *testing.T) (data, config) {
	planPath = filepath.Join(getGenerateDir(tfVersion), "test.tfplan")
	cfg := configure(t)

	var executable string
	switch tfVersion {
	case version.TF11:
		executable = "terraform11"
	case version.TF12:
		executable = "terraform"
	default:
		t.Fatalf("unknown tfVersion %s", tfVersion)
	}

	data := newData(tfVersion, cfg.project, cfg.credentials)

	generateConfigs(t, data, templateDir, getGenerateDir(tfVersion), "*.tf")
	generateConfigs(t, data, jsonTemplateDir, getJSONGenerateDir(tfVersion), "*.json")

	run(t, "rm", "-rf", ".terraform")
	run(t, executable, "fmt", getGenerateDir(tfVersion))
	run(t, executable, "init", getGenerateDir(tfVersion))
	run(t, executable, "plan", "--out", planPath, getGenerateDir(tfVersion))
	if tfVersion == version.TF12 {
		jsonPlanPath := filepath.Join(getGenerateDir(tfVersion), "test.tfplan.json")
		jsonOut, _ := run(t, executable, "show", "-json", planPath)
		f, err := os.Create(jsonPlanPath)
		if err != nil {
			t.Fatalf("error while creating file %s, error %v", jsonPlanPath, err)
		}
		_, err = f.Write(jsonOut)
		if err != nil {
			t.Fatalf("error while writing to file %s, error %v", jsonPlanPath, err)
		}
		// override plan path, to use it in a test
		planPath = jsonPlanPath
	}

	return data, cfg
}

type config struct {
	project     string
	credentials string
}

func configure(t *testing.T) config {
	var cfg config
	var ok bool

	cfg.project, ok = os.LookupEnv("TEST_PROJECT")
	if !ok {
		t.Fatal("missing required env var TEST_PROJECT")
	}

	cfg.credentials, ok = os.LookupEnv("TEST_CREDENTIALS")
	if ok {
		// Make credentials path relative to repo root rather than
		// test/ dir if it is a relative path.
		if !filepath.IsAbs(cfg.credentials) {
			cfg.credentials = filepath.Join("..", cfg.credentials)
		}
	} else {
		t.Log("missing env var TEST_CREDENTIALS, will try to use Application Default Credentials")
	}
	return cfg
}

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, name string, args ...string) ([]byte, []byte) {
	c := exec.Command(name, args...)
	var stderr, stdout bytes.Buffer
	c.Stderr, c.Stdout = &stderr, &stdout
	if err := c.Run(); err != nil {
		t.Fatalf("%s %s: %v, \n %s", name, strings.Join(args, " "), err, stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

func runWithCred(credFile string, name string, args ...string) (error, []byte, []byte) {
	cmd := exec.Command(name, args...)
	cmd.Env = []string{"GOOGLE_APPLICATION_CREDENTIALS=" + credFile}
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	return cmd.Run(), stdout.Bytes(), stderr.Bytes()
}

func generateConfigs(t *testing.T, data data, sourceDir string, targetDir string, selector string) {
	tmpls := template.Must(
		template.New("").
			Funcs(templateFuncs()).
			ParseGlob(filepath.Join(sourceDir, selector)))

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

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"pastLastSlash": func(s string) string {
			split := strings.Split(s, "/")
			return split[len(split)-1]
		},
	}
}

func requireEqualJSON(t *testing.T, expected []byte, provided interface{}) {
	providedJSON, err := json.Marshal(provided)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}

	require.JSONEq(t, string(expected), string(providedJSON))
}
