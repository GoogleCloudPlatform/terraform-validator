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
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <tfplan>",
	Short: "Validate resources in a Terraform plan by calling Forseti Config Validator.",
	Long: `Validate (terraform-validator validate) converts supported Terraform
resources (see: "terraform-validate list-supported-resources") into their CAI
(Cloud Asset Inventory) format and calls Forseti Config Validator,
returning the violations. If any violations are reported an exit code of 2
is set.

Example:
  terraform-validator validate ./example/terraform.tfplan \
    --project my-project \
    --ancestry organization/my-org/folder/my-folder \
    --policy-path ./path/to/my/gcv/policies
`,
	PreRunE: func(c *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing required argument <tfplan>")
		}
		if flags.validate.offline && flags.validate.ancestry == "" {
			return errors.New("please set ancestry via --ancestry in offline mode")
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		assets, err := tfgcv.ReadPlannedAssets(args[0], flags.validate.project,
			flags.validate.ancestry, flags.validate.offline)
		if err != nil {
			if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
				return errors.New("unable to parse provider project, please use --project flag")
			}
			return errors.Wrap(err, "converting tfplan to CAI assets")
		}

		auditResult, err := tfgcv.ValidateAssets(assets, flags.validate.policyPath)
		if err != nil {
			return errors.Wrap(err, "validating: FCV")
		}

		if len(auditResult.Violations) > 0 {
			if flags.validate.outputJSON {
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

			os.Exit(2)
		}

		if !flags.validate.outputJSON {
			fmt.Println("No violations found.")
		}
		return nil
	},
}
