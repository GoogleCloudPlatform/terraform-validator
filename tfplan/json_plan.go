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
	"github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func IsCreate(rc *tfjson.ResourceChange) bool {
	return len(rc.Change.Actions) == 1 && rc.Change.Actions[0] == "create"
}

func IsUpdate(rc *tfjson.ResourceChange) bool {
	return len(rc.Change.Actions) == 1 && rc.Change.Actions[0] == "update"
}

func IsDeleteCreate(rc *tfjson.ResourceChange) bool {
	return len(rc.Change.Actions) == 2 && rc.Change.Actions[0] == "delete"
}

func IsDelete(rc *tfjson.ResourceChange) bool {
	return len(rc.Change.Actions) == 1 && rc.Change.Actions[0] == "delete"
}

// compatibility shim until ResourceChange is expected by all callers.
func Kind(rc *tfjson.ResourceChange) string {
	return rc.Type
}

// ComposeTF12Resources is a thin wrapper around ReadResourceChanges as a compatibility shim.
// It needs to return a slice of non-pointers to the same type that Converter.AddResource
// takes a pointer to.
func ComposeTF12Resources(data []byte, _ map[string]*schema.Resource) ([]tfjson.ResourceChange, error) {
	rcps, err := ReadResourceChanges(data)

	if err != nil {
		return nil, err
	}

	var rcs []tfjson.ResourceChange
	for _, rcp := range rcps {
		rcs = append(rcs, *rcp)
	}
	return rcs, nil
}

// ReadResourceChanges returns the list of resource changes from a json plan
func ReadResourceChanges(data []byte) ([]*tfjson.ResourceChange, error) {
	plan := tfjson.Plan{}
	err := plan.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.Wrap(err, "reading JSON plan")
	}

	err = plan.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "validating JSON plan")
	}

	return plan.ResourceChanges, nil
}
