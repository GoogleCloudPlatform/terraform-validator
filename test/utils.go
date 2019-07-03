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
)

const (
	templateDir     = "tf_templates"
	generateDir     = "tf_generated"
	jsonTemplateDir = "json_templates"
	jsonGenerateDir = "json_generated"
)

var planPath = filepath.Join(generateDir, "test.tfplan")

// setup an end-to-end test.
// Pull in env vars (config).
// Build Google API resource values (data).
// Generate terraform files from templates using above data.
// Run terraform fmt/init/plan to create a .tfplan file.
// Return data and config.
func setup(t *testing.T) (data, config) {
	cfg := configure(t)

	data := newData(cfg.project, cfg.credentials)

	generateConfigs(t, data, templateDir, generateDir, "*.tf")
	generateConfigs(t, data, jsonTemplateDir, jsonGenerateDir, "*.json")

	run(t, "terraform", "fmt", generateDir)
	run(t, "terraform", "init", generateDir)
	run(t, "terraform", "plan",
		"--out", planPath,
		generateDir,
	)

	return data, cfg
}

type config struct {
	project     string
	credentials string
	policy      string
}

func configure(t *testing.T) config {
	var cfg config
	var ok bool

	cfg.project, ok = os.LookupEnv("TEST_PROJECT")
	if !ok {
		t.Fatal("missing required env var TEST_PROJECT")
	}

	cfg.credentials, ok = os.LookupEnv("TEST_CREDENTIALS")
	if !ok {
		t.Fatal("missing required env var TEST_CREDENTIALS")
	}

	cfg.policy, ok = os.LookupEnv("TEST_POLICY")
	if !ok {
		t.Fatal("missing required env var TEST_POLICY to specify the path to the policy library")
	}

	// Make credentials path relative to repo root rather than
	// test/ dir if it is a relative path.
	if !filepath.IsAbs(cfg.credentials) {
		cfg.credentials = filepath.Join("..", cfg.credentials)
	}

	return cfg
}

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, name string, args ...string) {
	c := exec.Command(name, args...)
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		t.Fatalf("%s %s: %v", name, strings.Join(args, " "), err)
	}
}

func runWithCred(t *testing.T, credFile string, name string, args ...string) (error, []byte, []byte) {
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
