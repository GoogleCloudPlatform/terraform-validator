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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

const maxChildModuleLevel = 10000

// jsonPlan structure used to parse Terraform 12 plan exported in json format by 'terraform show -json ./binary_plan.tfplan' command.
type jsonPlan struct {
	PlannedValues struct {
		RootModules struct {
			Resources    []jsonResource `json:"resources"`
			ChildModules []childModule  `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []struct {
		Address string `json:"address"`
		Change  struct {
			AfterUnknown map[string]interface{} `json:"after_unknown"`
		} `json:"change"`
	} `json:"resource_changes"`
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
	resources, err := readJSONResources(data)
	if err != nil {
		return nil, errors.Wrap(err, "read resources")
	}
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
	return instances, nil
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

// readJSONResources unmarshal json data to go struct
// and returns array of all jsonResources both from root and child modules.
func readJSONResources(data []byte) ([]jsonResource, error) {
	plan := jsonPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		return nil, err
	}
	var result []jsonResource
	for _, resource := range plan.PlannedValues.RootModules.Resources {
		resource.Module = "root"
		result = append(result, resource)
	}
	for _, module := range plan.PlannedValues.RootModules.ChildModules {
		result = append(result, resourcesFromModule(&module, 0)...)
	}

	resources := make(map[string]*jsonResource)
	for i, resource := range result {
		resources[resource.Address] = &result[i]
	}

	// Inject computed values into resources
	for _, change := range plan.ResourceChanges {
		resource, ok := resources[change.Address]
		if !ok {
			continue
		}
		err = injectComputedMap(change.Change.AfterUnknown, resource.Values)
	}

	return result, nil
}

func injectComputedMap(changes map[string]interface{}, resource map[string]interface{}) error {
	for key, val := range changes {
		switch val.(type) {
		case bool:
			if val.(bool) {
				resource[key] = "_computed_"
			}
		case []interface{}:
			resourceValue, ok := resource[key]
			if !ok {
				resourceValue = make([]interface{}, 0)
				log.Printf("Injected empty array at %v", key)
			}
			injectComputedArray(val.([]interface{}), resourceValue.([]interface{}))
		default:
			log.Printf("Unknown computed value found (key = %v): %v", key, val)
		}
	}
	return nil
}

func injectComputedArray(changeValue []interface{}, resourceValue []interface{}) error {
	for i, item := range changeValue {
		if len(resourceValue) <= i {
			resourceValue = append(resourceValue, make(map[string]interface{}))
			log.Printf("Injected empty value at %v", i)
		}
		resourceItem := resourceValue[i]
		injectComputedMap(item.(map[string]interface{}), resourceItem.(map[string]interface{}))
	}
	return nil
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
