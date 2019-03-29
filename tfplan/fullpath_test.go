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

	"github.com/stretchr/testify/assert"
)

func TestFullpath(t *testing.T) {
	cases := []struct {
		name           string
		modulePath     []string
		resource       string
		expected       Fullpath
		expectedString string
	}{
		{
			name:       "Root",
			modulePath: []string{"root"},
			resource:   "my_resource_kind.my-resource-name",
			expected: Fullpath{
				Name:   "my-resource-name",
				Kind:   "my_resource_kind",
				Module: "root",
			},
			expectedString: "root.my_resource_kind.my-resource-name",
		},
		{
			name:       "Submodule",
			modulePath: []string{"root", "my-module-name"},
			resource:   "my_resource_kind.my-resource-name",
			expected: Fullpath{
				Name:   "my-resource-name",
				Kind:   "my_resource_kind",
				Module: "root.my-module-name",
			},
			expectedString: "root.my-module-name.my_resource_kind.my-resource-name",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fp := newFullpath(c.modulePath, c.resource)
			assert.Equal(t, c.expected, fp, "Fullpath")
			assert.Equal(t, c.expectedString, fp.String(), "String()")
		})
	}
}
