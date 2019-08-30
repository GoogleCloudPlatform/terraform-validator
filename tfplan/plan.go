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
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Resource is the terraform representation of a resource.
type Resource struct {
	Path Fullpath
	*fieldGetter
}

// Kind returns the type of resource (i.e. "google_storage_bucket").
func (r *Resource) Kind() string {
	return r.Path.Kind
}

// Provider derives the provider from the resource prefix.
// i.e. resource="google_storage_bucket" --> provider="google".
func (r *Resource) Provider() string {
	// NOTE: In order to differentiate between "google" and "google-beta"
	// the resource fields will need to be inspected (i.e. "provider").
	return strings.Split(r.Kind(), "_")[0]
}

// ComposeResources inspects a plan and returns the planned resources that match
// the provided resource schema map.
func ComposeResources(plan *terraform.Plan, schemas map[string]*schema.Resource) []Resource {
	instances := make([]Resource, 0)
	for path, sd := range mergeStateDiffs(plan) {
		schm, ok := schemas[path.Kind]
		if !ok {
			// Unsupported in given provider schm.
			continue
		}

		instances = append(instances, Resource{
			Path:        path,
			fieldGetter: newFieldGetter(schm.Schema, sd.State, sd.Diff),
		})
	}

	return instances
}

type stateDiff struct {
	State *terraform.InstanceState
	Diff  *terraform.InstanceDiff
}

// mergeStateDiffs creates a map of all the resources contained in a plan
// by stitching together the state and diff to create a unified view of the
// future resource that is planned.
func mergeStateDiffs(plan *terraform.Plan) map[Fullpath]stateDiff {
	res := make(map[Fullpath]stateDiff)

	for _, mod := range plan.State.Modules {
		for name, state := range mod.Resources {
			fp := newFullpath(mod.Path, name)
			res[fp] = stateDiff{
				State: state.Primary,
			}
		}
	}

	for _, mod := range plan.Diff.Modules {
		for name, diff := range mod.Resources {
			fp := newFullpath(mod.Path, name)
			if data, ok := res[fp]; ok {
				data.Diff = diff
				res[fp] = data
			} else {
				res[fp] = stateDiff{
					Diff: diff,
				}
			}
		}
	}

	return res
}
