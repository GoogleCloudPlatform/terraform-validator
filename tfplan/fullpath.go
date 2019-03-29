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

import "strings"

// Fullpath contains both the path of the module and the resource within
// the module. All fields are comparable so it can be used as a map key.
type Fullpath struct {
	Kind   string
	Name   string
	Module string
}

// String returns a concatenated full resource path.
func (fp Fullpath) String() string {
	return strings.Join([]string{fp.Module, fp.Kind, fp.Name}, ".")
}

// newFullpath makes it easy to create a fullpath from the way the parts are
// stored in a terraform.Plan ([]string for module path and "kind.name" for
// the resource).
func newFullpath(modPath []string, res string) Fullpath {
	splitRes := strings.Split(res, ".")
	return Fullpath{
		Kind:   splitRes[0],
		Name:   splitRes[1],
		Module: strings.Join(modPath, "."),
	}
}
