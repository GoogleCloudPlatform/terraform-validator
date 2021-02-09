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
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// jsonPlan structure used to parse Terraform 12 plan exported in json format by 'terraform show -json ./binary_plan.tfplan' command.
// https://www.terraform.io/docs/internals/json-format.html#plan-representation
type jsonPlan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       struct {
		// Valid actions values are:
		//    ["no-op"]
		//    ["create"]
		//    ["read"]
		//    ["update"]
		//    ["delete", "create"]
		//    ["create", "delete"]
		//    ["delete"]
		Actions []string               `json:"actions"`
		Before  map[string]interface{} `json:"before"`
		After   map[string]interface{} `json:"after"`
	} `json:"change"`
}

// isCreate returns true if the action on the resource is ["create"].
func (c *ResourceChange) isCreate() bool {
	return len(c.Change.Actions) == 1 && c.Change.Actions[0] == "create"
}

// compatibility shim until ResourceChange is expected by all callers.
func (c *ResourceChange) Kind() string {
	return c.Type
}

// ComposeTF12Resources is a thin wrapper around ReadResourceChanges as a compatibility shim.
func ComposeTF12Resources(data []byte, schemas map[string]*schema.Resource) ([]ResourceChange, error) {
	return ReadResourceChanges(data)
}

// ReadResourceChanges returns the list of resource changes from a json plan
func ReadResourceChanges(data []byte) ([]ResourceChange, error) {
	plan := jsonPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		return nil, errors.Wrap(err, "reading JSON plan")
	}

	return plan.ResourceChanges, nil
}
