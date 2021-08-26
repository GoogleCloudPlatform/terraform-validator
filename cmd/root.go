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
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// LoggerStdErr used by commands to print errors and warnings
var LoggerStdErr = log.New(os.Stderr, "", log.LstdFlags)

const rootCmdDesc = `
Validate that a terraform plan conforms to a Constraint Framework 
policy library written to expect Google CAI (Cloud Asset Inventory) data.

Supported Terraform versions = 0.12+`

type rootOptions struct {
	verbose bool
}

func newRootCmd() *cobra.Command {
	o := &rootOptions{}

	cmd := &cobra.Command{
		Use:   "terraform-validator",
		Short: "Validate that a terraform plan conforms to Constraint Framework policies",
		Long: rootCmdDesc,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !o.verbose {
				// Suppress chatty packages.
				log.SetOutput(ioutil.Discard)
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&o.verbose, "verbose", false, "Log output to stderr")

	cmd.AddCommand(newConvertCmd())
	cmd.AddCommand(newListSupportedResourcesCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

// Execute is the entry-point for all commands.
// This lets us keep all new command functions private.
func Execute() {
	rootCmd := newRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
