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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// jsonPlan structure used to parse Terraform 12 plan exported in json format by 'terraform show -json ./binary_plan.tfplan' command.
type jsonPlan struct {
	PlannedValues struct {
		RootModules struct {
			Resources    []jsonResource
			ChildModules []struct {
				Address   string
				Resources []jsonResource
			} `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
}

// jsonResource represent single Terraform resource definition.
type jsonResource struct {
	Module       string
	Name         string
	Address      string
	Mode         string
	Kind         string `json:"type"`
	ProviderName string `json:"provider_name"`
	Values       map[string]interface{}
}

// jsonResourceFieldReader structure used to search retrieve fields from jsonResource.Values map of maps.
type jsonResourceFieldReader struct {
	Source jsonResource
	Schema map[string]*schema.Schema
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

		instances = append(instances, Resource{
			Path:        Fullpath{r.Kind, r.Name, r.Module},
			fieldGetter: newJSONFieldGetter(sch.Schema, r),
		})
	}
	return instances, nil
}

// ReadField are responsible for decoding fields out of data into
// the proper typed representation. jsonResourceFieldReader uses this to query data from json representation of
// Terraform 12 plan.
// See github.com/hashicorp/terraform/helper/schema.FieldReader interface for details
func (r jsonResourceFieldReader) ReadField(address []string) (schema.FieldReadResult, error) {
	addr := strings.Join(address, ".")
	schemaList := jsonAddrToSchema(address, r.Schema)
	if len(schemaList) == 0 {
		return schema.FieldReadResult{}, nil
	}

	var returnVal interface{}
	var ok bool
	sch := schemaList[len(schemaList)-1]

	if sch == nil {
		return schema.FieldReadResult{
			Value:  nil,
			Exists: false,
		}, nil
	}

	switch sch.Type {
	case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
		returnVal, ok = getValue(addr, r.Source, *sch)
	case schema.TypeSet:
		returnVal, ok = getValue(addr, r.Source, *sch)
		if returnVal == nil {
			returnVal = []interface{}{}
		}
		f := hashInterface
		returnVal = schema.NewSet(f, returnVal.([]interface{}))
	case schema.TypeMap:
		returnVal, ok = getValue(addr, r.Source, *sch)
		if returnVal == nil {
			returnVal = map[string]interface{}{}
		}
	default:
		panic(fmt.Sprintf("Unknown type: %s", sch.Type))
	}
	return schema.FieldReadResult{
		Value:  returnVal,
		Exists: ok,
	}, nil
}

func getValue(addr string, source jsonResource, sch schema.Schema) (value interface{}, ok bool) {
	value = source.Values[addr]
	if ok = value == nil; ok {
		value = sch.Default
	}
	return value, ok
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
		name := strings.SplitAfterN(module.Address, ".", 2)[1]
		for _, resource := range module.Resources {
			resource.Module = name
			result = append(result, resource)
		}
	}

	return result, nil
}

func newJSONFieldGetter(sch map[string]*schema.Schema, resource jsonResource) *fieldGetter {
	return &fieldGetter{
		rdr:    jsonResourceFieldReader{resource, sch},
		schema: sch,
	}
}

// jsonAddrToSchema returns shema.Shema object for string address representation, in most common case address is just
func jsonAddrToSchema(addr []string, schemaMap map[string]*schema.Schema) []*schema.Schema {
	var result []*schema.Schema
	key := strings.Join(addr, ".")
	result = append(result, schemaMap[key])
	return result
}

// Function hashInterface returns unique in ID for any given object.
// Function used by schema.Set instance, see github.com/hashicorp/terraform/helper/schema.SchemaSetFunc type for details.
func hashInterface(s interface{}) int {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(s)
	if err != nil {
		panic(fmt.Sprintf("error creating hashInterface function: %v", err))
	}
	h := fnv.New32a()
	_, err = h.Write(b.Bytes())
	if err != nil {
		panic(fmt.Sprintf("error creating hashInterface function: %v", err))
	}
	return int(h.Sum32())
}
