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
	"os"

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert <tfplan>",
	Short: "Convert resources in a Terraform plan to their Google CAI representation.",
	Long: `Convert (terraform-validator convert) will convert a Terraform plan file
into CAI (Cloud Asset Inventory) resources and output them as a JSON array.

Note:
  Only supported resources will be converted. Non supported resources are
  omitted from results.
  Run "terraform-validator list-supported-resources" to see all supported
  resources.

Example:
  terraform-validator convert ./example/terraform.tfplan --project my-project \
    --ancestry organization/my-org/folder/my-folder
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
		ctx := context.Background()
		assets, err := tfgcv.ReadPlannedAssets(ctx, args[0], flags.convert.project, flags.convert.ancestry, flags.convert.offline)
		if err != nil {
			if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
				return errors.New("unable to parse provider project, please use --project flag")
			}
			return errors.Wrap(err, "converting tfplan to CAI assets")
		}

		if err := json.NewEncoder(os.Stdout).Encode(assets); err != nil {
			return errors.Wrap(err, "encoding json")
		}

		return nil
	},
}
