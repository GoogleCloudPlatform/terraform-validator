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

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const validateDesc = `
Validate that a terraform plan conforms to a Constraint Framework
policy library written to expect Google CAI (Cloud Asset Inventory) data. 
Unsupported terraform resources (see: "terraform-validate list-supported-resources")
are skipped.

Policy violations will result in an exit code of 2.

Example:
  terraform-validator validate ./example/terraform.tfplan \
    --project my-project \
    --ancestry organization/my-org/folder/my-folder \
    --policy-path ./path/to/my/gcv/policies
`

type validateOptions struct {
	project           string
	ancestry          string
	offline           bool
	policyPath        string
	outputJSON        bool
	dryRun            bool
	rootOptions       *rootOptions
	readPlannedAssets tfgcv.ReadPlannedAssetsFunc
	validateAssets    tfgcv.ValidateAssetsFunc
}

func newValidateCmd(rootOptions *rootOptions) *cobra.Command {
	o := &validateOptions{
		rootOptions:       rootOptions,
		readPlannedAssets: tfgcv.ReadPlannedAssets,
		validateAssets:    tfgcv.ValidateAssets,
	}

	cmd := &cobra.Command{
		Use:   "validate TFPLAN_JSON --policy-path=/path/to/policy/library",
		Short: "Validate that a terraform plan conforms to Constraint Framework policies",
		Long:  validateDesc,
		PreRunE: func(c *cobra.Command, args []string) error {
			return o.validateArgs(args)
		},
		RunE: func(c *cobra.Command, args []string) error {
			if o.dryRun {
				return nil
			}
			return o.run(args[0])
		},
	}

	cmd.Flags().StringVar(&o.policyPath, "policy-path", "", "Path to directory containing validation policies")
	cmd.MarkFlagRequired("policy-path")
	cmd.Flags().StringVar(&o.project, "project", "", "Provider project override (override the default project configuration assigned to the google terraform provider when validating resources)")
	cmd.Flags().StringVar(&o.ancestry, "ancestry", "", "Override the ancestry location of the project when validating resources")
	cmd.Flags().BoolVar(&o.offline, "offline", false, "Do not make network requests")
	cmd.Flags().BoolVar(&o.outputJSON, "output-json", false, "Print violations as JSON")
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Only parse & validate args")
	cmd.Flags().MarkHidden("dry-run")

	return cmd
}

func (o *validateOptions) validateArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("missing required argument TFPLAN_JSON")
	}
	if o.offline && o.ancestry == "" {
		return errors.New("please set ancestry via --ancestry in offline mode")
	}
	return nil
}

func (o *validateOptions) run(plan string) error {
	ctx := context.Background()
	assets, err := o.readPlannedAssets(ctx, plan, o.project, o.ancestry, o.offline, false, o.rootOptions.errorLogger)
	if err != nil {
		if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
			return errors.New("unable to parse provider project, please use --project flag")
		}
		return errors.Wrap(err, "converting tfplan to CAI assets")
	}

	auditResult, err := o.validateAssets(ctx, assets, o.policyPath)
	if err != nil {
		return errors.Wrap(err, "validating: FCV")
	}

	if o.rootOptions.useStructuredLogging {
		msg := "No violations found"
		if len(auditResult.Violations) > 0 {
			msg = "Violations found"
		}
		o.rootOptions.outputLogger.Info(
			msg,
			zap.Any("resource_body", auditResult.Violations),
		)
		if len(auditResult.Violations) > 0 {
			return errViolations
		}
		return nil
	}

	// Legacy behavior
	if len(auditResult.Violations) > 0 {
		if o.outputJSON {
			marshaller := &jsonpb.Marshaler{}
			if err := marshaller.Marshal(os.Stdout, auditResult); err != nil {
				return errors.Wrap(err, "marshalling violations to json")
			}
		} else {
			fmt.Print("Found Violations:\n\n")
			for _, v := range auditResult.Violations {
				fmt.Printf("Constraint %v on resource %v: %v\n\n",
					v.Constraint,
					v.Resource,
					v.Message,
				)
			}
		}

		return errViolations
	}

	if !o.outputJSON {
		fmt.Println("No violations found.")
	}
	return nil
}
