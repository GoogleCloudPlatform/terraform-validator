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
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// LoggerStdErr used by commands to print errors and warnings
var LoggerStdErr = log.New(os.Stderr, "", log.LstdFlags)

func init() {
	rootCmd.PersistentFlags().BoolVar(&flags.verbose, "verbose", false, "Log output to stderr")

	validateCmd.Flags().StringVar(&flags.validate.policyPath, "policy-path", "", "Path to directory containing validation policies")
	validateCmd.MarkFlagRequired("policy-path")
	validateCmd.Flags().StringVar(&flags.validate.project, "project", "", "Provider project override (override the default project configuration assigned to the google terraform provider when validating resources)")
	validateCmd.Flags().StringVar(&flags.validate.ancestry, "ancestry", "", "Override the ancestry location of the project when validating resources")
	validateCmd.Flags().BoolVar(&flags.validate.offline, "offline", false, "Do not make network requests")
	validateCmd.Flags().BoolVar(&flags.validate.outputJSON, "output-json", false, "Print violations as JSON")

	convertCmd.Flags().StringVar(&flags.convert.project, "project", "", "Provider project override (override the default project configuration assigned to the google terraform provider when converting resources)")
	convertCmd.Flags().StringVar(&flags.convert.ancestry, "ancestry", "", "Override the ancestry location of the project when validating resources")
	convertCmd.Flags().BoolVar(&flags.convert.offline, "offline", false, "Do not make network requests")

	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(listSupportedResourcesCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

// NOTE: We use a pkg-level var here instead of github.com/spf13/viper
// to establish a pattern of passing down config rather than accessing it
// globally.
var flags struct {
	// Common flags
	verbose bool

	// flags that correspond to subcommands:
	convert struct {
		project  string
		ancestry string
		offline    bool
	}
	validate struct {
		project    string
		ancestry   string
		offline    bool
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
	Short: "Validate that a terraform plan conforms to Constraint Framework policies.",
	Long: fmt.Sprintf(`Validate that a terraform plan conforms to a Constraint Framework 
policy library written to expect Google CAI (Cloud Asset Inventory) data.

Supported Terraform versions = 0.12+`),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !flags.verbose {
			// Suppress chatty packages.
			log.SetOutput(ioutil.Discard)
		}
		return nil
	},
}
