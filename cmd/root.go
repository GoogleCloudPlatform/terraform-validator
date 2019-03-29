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

func init() {
	validateCmd.PersistentFlags().BoolVar(&flags.verbose, "verbose", false, "Log output to stderr")

	validateCmd.Flags().StringVar(&flags.validate.policyPath, "policy-path", "", "Path to directory containing validation policies")
	validateCmd.MarkFlagRequired("policy-path")
	validateCmd.Flags().StringVar(&flags.validate.project, "project", "", "Provider project")
	validateCmd.MarkFlagRequired("project")

	convertCmd.Flags().StringVar(&flags.convert.project, "project", "", "Provider project")
	convertCmd.MarkFlagRequired("project")

	validateCmd.Flags().BoolVar(&flags.validate.outputJSON, "output-json", false, "Print violations as JSON")

	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listSupportedResourcesCmd)
}

// NOTE: We use a pkg-level var here instead of github.com/spf13/viper
// to establish a pattern of passing down config rather than accessing it
// globally.
var flags struct {
	// Common flags
	verbose bool

	// flags that correspond to subcommands:
	convert struct {
		project string
	}
	validate struct {
		project    string
		policyPath string
		outputJSON bool
	}
	listSupportedResources struct{}
}

// Execute is the entry-point for all commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "terraform-validator",
	Short: "Validate terraform plans using Forseti Config Validator.",
	Long: `Validate terraform plans by converting terraform resources
to their Google CAI (Cloud Asset Inventory) format and passing them through
Forseti Config Validator.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !flags.verbose {
			// Suppress chatty packages.
			log.SetOutput(ioutil.Discard)
		}
	},
}
