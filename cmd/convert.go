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
	"encoding/json"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const convertDesc = `
This command will convert a Terraform plan file into CAI (Cloud Asset Inventory)
resources and output them as a JSON array.

Note:
  Only supported resources will be converted. Non supported resources are
  omitted from results.
  Run "terraform-validator list-supported-resources" to see all supported
  resources.

Example:
  terraform-validator convert ./example/terraform.tfplan --project my-project \
    --ancestry organization/my-org/folder/my-folder
`

type convertOptions struct {
	project           string
	ancestry          string
	offline           bool
	rootOptions       *rootOptions
	readPlannedAssets tfgcv.ReadPlannedAssetsFunc
	outputPath        string
	dryRun            bool
}

func newConvertCmd(rootOptions *rootOptions) *cobra.Command {
	o := &convertOptions{
		rootOptions:       rootOptions,
		readPlannedAssets: tfgcv.ReadPlannedAssets,
	}

	cmd := &cobra.Command{
		Use:   "convert TFPLAN_JSON",
		Short: "convert a Terraform plan to Google CAI assets",
		Long:  convertDesc,
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

	cmd.Flags().StringVar(&o.project, "project", "", "Provider project override (override the default project configuration assigned to the google terraform provider when converting resources)")
	cmd.Flags().StringVar(&o.ancestry, "ancestry", "", "Override the ancestry location of the project when validating resources")
	cmd.Flags().BoolVar(&o.offline, "offline", false, "Do not make network requests")
	cmd.Flags().StringVar(&o.outputPath, "output-path", "", "If specified, write the convert result into the specified output file")
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Only parse & validate args")
	cmd.Flags().MarkHidden("dry-run")

	return cmd
}

func (o *convertOptions) validateArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("missing required argument TFPLAN_JSON")
	}
	if o.offline && o.ancestry == "" {
		return errors.New("please set ancestry via --ancestry in offline mode")
	}
	return nil
}

func (o *convertOptions) run(plan string) error {
	ctx := context.Background()
	ancestryCache := map[string]string{
		o.project: o.ancestry,
	}
	userAgent := fmt.Sprintf("config-validator-tf/%s", version.BuildVersion())
	assets, err := o.readPlannedAssets(ctx, plan, o.project, ancestryCache, o.offline, false, o.rootOptions.errorLogger, userAgent)
	if err != nil {
		if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
			return errors.New("unable to parse provider project, please use --project flag")
		}
		return errors.Wrap(err, "converting tfplan to CAI assets")
	}

	if len(o.outputPath) > 0 {
		f, err := os.OpenFile(o.outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(assets); err != nil {
			return errors.Wrap(err, "encoding json")
		}
		return nil
	}

	if o.rootOptions.useStructuredLogging {
		o.rootOptions.outputLogger.Info(
			"converted resources",
			zap.Any("resource_body", assets),
		)
		return nil
	}

	// Legacy behavior
	if err := json.NewEncoder(os.Stdout).Encode(assets); err != nil {
		return errors.Wrap(err, "encoding json")
	}
	return nil
}
