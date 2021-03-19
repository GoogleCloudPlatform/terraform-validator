package test

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
)

const (
	samplePolicyPath       = "../testdata/sample_policies"
	defaultAncestry        = "organization/12345/folder/67890"
	defaultOrganization    = "12345"
	defaultFolder          = "67890"
	defaultProject         = "foobar"
	defaultProviderVersion = "3.51.0"
)

var (
	data      *testData
	tfvBinary string
	tmpDir    = os.TempDir()
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
	Time     map[string]int64
	OrgID    string
	FolderID string
	Ancestry string
}

// init initializes the variables used for testing. As tests rely on
// environment variables, the parsing of those are only done once.
func init() {
	// don't raise errors in glog

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
	org, ok := os.LookupEnv("TEST_ORG_ID")
	if !ok {
		log.Printf("Missing required env var TEST_ORG_ID. Default (%s) will be used.", defaultOrganization)
		org = defaultOrganization
	}
	folder, ok := os.LookupEnv("TEST_FOLDER_ID")
	if !ok {
		log.Printf("Missing required env var TEST_FOLDER_ID. Default (%s) will be used.", defaultFolder)
		folder = defaultFolder
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
	providerVersion := defaultProviderVersion
	//As time is not information in terraform resource data, time is rounded for testing purposes
	currentTime := time.Now()
	currentTime = currentTime.Round(time.Minute)
	data = &testData{
		TFVersion: "0.12",
		Provider: map[string]string{
			"version":     providerVersion,
			"project":     project,
			"credentials": credentials,
		},
		Time: map[string]int64{
			"timeSeconds": currentTime.Unix(),
			"timeNanos":   currentTime.UnixNano(),
		},
		Project: map[string]string{
			"Name":               "My Project Name",
			"ProjectId":          "my-project-id",
			"BillingAccountName": "012345-567890-ABCDEF",
			"Number":             "1234567890",
		},
		OrgID:    org,
		FolderID: folder,
		Ancestry: ancestry,
	}
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
		log.Printf("Successfully created file %v", path)
	}
}

func normalizeAssets(t *testing.T, got []google.Asset, offline bool) []byte {
	t.Helper()
	if !offline {
		// remove the ancestry as the value of that is dependent on project,
		// and is not important for the test.
		for i := range got {
			got[i].Ancestry = ""
		}
	}
	// Replace placeholder in names.
	re := regexp.MustCompile(`/placeholder-.*`)
	for i := range got {
		got[i].Name = re.ReplaceAllString(got[i].Name, "/placeholder-foobar")
	}
	// compare assets
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	return gotJSON
}
