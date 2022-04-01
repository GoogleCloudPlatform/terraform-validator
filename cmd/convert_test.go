package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func testAssets() []google.Asset {
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
				},
			},
		},
	}
}

func MockReadPlannedAssets(ctx context.Context, path, project string, ancestryCache map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger) ([]google.Asset, error) {
	return testAssets(), nil
}

func TestConvertRun(t *testing.T) {
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

	err := o.run("/path/to/plan")
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
	expectedAssetJSON, _ := json.Marshal(testAssets())
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

	err := o.run("/path/to/plan")
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
	expectedAssetJSON, _ := json.Marshal(testAssets())
	json.Unmarshal(expectedAssetJSON, &expectedAssets)
	a.Equal(expectedAssets, gotAssets)
}
