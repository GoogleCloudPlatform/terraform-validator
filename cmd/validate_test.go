package cmd

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/stretchr/testify/assert"
)

func MockValidateAssetsNoViolations(ctx context.Context, assets []google.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return &validator.AuditResponse{
		Violations: []*validator.Violation{},
	}, nil
}

func MockValidateAssetsWithViolations(ctx context.Context, assets []google.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return &validator.AuditResponse{
		Violations: []*validator.Violation{
			&validator.Violation{
				Constraint: "GCPAlwaysViolatesConstraintV1.always_violates_project_match_target",
				Resource:   "//bigtable.googleapis.com/projects/my-project/instances/tf-instance",
				Message:    "Constraint GCPAlwaysViolatesConstraintV1.always_violates_project_match_target on resource //bigtable.googleapis.com/projects/my-project/instances/tf-instance",
				Severity:   "high",
			},
		},
	}, nil
}

func TestValidateRun(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := true
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
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

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	a.Len(errorLogs, 0)
	a.Len(outputLogs, 1)

	// On a run with violations, we should see a list of violations in the resource_body field
	fields := outputLogs[0].ContextMap()
	a.Contains(fields, "resource_body")
	a.Len(fields["resource_body"], 1)
	a.IsType([]*validator.Violation{}, fields["resource_body"])
}

func TestValidateRunNoViolations(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := true
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
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

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	a.Len(errorLogs, 0)
	a.Len(outputLogs, 1)

	// On a run with no violations, we should see an empty list of violations in the resource_body field
	fields := outputLogs[0].ContextMap()
	a.Contains(fields, "resource_body")
	a.Len(fields["resource_body"], 0)
	a.IsType([]*validator.Violation{}, fields["resource_body"])
}

func TestValidateRunLegacy(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := false
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
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

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	// On a legacy run with validation errors, loggers should not be used.
	a.Len(errorLogs, 0)
	a.Len(outputLogs, 0)
}

func TestValidateRunNoViolationsLegacy(t *testing.T) {
	a := assert.New(t)
	verbose := true
	useStructuredLogging := false
	errorLogger, errorObservedLogs := newTestErrorLogger(verbose, useStructuredLogging)
	outputLogger, outputObservedLogs := newTestOutputLogger()
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

	errorLogs := errorObservedLogs.AllUntimed()
	outputLogs := outputObservedLogs.AllUntimed()

	// On a legacy run with no validation errors, loggers should not be used.
	a.Len(errorLogs, 0)
	a.Len(outputLogs, 0)
}
