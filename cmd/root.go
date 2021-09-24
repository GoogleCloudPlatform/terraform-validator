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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const rootCmdDesc = `
Validate that a terraform plan conforms to a Constraint Framework 
policy library written to expect Google CAI (Cloud Asset Inventory) data.

Supported Terraform versions = 0.12+`

type rootOptions struct {
	verbose     bool
	errorLogger *zap.Logger
}

func newRootCmd() (*cobra.Command, *zap.Logger, error) {
	o := &rootOptions{}

	cmd := &cobra.Command{
		Use:           "terraform-validator",
		Short:         "Validate that a terraform plan conforms to Constraint Framework policies",
		Long:          rootCmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolVar(&o.verbose, "verbose", false, "Log additional output")

	useStructuredLogging := os.Getenv("USE_STRUCTURED_LOGGING") == "true"
	errorLogger := newErrorLogger(o.verbose, useStructuredLogging, zapcore.Lock(os.Stderr))
	defer errorLogger.Sync()
	zap.RedirectStdLog(errorLogger)
	o.errorLogger = errorLogger

	outputLogger := newOutputLogger(zapcore.Lock(os.Stdout))
	defer outputLogger.Sync()

	cmd.AddCommand(newConvertCmd(errorLogger, outputLogger, useStructuredLogging))
	cmd.AddCommand(newListSupportedResourcesCmd())
	cmd.AddCommand(newValidateCmd(errorLogger, outputLogger, useStructuredLogging))
	cmd.AddCommand(newVersionCmd())

	return cmd, errorLogger, nil
}

// Execute is the entry-point for all commands.
// This lets us keep all new command functions private.
func Execute() {
	rootCmd, logger, err := newRootCmd()

	if err != nil {
		fmt.Printf("Error creating root logger: %s", err)
		os.Exit(1)
	}

	err = rootCmd.Execute()

	if err == nil {
		os.Exit(0)
	} else if errors.Is(err, errViolations) {
		os.Exit(2)
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
