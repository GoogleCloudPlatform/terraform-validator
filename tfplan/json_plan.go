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
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const maxChildModuleLevel = 10000

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
func ComposeTF12Resources(data []byte, schemas map[string]*schema.Resource) ([]Resource, error) {
	plan, err := readPlan(data)
	if err != nil {
		return nil, err
	}
	resources := readPlannedJSONResources(plan)
	return jsonToResources(resources, schemas), nil
}

// ComposeCurrentTF12Resources inspects a plan and returns the current resources that match the provided resource schema map.
// This works the same as ComposeTF12Resources but operates on current rather than planned resources.
func ComposeCurrentTF12Resources(data []byte, schemas map[string]*schema.Resource) ([]Resource, error) {
	plan, err := readPlan(data)
	if err != nil {
		return nil, err
	}
	resources := readCurrentJSONResources(plan)
	return jsonToResources(resources, schemas), nil
}

// jsonToResource converts the jsonResources to tfplan.Resource using the provided schemas.
// Any resources not supported by the schemas are silently skipped.
func jsonToResources(resources []jsonResource, schemas map[string]*schema.Resource) []Resource {
	var instances []Resource
	for _, r := range resources {
		sch, ok := schemas[r.Kind]
		if !ok {
			// Unsupported in given provider schema.
			continue
		}
		res := map[string]string{}
		var address []string
		attributes(r.Values, address, res, sch.Schema)
		reader := &schema.MapFieldReader{
			Map:    schema.BasicMapReader(res),
			Schema: sch.Schema,
		}
		instances = append(instances, Resource{
			Path: Fullpath{r.Kind, r.Name, r.Module},
			fieldGetter: &fieldGetter{
				rdr:    reader,
				schema: sch.Schema,
			},
		})
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

// attributes function takes json parsed JSON object (value param) and fill map[string]string with it's
// content (res param) for example JSON:
//
//	{
//		"foo": {
//			"name" : "value"
//		},
//	  "list": ["item1", "item2"]
//	}
//
// will be translated to map with following key/value set:
//
//	foo.name => "value"
//	list.# => 2
//	list.0 => "item1"
//	list.1 => "item2"
//
// Map above will be passed to schema.BasicMapReader that have all appropriate logic to read fields
// correctly during conversion to CAI.
func attributes(value interface{}, address []string, res map[string]string, schemas map[string]*schema.Schema) {
	schemaArr := addrToSchema(address, schemas)
	if len(schemaArr) == 0 {
		return
	}
	sch := schemaArr[len(schemaArr)-1]
	addr := strings.Join(address, ".")
	// int is special case, can't use handle it in main switch because number will be always parsed from JSON as float
	// need to identify it by schema.TypeInt and convert to int from int or float
	if sch.Type == schema.TypeInt && value != nil {
		switch value.(type) {
		case int:
			res[addr] = strconv.Itoa(value.(int))
		case float64:
			res[addr] = strconv.Itoa(int(value.(float64)))
		case float32:
			res[addr] = strconv.Itoa(int(value.(float32)))
		}
		return
	}

	switch value.(type) {
	case nil:
		defaultValue, err := sch.DefaultValue()
		if err != nil {
			panic(fmt.Sprintf("error getting default value for %v", address))
		}
		if defaultValue == nil {
			defaultValue = sch.ZeroValue()
		}
		attributes(defaultValue, address, res, schemas)
	case float64:
		res[addr] = strconv.FormatFloat(value.(float64), 'f', 6, 64)
	case float32:
		res[addr] = strconv.FormatFloat(value.(float64), 'f', 6, 32)
	case string:
		res[addr] = value.(string)
	case bool:
		res[addr] = strconv.FormatBool(value.(bool))
	case int:
		res[addr] = strconv.Itoa(value.(int))
	case []interface{}:
		arr := value.([]interface{})
		countAddr := addr + ".#"
		res[countAddr] = strconv.Itoa(len(arr))
		for i, e := range arr {
			addr := append(address, strconv.Itoa(i))
			attributes(e, addr, res, schemas)
		}
	case map[string]interface{}:
		m := value.(map[string]interface{})
		for k, v := range m {
			addr := append(address, k)
			attributes(v, addr, res, schemas)
		}
	case *schema.Set:
		set := value.(*schema.Set)
		attributes(set.List(), address, res, schemas)
	default:
		panic(fmt.Sprintf("unrecognized type %T", value))
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

// readCurrentJSONResources constructs jsonResources for the current state.
func readCurrentJSONResources(plan jsonPlan) []jsonResource {
	var result []jsonResource
	for _, c := range plan.ResourceChanges {
		// Ignore resources being created because they don't currently exist.
		if c.isCreate() {
			continue
		}
		result = append(result, jsonResource{
			Name:         c.Name,
			Address:      c.Address,
			Mode:         c.Mode,
			Kind:         c.Kind,
			ProviderName: c.ProviderName,
			Values:       c.Change.Before,
		})
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
