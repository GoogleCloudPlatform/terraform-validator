package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/config-validator/pkg/api/validator"
	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
	"github.com/stretchr/testify/assert"
)

func testNoViolations() []*validator.Violation {
	return []*validator.Violation{}
}

func MockValidateAssetsNoViolations(ctx context.Context, assets []google.Asset, policyRootPath string) ([]*validator.Violation, error) {
	return testNoViolations(), nil
}

func testWithViolations() []*validator.Violation {
	return []*validator.Violation{
		&validator.Violation{
			Constraint: "GCPAlwaysViolatesConstraintV1.always_violates_project_match_target",
			Resource:   "//bigtable.googleapis.com/projects/my-project/instances/tf-instance",
			Message:    "Constraint GCPAlwaysViolatesConstraintV1.always_violates_project_match_target on resource //bigtable.googleapis.com/projects/my-project/instances/tf-instance",
			Severity:   "high",
		},
	}
}

func MockValidateAssetsWithViolations(ctx context.Context, assets []google.Asset, policyRootPath string) ([]*validator.Violation, error) {
	return testWithViolations(), nil
}

func TestValidateRun(t *testing.T) {
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
	o := validateOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		policyPath:        "",
		outputJSON:        false,
		dryRun:            false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
		validateAssets:    MockValidateAssetsWithViolations,
	}

	inputPath := createEmptyFile(t, []byte{'0'})
	err := o.run(inputPath)
	a.ErrorIs(err, errViolations)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.Bytes()

	a.Equal("", errorJSON)

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	// On a run with violations, we should see a list of violations in the resource_body field
	a.Contains(output, "resource_body")
	a.Len(output["resource_body"], 1)

	var expectedViolations []interface{}
	expectedViolationsJSON, _ := json.Marshal(testWithViolations())
	json.Unmarshal(expectedViolationsJSON, &expectedViolations)
	a.Equal(expectedViolations, output["resource_body"])
}

func TestValidateRunNoViolations(t *testing.T) {
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
	o := validateOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		policyPath:        "",
		outputJSON:        false,
		dryRun:            false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
		validateAssets:    MockValidateAssetsNoViolations,
	}

	inputPath := createEmptyFile(t, []byte{'0'})
	err := o.run(inputPath)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.Bytes()

	a.Equal("", errorJSON)

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	// On a run with no violations, we should see an empty list of violations in the resource_body field
	a.Contains(output, "resource_body")
	a.Len(output["resource_body"], 0)
}

func TestValidateRunLegacy(t *testing.T) {
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
	o := validateOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		policyPath:        "",
		outputJSON:        false,
		dryRun:            false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
		validateAssets:    MockValidateAssetsWithViolations,
	}

	inputPath := createEmptyFile(t, []byte{'0'})
	err := o.run(inputPath)
	a.ErrorIs(err, errViolations)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a legacy run with validation errors, loggers should not be used.
	a.Equal("", errorJSON)
	a.Equal("", outputJSON)
}

func TestValidateRunNoViolationsLegacy(t *testing.T) {
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
	o := validateOptions{
		project:           "",
		ancestry:          "",
		offline:           false,
		policyPath:        "",
		outputJSON:        false,
		dryRun:            false,
		rootOptions:       ro,
		readPlannedAssets: MockReadPlannedAssets,
		validateAssets:    MockValidateAssetsNoViolations,
	}

	inputPath := createEmptyFile(t, []byte{'0'})
	err := o.run(inputPath)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a legacy run with no validation errors, loggers should not be used.
	a.Equal("", errorJSON)
	a.Equal("", outputJSON)
}

func TestValidateRunWithConvertedAssets(t *testing.T) {
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
	o := validateOptions{
		project:        "",
		ancestry:       "",
		offline:        false,
		policyPath:     "",
		outputJSON:     false,
		dryRun:         false,
		rootOptions:    ro,
		validateAssets: MockValidateAssetsNoViolations,
	}

	inputPath := path.Join(t.TempDir(), "testfile.json")
	data, err := json.Marshal(testAssets(inputPath, "", "", "", map[string]string{}, false, false, errorLogger, fmt.Sprintf("config-validator-tf/%s", version.BuildVersion())))
	if err != nil {
		t.Fatalf("Failed to marshal assets: %s", err)
	}

	if err := ioutil.WriteFile(inputPath, data, os.ModePerm); err != nil {
		t.Fatalf("Failed to write file %s: %s", inputPath, err)
	}
	err = o.run(inputPath)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a legacy run with no validation errors, loggers should not be used.
	a.Equal("", errorJSON)
	a.Equal("", outputJSON)
}

func TestValidateRun_passesCorrectArguments(t *testing.T) {
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
			errorLogger, _ := newTestErrorLogger(verbosity, useStructuredLogging)
			outputLogger, _ := newTestOutputLogger()
			inputPath := createEmptyFile(t, []byte{'0'})
			wantAncestryCache := map[string]string{}
			if c.wantProject != "" {
				wantAncestryCache[c.wantProject] = c.wantAncestry
			}
			expectedAssets := testAssets(inputPath, c.wantProject, c.wantZone, c.wantRegion, wantAncestryCache, false, false, errorLogger, fmt.Sprintf("config-validator-tf/%s", version.BuildVersion()))

			ro := &rootOptions{
				verbosity:            verbosity,
				useStructuredLogging: useStructuredLogging,
				errorLogger:          errorLogger,
				outputLogger:         outputLogger,
			}
			o := validateOptions{
				project:           c.project,
				ancestry:          c.ancestry,
				offline:           false,
				policyPath:        "",
				outputJSON:        false,
				dryRun:            false,
				rootOptions:       ro,
				readPlannedAssets: MockReadPlannedAssets,
				validateAssets: func(ctx context.Context, assets []google.Asset, policyRootPath string) ([]*validator.Violation, error) {
					a.Equal(expectedAssets, assets)
					return MockValidateAssetsNoViolations(ctx, assets, policyRootPath)
				},
			}

			o.run(inputPath)
		})
	}
}

func createEmptyFile(t *testing.T, data []byte) string {
	assetPath := path.Join(t.TempDir(), "testfile.json")
	if err := ioutil.WriteFile(assetPath, data, os.ModePerm); err != nil {
		t.Fatalf("Failed to write file %s: %s", assetPath, err)
	}
	return assetPath
}
