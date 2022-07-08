package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/version"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func resetEnvKeys() []string {
	return []string{
		"GOOGLE_ZONE",
		"GCLOUD_ZONE",
		"CLOUDSDK_COMPUTE_ZONE",
		"GOOGLE_REGION",
		"GCLOUD_REGION",
		"CLOUDSDK_COMPUTE_REGION",
	}
}

func testAssets(path, project, zone, region string, ancestry map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger, userAgent string) []google.Asset {
	return []google.Asset{
		google.Asset{
			Name: "//compute.googleapis.com/projects/my-project/zones/us-central1-a/disks/test-disk",
			Type: "compute.googleapis.com/Disk",
			Resource: &google.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Disk",
				Data: map[string]interface{}{
					"labels": map[string]string{
						"environment": "dev",
					},
					"name":                   "test-disk",
					"physicalBlockSizeBytes": 4096,
					"sourceImage":            "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
					"type":                   "projects/my-project/zones/us-central1-a/diskTypes/pd-ssd",
					"zone":                   "projects/my-project/global/zones/us-central1-a",
					"arguments": map[string]interface{}{
						"path":             path,
						"project":          project,
						"zone":             zone,
						"region":           region,
						"ancestry":         ancestry,
						"offline":          offline,
						"convertUnchanged": false,
						"errorLoggger":     errorLogger,
						"userAgent":        userAgent,
					},
				},
			},
		},
	}
}

func MockReadPlannedAssets(ctx context.Context, path, project, zone, region string, ancestry map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger, userAgent string) ([]google.Asset, error) {
	return testAssets(path, project, zone, region, ancestry, offline, convertUnchanged, errorLogger, userAgent), nil
}

func TestConvertRun(t *testing.T) {
	for _, k := range resetEnvKeys() {
		k := k
		originalValue, isSet := os.LookupEnv(k)
		if isSet {
			defer os.Setenv(k, originalValue)
		} else {
			defer os.Unsetenv(k)
		}
		err := os.Setenv(k, "")
		if err != nil {
			t.Fatalf("error clearing env var %s: %s", k, err)
		}
	}

	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := true
	errorLogger, errorBuf := newTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	ro := &rootOptions{
		verbosity:            verbosity,
		useStructuredLogging: useStructuredLogging,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
	}
	o := convertOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
	}

	path := "/path/to/plan"
	err := o.run(path)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.Bytes()

	a.Equal(errorJSON, "")

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	// On a successful run, we should see a list of google assets in the resource_body field
	a.Contains(output, "resource_body")
	a.Len(output["resource_body"], 1)

	var expectedAssets []interface{}
	expectedAssetJSON, _ := json.Marshal(testAssets(path, "", "", "", map[string]string{}, false, false, errorLogger, fmt.Sprintf("config-validator-tf/%s", version.BuildVersion())))
	json.Unmarshal(expectedAssetJSON, &expectedAssets)
	a.Equal(expectedAssets, output["resource_body"])
}

func TestConvertRunLegacy(t *testing.T) {
	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := false
	errorLogger, errorBuf := newTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	ro := &rootOptions{
		verbosity:            verbosity,
		useStructuredLogging: useStructuredLogging,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
	}
	o := convertOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
	}

	err := o.run("/path/to/plan")
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a successful legacy run, we don't output anything via loggers.
	a.Equal(errorJSON, "")
	a.Equal(outputJSON, "")
}

func TestConvertRunOutputFile(t *testing.T) {
	for _, k := range resetEnvKeys() {
		k := k
		originalValue, isSet := os.LookupEnv(k)
		if isSet {
			defer os.Setenv(k, originalValue)
		} else {
			defer os.Unsetenv(k)
		}
		err := os.Setenv(k, "")
		if err != nil {
			t.Fatalf("error clearing env var %s: %s", k, err)
		}
	}

	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := false
	errorLogger, errorBuf := newTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	ro := &rootOptions{
		verbosity:            verbosity,
		useStructuredLogging: useStructuredLogging,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
	}
	outputPath := path.Join(t.TempDir(), "converted.json")
	o := convertOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
		outputPath:        outputPath,
	}

	path := "/path/to/plan"
	err := o.run(path)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	a.Equal(errorJSON, "")
	a.Equal(outputJSON, "")

	b, err := ioutil.ReadFile(outputPath)
	if err != nil {
		a.Failf("Unable to read file %s: %s", outputPath, err)
	}
	var gotAssets []interface{}
	err = json.Unmarshal(b, &gotAssets)
	if err != nil {
		a.Failf("Failed to unmarshal file %s: %s", outputPath, err)
	}

	var expectedAssets []interface{}
	expectedAssetJSON, _ := json.Marshal(testAssets(path, "", "", "", map[string]string{}, false, false, errorLogger, fmt.Sprintf("config-validator-tf/%s", version.BuildVersion())))
	json.Unmarshal(expectedAssetJSON, &expectedAssets)
	a.Equal(expectedAssets, gotAssets)
}

func TestConvertRun_passesCorrectArguments(t *testing.T) {
	cases := []struct {
		name         string
		project      string
		ancestry     string
		envKey       string
		envValue     string
		wantAncestry string
		wantProject  string
		wantZone     string
		wantRegion   string
	}{
		{
			name:        "project",
			project:     "my-project",
			wantProject: "my-project",
		},
		{
			name:         "project with ancestry",
			project:      "my-project",
			ancestry:     "organizations/1234/folders/5678",
			wantProject:  "my-project",
			wantAncestry: "organizations/1234/folders/5678",
		},
		{
			name:     "GOOGLE_ZONE",
			envKey:   "GOOGLE_ZONE",
			envValue: "whatever",
			wantZone: "whatever",
		},
		{
			name:     "GCLOUD_ZONE",
			envKey:   "GCLOUD_ZONE",
			envValue: "whatever",
			wantZone: "whatever",
		},
		{
			name:     "CLOUDSDK_COMPUTE_ZONE",
			envKey:   "CLOUDSDK_COMPUTE_ZONE",
			envValue: "whatever",
			wantZone: "whatever",
		},
		{
			name:       "GOOGLE_REGION",
			envKey:     "GOOGLE_REGION",
			envValue:   "whatever",
			wantRegion: "whatever",
		},
		{
			name:       "GCLOUD_REGION",
			envKey:     "GCLOUD_REGION",
			envValue:   "whatever",
			wantRegion: "whatever",
		},
		{
			name:       "CLOUDSDK_COMPUTE_REGION",
			envKey:     "CLOUDSDK_COMPUTE_REGION",
			envValue:   "whatever",
			wantRegion: "whatever",
		},
	}

	for _, k := range resetEnvKeys() {
		k := k
		originalValue, isSet := os.LookupEnv(k)
		if isSet {
			defer os.Setenv(k, originalValue)
		} else {
			defer os.Unsetenv(k)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// clear env vars before each test
			for _, k := range resetEnvKeys() {
				err := os.Setenv(k, "")
				if err != nil {
					t.Fatalf("error clearing env var %s: %s", k, err)
				}
			}
			if c.envKey != "" {
				err := os.Setenv(c.envKey, c.envValue)
				if err != nil {
					t.Fatalf("error setting env var %s=%s: %s", c.envKey, c.envValue, err)
				}

			}

			a := assert.New(t)
			verbosity := "debug"
			useStructuredLogging := true
			errorLogger, errorBuf := newTestErrorLogger(verbosity, useStructuredLogging)
			outputLogger, outputBuf := newTestOutputLogger()
			ro := &rootOptions{
				verbosity:            verbosity,
				useStructuredLogging: useStructuredLogging,
				errorLogger:          errorLogger,
				outputLogger:         outputLogger,
			}
			o := convertOptions{
				project:           c.project,
				ancestry:          c.ancestry,
				offline:           false,
				rootOptions:       ro,
				readPlannedAssets: MockReadPlannedAssets,
			}

			path := "/path/to/plan"
			err := o.run(path)
			a.Nil(err)

			errorJSON := errorBuf.String()
			outputJSON := outputBuf.Bytes()

			a.Equal(errorJSON, "")

			var output map[string]interface{}
			json.Unmarshal(outputJSON, &output)

			// On a successful run, we should see a list of google assets in the resource_body field
			a.Contains(output, "resource_body")
			a.Len(output["resource_body"], 1)

			wantAncestryCache := map[string]string{}
			if c.wantProject != "" {
				wantAncestryCache[c.wantProject] = c.wantAncestry
			}
			var expectedAssets []interface{}
			expectedAssetJSON, _ := json.Marshal(testAssets(path, c.wantProject, c.wantZone, c.wantRegion, wantAncestryCache, false, false, errorLogger, fmt.Sprintf("config-validator-tf/%s", version.BuildVersion())))
			json.Unmarshal(expectedAssetJSON, &expectedAssets)
			a.Equal(expectedAssets, output["resource_body"])
		})
	}
}
