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

package google

import (
	"context"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-json"
	"github.com/GoogleCloudPlatform/terraform-validator/ancestrymanager"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const testProject = "test-project"

func newTestConverter() (*Converter, error) {
	ctx := context.Background()
	ancestry := ""
	project := testProject
	offline := true
	ancestryManager, err := ancestrymanager.New(context.Background(), project, ancestry, offline)
	if err != nil {
		return nil, errors.Wrap(err, "constructing resource manager client")
	}
	c, err := NewConverter(ctx, ancestryManager, project, "", offline)
	if err != nil {
		return nil, errors.Wrap(err, "building converter")
	}
	return c, nil
}

func TestSortByName(t *testing.T) {
	cases := []struct {
		name           string
		unsorted       []Asset
		expectedSorted []Asset
	}{
		{
			name:           "Empty",
			unsorted:       []Asset{},
			expectedSorted: []Asset{},
		},
		{
			name:           "BCAtoABC",
			unsorted:       []Asset{{Name: "b", Type: "b-type"}, {Name: "c", Type: "c-type"}, {Name: "a", Type: "a-type"}},
			expectedSorted: []Asset{{Name: "a", Type: "a-type"}, {Name: "b", Type: "b-type"}, {Name: "c", Type: "c-type"}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assets := c.unsorted
			sort.Sort(byName(assets))
			assert.EqualValues(t, c.expectedSorted, assets)
		})
	}
}

func TestAddResourceChanges_unknownResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_unknown.foo",
		Mode:         "managed",
		Type:         "google_unknown",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"change"},
			Before:  nil,
			After:   nil,
		},
	}
	c, err := newTestConverter()
	assert.Nil(t, err)
	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
}

func TestAddResourceChanges_unsupportedResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_unknown.foo",
		Mode:         "managed",
		Type:         "google_unsupported",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"change"},
			Before:  nil,
			After:   nil,
		},
	}
	c, err := newTestConverter()
	assert.Nil(t, err)

	// fake that this resource is known to the provider; it will never be "supported" by the
	// converter.
	c.schema.ResourcesMap[rc.Type] = c.schema.ResourcesMap["google_compute_disk"]

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
}

func TestAddResourceChanges_noopIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_compute_disk.foo",
		Mode:         "managed",
		Type:         "google_compute_disk",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"no-op"},
			Before:  nil,
			After:   nil,
		},
	}
	c, err := newTestConverter()
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
}

func TestAddResourceChanges_deleteProcessed(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_compute_disk.foo",
		Mode:         "managed",
		Type:         "google_compute_disk",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"delete"},
			Before:  map[string]interface{}{
				"project": testProject,
				"name":    "test-disk",
				"type":    "pd-ssd",
				"zone":    "us-central1-a",
				"image":   "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
				"labels": map[string]interface{}{
					"environment": "dev",
				},
				"physical_block_size_bytes": 4096,
			},
			After:   nil,
		},
	}
	c, err := newTestConverter()
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
}

func TestAddResourceChanges_createOrUpdateOrDeleteCreateProcessed(t *testing.T) {
	cases := []struct {
		name    string
		actions tfjson.Actions
	}{
		{
			name:    "Create",
			actions: tfjson.Actions{"create"},
		},
		{
			name:    "Update",
			actions: tfjson.Actions{"update"},
		},
		{
			name:    "DeleteCreate",
			actions: tfjson.Actions{"delete", "create"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rc := tfjson.ResourceChange{
				Address:      "whatever.google_compute_disk.foo",
				Mode:         "managed",
				Type:         "google_compute_disk",
				Name:         "foo",
				ProviderName: "google",
				Change: &tfjson.Change{
					Actions: c.actions,
					Before:  nil, // Ignore Before because it's unused
					After: map[string]interface{}{
						"project": testProject,
						"name":    "test-disk",
						"type":    "pd-ssd",
						"zone":    "us-central1-a",
						"image":   "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
						"labels": map[string]interface{}{
							"environment": "dev",
						},
						"physical_block_size_bytes": 4096,
					},
				},
			}
			c, err := newTestConverter()
			assert.Nil(t, err)

			err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
			assert.Nil(t, err)

			caiKey := "compute.googleapis.com/Disk//compute.googleapis.com/projects/test-project/zones/us-central1-a/disks/test-disk"
			assert.Contains(t, c.assets, caiKey)
		})
	}
}
