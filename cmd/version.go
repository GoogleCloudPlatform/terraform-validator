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

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"

	"github.com/spf13/cobra"
)

type versionOptions struct{}

func newVersionCmd() *cobra.Command {
	o := versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display Terraform Validator version.",
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}

	return cmd
}

func (o *versionOptions) run() error {
	fmt.Printf("Build version: %s\n", tfgcv.BuildVersion())
	return nil
}
