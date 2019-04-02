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
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	templateDir = "tf_templates"
	generateDir = "tf_generated"
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

	generateTFConfigs(t, data)

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

func generateTFConfigs(t *testing.T, data data) {
	tmpls := template.Must(
		template.New("").
			Funcs(templateFuncs()).
			ParseGlob(filepath.Join(templateDir, "*.tf")))

	for _, tmpl := range tmpls.Templates() {
		if tmpl.Name() == "" {
			continue // Skip base template.
		}
		path := filepath.Join(generateDir, tmpl.Name())

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

// jsonify converts a value into another via JSON as an intermediary
// serialization.
func jsonify(t *testing.T, from, to interface{}) {
	btys, err := json.Marshal(from)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	if err := json.Unmarshal(btys, to); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}
}
