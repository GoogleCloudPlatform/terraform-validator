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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestAncestryPath(t *testing.T) {
	cases := []struct {
		name           string
		input          []*cloudresourcemanager.Ancestor
		expectedOutput string
	}{
		{
			name:           "Empty",
			input:          []*cloudresourcemanager.Ancestor{},
			expectedOutput: "",
		},
		{
			name: "ProjectOrganization",
			input: []*cloudresourcemanager.Ancestor{
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-prj",
						Type: "project",
					},
				},
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-org",
						Type: "organization",
					},
				},
			},
			expectedOutput: "organization/my-org/project/my-prj",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := ancestryPath(c.input)
			if output != c.expectedOutput {
				t.Errorf("expected output %q, got %q", c.expectedOutput, output)
			}
		})
	}
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
