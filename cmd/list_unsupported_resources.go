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
	"sort"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	provider "github.com/hashicorp/terraform-provider-google/google"
	"github.com/spf13/cobra"
)

type listUnsupportedResourcesOptions struct{}

func newListUnsupportedResourcesCmd() *cobra.Command {
	o := listUnsupportedResourcesOptions{}

	cmd := &cobra.Command{
		Use:   "list-unsupported-resources",
		Short: "List unsupported terraform resources.",
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
		Hidden: true,
	}
	return cmd
}

func (o *listUnsupportedResourcesOptions) run() error {
	// Get a map of supported terraform resources
	converters := resources.ResourceConverters()

	// Get a sorted list of unsupported resources
	schema := provider.Provider()
	unsupported := make([]string, 0, len(schema.ResourcesMap))
	for k := range schema.ResourcesMap {
		if _, ok := converters[k]; !ok {
			unsupported = append(unsupported, k)
		}
	}
	sort.Strings(unsupported)

	// go through and print
	for _, resource := range unsupported {
		fmt.Println(resource)
	}

	return nil
}
