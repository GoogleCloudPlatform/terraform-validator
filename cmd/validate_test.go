package cmd

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/stretchr/testify/assert"
)

func testAuditResponseNoViolations() *validator.AuditResponse {
	return &validator.AuditResponse{
		Violations: []*validator.Violation{},
	}
}

func MockValidateAssetsNoViolations(ctx context.Context, assets []google.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return testAuditResponseNoViolations(), nil
}

func testAuditResponseWithViolations() *validator.AuditResponse {
	return &validator.AuditResponse{
		Violations: []*validator.Violation{
			&validator.Violation{
				Constraint: "GCPAlwaysViolatesConstraintV1.always_violates_project_match_target",
				Resource:   "//bigtable.googleapis.com/projects/my-project/instances/tf-instance",
				Message:    "Constraint GCPAlwaysViolatesConstraintV1.always_violates_project_match_target on resource //bigtable.googleapis.com/projects/my-project/instances/tf-instance",
				Severity:   "high",
			},
		},
	}
}

func MockValidateAssetsWithViolations(ctx context.Context, assets []google.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return testAuditResponseWithViolations(), nil
}

func TestValidateRun(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := true
	errorLogger, errorBuf := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	o := validateOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		policyPath:           "",
		outputJSON:           false,
		dryRun:               false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
		validateAssets:       MockValidateAssetsWithViolations,
	}

	err := o.run("/path/to/plan")
	a.ErrorIs(err, errViolations)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.Bytes()

	a.Equal("", errorJSON)

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	// On a run with violations, we should see a list of violations in the resource_body field
	a.Contains(output, "resource_body")
	a.Len(output["resource_body"], 1)

	var expectedAuditResponse map[string]interface{}
	expectedAuditResponseJSON, _ := json.Marshal(testAuditResponseWithViolations())
	json.Unmarshal(expectedAuditResponseJSON, &expectedAuditResponse)
	a.Equal(expectedAuditResponse["violations"], output["resource_body"])
}

func TestValidateRunNoViolations(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := true
	errorLogger, errorBuf := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	o := validateOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		policyPath:           "",
		outputJSON:           false,
		dryRun:               false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
		validateAssets:       MockValidateAssetsNoViolations,
	}

	err := o.run("/path/to/plan")
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
	verbose := true
	useStructuredLogging := false
	errorLogger, errorBuf := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	o := validateOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		policyPath:           "",
		outputJSON:           false,
		dryRun:               false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
		validateAssets:       MockValidateAssetsWithViolations,
	}

	err := o.run("/path/to/plan")
	a.ErrorIs(err, errViolations)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a legacy run with validation errors, loggers should not be used.
	a.Equal("", errorJSON)
	a.Equal("", outputJSON)
}

func TestValidateRunNoViolationsLegacy(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := false
	errorLogger, errorBuf := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputBuf := newTestOutputLogger()
	o := validateOptions{
		project:              "",
		ancestry:             "",
		offline:              false,
		policyPath:           "",
		outputJSON:           false,
		dryRun:               false,
		errorLogger:          errorLogger,
		outputLogger:         outputLogger,
		useStructuredLogging: useStructuredLogging,
		readPlannedAssets:    MockReadPlannedAssets,
		validateAssets:       MockValidateAssetsNoViolations,
	}

	err := o.run("/path/to/plan")
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a legacy run with no validation errors, loggers should not be used.
	a.Equal("", errorJSON)
	a.Equal("", outputJSON)
}
