package cmd

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/stretchr/testify/assert"
)

func MockReadPlannedAssets(ctx context.Context, path, project, ancestry string, offline, convertUnchanged bool) ([]google.Asset, error) {
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
	}, nil
}

func TestConvertRun(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := true
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
	o := convertOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
	}

	err := o.run("/path/to/plan")
	a.Nil(err)

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	a.Len(errorLogs, 0)
	a.Len(outputLogs, 1)

	// On a successful run, we should see a list of google assets in the resource_body field
	fields := outputLogs[0].ContextMap()
	a.Contains(fields, "resource_body")
	a.Len(fields["resource_body"], 1)
	a.IsType([]google.Asset{}, fields["resource_body"])
}

func TestConvertRunLegacy(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := false
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
	o := convertOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
	}

	err := o.run("/path/to/plan")
	a.Nil(err)

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	// On a successful legacy run, we don't output anything via loggers.
	a.Len(errorLogs, 0)
	a.Len(outputLogs, 0)
}
