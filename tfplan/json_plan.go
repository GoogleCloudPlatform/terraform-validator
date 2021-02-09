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
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const maxChildModuleLevel = 10000

// jsonPlan structure used to parse Terraform 12 plan exported in json format by 'terraform show -json ./binary_plan.tfplan' command.
// https://www.terraform.io/docs/internals/json-format.html#plan-representation
type jsonPlan struct {
	PlannedValues struct {
		RootModules struct {
			Resources    []jsonResource `json:"resources"`
			ChildModules []childModule  `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []jsonResourceChange `json:"resource_changes"`
}

type jsonResourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Kind         string `json:"type"`
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
func (c *jsonResourceChange) isCreate() bool {
	return len(c.Change.Actions) == 1 && c.Change.Actions[0] == "create"
}

type childModule struct {
	Address      string         `json:"address"`
	Resources    []jsonResource `json:"resources"`
	ChildModules []childModule  `json:"child_modules"`
}

// jsonResource represent single Terraform resource definition.
type jsonResource struct {
	Module       string                 `json:"module"`
	Name         string                 `json:"name"`
	Address      string                 `json:"address"`
	Mode         string                 `json:"mode"`
	Kind         string                 `json:"type"`
	ProviderName string                 `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}

// ComposeTF12Resources inspects a plan and returns the planned resources that match the provided resource schema map.
// ComposeTF12Resources works in a same way as tfplan.ComposeResources and returns array of tfplan.Resource
func ComposeTF12Resources(data []byte, schemas map[string]*schema.Resource) ([]google.FakeResourceData, error) {
	plan, err := readPlan(data)
	if err != nil {
		return nil, err
	}
	resources := readPlannedJSONResources(plan)
	return jsonToResources(resources, schemas), nil
}

// jsonToResource converts the jsonResources to tfplan.Resource using the provided schemas.
// Any resources not supported by the schemas are silently skipped.
func jsonToResources(resources []jsonResource, schemas map[string]*schema.Resource) []google.FakeResourceData {
	var instances []google.FakeResourceData
	for _, r := range resources {
		sch, ok := schemas[r.Kind]
		if !ok {
			// Unsupported in given provider schema.
			continue
		}

		instances = append(instances, google.NewFakeResourceData(
			r.Kind,
			sch.Schema,
			r.Values,
		))
	}
	return instances
}

func jsonResourceFromChange(c jsonResourceChange) jsonResource {
	return jsonResource{
		Name:         c.Name,
		Address:      c.Address,
		Mode:         c.Mode,
		Kind:         c.Kind,
		ProviderName: c.ProviderName,
		Values:       c.Change.Before,
	}
}

// readPlan converts the raw bytes into a jsonPlan
func readPlan(data []byte) (jsonPlan, error) {
	plan := jsonPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		return jsonPlan{}, errors.Wrap(err, "read resources")
	}
	return plan, nil
}

// readPlannedJSONResources reads a jsonPlan and returns an array of all
// planned jsonResources.
// This includes resources from both from root and child modules.
func readPlannedJSONResources(plan jsonPlan) []jsonResource {
	var result []jsonResource
	for _, resource := range plan.PlannedValues.RootModules.Resources {
		resource.Module = "root"
		result = append(result, resource)
	}
	for _, module := range plan.PlannedValues.RootModules.ChildModules {
		result = append(result, resourcesFromModule(&module, 0)...)
	}
	return result
}

func resourcesFromModule(module *childModule, level int) []jsonResource {
	if level > maxChildModuleLevel {
		log.Printf("The configuration has more than %d level of modules. Modules with a depth more than %d will be ignored.", maxChildModuleLevel, maxChildModuleLevel)
		return nil
	}
	pcs := strings.SplitAfterN(module.Address, ".", 2)
	if len(pcs) < 2 {
		return nil
	}
	name := pcs[1]
	var result []jsonResource
	for _, resource := range module.Resources {
		resource.Module = name
		result = append(result, resource)
	}
	for _, c := range module.ChildModules {
		result = append(result, resourcesFromModule(&c, level+1)...)
	}
	return result
}
