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

package tfplan

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestComposeResources(t *testing.T) {
	plan := &terraform.Plan{
		State: &terraform.State{
			Modules: []*terraform.ModuleState{
				{
					Path: []string{"root"},
					Resources: map[string]*terraform.ResourceState{
						"google_test_resource.test": &terraform.ResourceState{
							Primary: &terraform.InstanceState{
								ID: "my-id",
								Attributes: map[string]string{
									"my_number": "42",
								},
							},
						},
					},
				},
			},
		},
		Diff: &terraform.Diff{
			Modules: []*terraform.ModuleDiff{
				{
					Path: []string{"root"},
					Resources: map[string]*terraform.InstanceDiff{
						"google_test_resource.test": &terraform.InstanceDiff{
							Attributes: map[string]*terraform.ResourceAttrDiff{
								"my_number": &terraform.ResourceAttrDiff{
									New: "42",
								},
							},
						},
					},
				},
			},
		},
	}

	schemas := map[string]*schema.Resource{
		"google_test_resource": {
			Schema: map[string]*schema.Schema{
				"my_number": {
					Type: schema.TypeInt,
				},
			},
		},
	}

	resources := ComposeResources(plan, schemas)
	require.Len(t, resources, 1)
	require.Equal(t, 42, resources[0].Get("my_number"))
	require.Equal(t, "my-id", resources[0].Id())
}
