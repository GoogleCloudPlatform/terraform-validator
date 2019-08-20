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
	"strconv"
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
func (r jsonResourceFieldReader) ReadField1(address []string) (schema.FieldReadResult, error) {
	addr := strings.Join(address, ".")
	schemaList := addrToSchema(address, r.Schema)
	if len(schemaList) == 0 {
		return schema.FieldReadResult{}, nil
	}
	sch := schemaList[len(schemaList)-1]
	return readFieldInternal(addr, r.Source.Values, sch)
}

func (r *jsonResourceFieldReader) ReadField(address []string) (schema.FieldReadResult, error) {
	k := strings.Join(address, ".")
	schemaList := addrToSchema(address, r.Schema)
	if len(schemaList) == 0 {
		return schema.FieldReadResult{}, nil
	}

	sch := schemaList[len(schemaList)-1]
	switch sch.Type {
	case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
		return r.readPrimitive(address, sch)
	case schema.TypeList:
		return readListField(r, address, sch)
	case schema.TypeMap:
		return r.readMap(k, sch)
	case schema.TypeSet:
		return r.readSet(address, sch)
	/*
		case typeObject:
			return readObjectField(r, address, schema.Elem.(map[string]*Schema))
	*/

	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}
}

func (r *jsonResourceFieldReader) readPrimitive(
	address []string, sch *schema.Schema) (schema.FieldReadResult, error) {
	k := strings.Join(address, ".")
	result, ok := r.Source.Values[k]
	if !ok {
		return schema.FieldReadResult{}, nil
	}

	returnVal, err := stringToPrimitive(result.(string), false, sch)
	if err != nil {
		return schema.FieldReadResult{}, err
	}

	return schema.FieldReadResult{
		Value:  returnVal,
		Exists: true,
	}, nil
}

// readListField is a generic method for reading a list field out of a
// a FieldReader. It does this based on the assumption that there is a key
// "foo.#" for a list "foo" and that the indexes are "foo.0", "foo.1", etc.
// after that point.
func readListField(
	r schema.FieldReader, addr []string, sch *schema.Schema) (schema.FieldReadResult, error) {
	addrPadded := make([]string, len(addr)+1)
	copy(addrPadded, addr)
	addrPadded[len(addrPadded)-1] = "#"

	// Get the number of elements in the list
	countResult, err := r.ReadField(addrPadded)
	if err != nil {
		return schema.FieldReadResult{}, err
	}
	if !countResult.Exists {
		// No count, means we have no list
		countResult.Value = 0
	}

	// If we have an empty list, then return an empty list
	if countResult.Computed || countResult.Value.(int) == 0 {
		return schema.FieldReadResult{
			Value:    []interface{}{},
			Exists:   countResult.Exists,
			Computed: countResult.Computed,
		}, nil
	}

	// Go through each count, and get the item value out of it
	result := make([]interface{}, countResult.Value.(int))
	for i, _ := range result {
		is := strconv.FormatInt(int64(i), 10)
		addrPadded[len(addrPadded)-1] = is
		rawResult, err := r.ReadField(addrPadded)
		if err != nil {
			return schema.FieldReadResult{}, err
		}
		if !rawResult.Exists {
			// This should never happen, because by the time the data
			// gets to the FieldReaders, all the defaults should be set by
			// Schema.
			rawResult.Value = nil
		}

		result[i] = rawResult.Value
	}

	return schema.FieldReadResult{
		Value:  result,
		Exists: true,
	}, nil
}

func (r *jsonResourceFieldReader) readMap(k string, sch *schema.Schema) (schema.FieldReadResult, error) {
	result := make(map[string]interface{})
	resultSet := false

	// If the name of the map field is directly in the map with an
	// empty string, it means that the map is being deleted, so mark
	// that is is set.
	if v, ok := r.Source.Values[k]; ok && v == "" {
		resultSet = true
	}

	prefix := k + "."
	for k, v := range r.Source.Values {
		if strings.HasPrefix(k, prefix) {
			resultSet = true

			key := k[len(prefix):]
			if key != "%" && key != "#" {
				result[key] = v
			}
		}
	}

	err := mapValuesToPrimitive(k, result, sch)
	if err != nil {
		return schema.FieldReadResult{}, nil
	}

	var resultVal interface{}
	if resultSet {
		resultVal = result
	}

	return schema.FieldReadResult{
		Value:  resultVal,
		Exists: resultSet,
	}, nil
}

func readFieldInternal(addr string, values map[string]interface{}, sch *schema.Schema) (schema.FieldReadResult, error) {

	var returnVal interface{}
	var ok bool

	if sch == nil {
		return schema.FieldReadResult{
			Value:  nil,
			Exists: false,
		}, nil
	}

	switch sch.Type {
	case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
		returnVal, ok = getValue(addr, values, *sch)
	case schema.TypeSet:
		returnVal, ok = getValue(addr, values, *sch)
		if returnVal == nil {
			returnVal = []interface{}{}
		}
		f := hashInterface
		returnVal = schema.NewSet(f, returnVal.([]interface{}))
	case schema.TypeMap:
		returnVal, ok = getValue(addr, values, *sch)
		if returnVal == nil {
			returnVal = map[string]interface{}{}
		}
	case schema.TypeList:
		returnVal, ok = getValue(addr, values, *sch)
		if returnVal == nil {
			returnVal = []interface{}{}
		}
	default:
		panic(fmt.Sprintf("Unknown type: %s", sch.Type))
	}
	return schema.FieldReadResult{
		Value:  returnVal,
		Exists: ok,
	}, nil

}

func getValue(addr string, values map[string]interface{}, sch schema.Schema) (value interface{}, ok bool) {
	value = values[addr]
	if ok = value != nil; !ok {
		value = sch.Default
	}

	if ok && sch.Type == schema.TypeList {
		for _, item := range value.([]interface{}) {
			result := item.(map[string]interface{})
			res := sch.Elem.(*schema.Resource)
			for attrName, attrSchema := range res.Schema {
				readResult, _ := readFieldInternal(attrName, result, attrSchema)
				result[attrName] = readResult.ValueOrZero(attrSchema)
			}
		}
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

func (r jsonResourceFieldReader) readSet(
	address []string, sch *schema.Schema) (schema.FieldReadResult, error) {
	// Get the number of elements in the list
	countRaw, err := r.readPrimitive(
		append(address, "#"), &schema.Schema{Type: schema.TypeInt})
	if err != nil {
		return schema.FieldReadResult{}, err
	}
	if !countRaw.Exists {
		// No count, means we have no list
		countRaw.Value = 0
	}

	// Create the set that will be our result
	set := sch.ZeroValue().(*schema.Set)

	// If we have an empty list, then return an empty list
	if countRaw.Computed || countRaw.Value.(int) == 0 {
		return schema.FieldReadResult{
			Value:    set,
			Exists:   countRaw.Exists,
			Computed: countRaw.Computed,
		}, nil
	}

	// Go through the map and find all the set items
	prefix := strings.Join(address, ".") + "."
	countExpected := countRaw.Value.(int)
	countActual := make(map[string]struct{})
	for k, _ := range r.Source.Values {
		if !strings.HasPrefix(k, prefix) {
			return true
		}
		if strings.HasPrefix(k, prefix+"#") {
			// Ignore the count field
			return true
		}

		// Split the key, since it might be a sub-object like "idx.field"
		parts := strings.Split(k[len(prefix):], ".")
		idx := parts[0]

		var raw schema.FieldReadResult
		raw, err = r.ReadField(append(address, idx))
		if err != nil {
			return false
		}
		if !raw.Exists {
			// This shouldn't happen because we just verified it does exist
			panic("missing field in set: " + k + "." + idx)
		}

		set.Add(raw.Value)

		// Due to the way multimap readers work, if we've seen the number
		// of fields we expect, then exit so that we don't read later values.
		// For example: the "set" map might have "ports.#", "ports.0", and
		// "ports.1", but the "state" map might have those plus "ports.2".
		// We don't want "ports.2"
		countActual[idx] = struct{}{}
		if len(countActual) >= countExpected {
			return false
		}

		return true
	}
	completed := r.Map.Range(func(k, _ string) bool {
		if !strings.HasPrefix(k, prefix) {
			return true
		}
		if strings.HasPrefix(k, prefix+"#") {
			// Ignore the count field
			return true
		}

		// Split the key, since it might be a sub-object like "idx.field"
		parts := strings.Split(k[len(prefix):], ".")
		idx := parts[0]

		var raw schema.FieldReadResult
		raw, err = r.ReadField(append(address, idx))
		if err != nil {
			return false
		}
		if !raw.Exists {
			// This shouldn't happen because we just verified it does exist
			panic("missing field in set: " + k + "." + idx)
		}

		set.Add(raw.Value)

		// Due to the way multimap readers work, if we've seen the number
		// of fields we expect, then exit so that we don't read later values.
		// For example: the "set" map might have "ports.#", "ports.0", and
		// "ports.1", but the "state" map might have those plus "ports.2".
		// We don't want "ports.2"
		countActual[idx] = struct{}{}
		if len(countActual) >= countExpected {
			return false
		}

		return true
	})
	if !completed && err != nil {
		return schema.FieldReadResult{}, err
	}

	return schema.FieldReadResult{
		Value:  set,
		Exists: true,
	}, nil
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

func stringToPrimitive(
	value string, computed bool, sch *schema.Schema) (interface{}, error) {
	var returnVal interface{}
	switch sch.Type {
	case schema.TypeBool:
		if value == "" {
			returnVal = false
			break
		}
		if computed {
			break
		}

		v, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}

		returnVal = v
	case schema.TypeFloat:
		if value == "" {
			returnVal = 0.0
			break
		}
		if computed {
			break
		}

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		returnVal = v
	case schema.TypeInt:
		if value == "" {
			returnVal = 0
			break
		}
		if computed {
			break
		}

		v, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return nil, err
		}

		returnVal = int(v)
	case schema.TypeString:
		returnVal = value
	default:
		panic(fmt.Sprintf("Unknown type: %s", sch.Type))
	}

	return returnVal, nil
}

// convert map values to the proper primitive type based on schema.Elem
func mapValuesToPrimitive(k string, m map[string]interface{}, sch *schema.Schema) error {
	elemType, err := getValueType(k, sch)
	if err != nil {
		return err
	}

	switch elemType {
	case schema.TypeInt, schema.TypeFloat, schema.TypeBool:
		for k, v := range m {
			vs, ok := v.(string)
			if !ok {
				continue
			}

			v, err := stringToPrimitive(vs, false, &schema.Schema{Type: elemType})
			if err != nil {
				return err
			}

			m[k] = v
		}
	}
	return nil
}

func getValueType(k string, sch *schema.Schema) (schema.ValueType, error) {
	if sch.Elem == nil {
		return schema.TypeString, nil
	}
	if vt, ok := sch.Elem.(schema.ValueType); ok {
		return vt, nil
	}

	// If a Schema is provided to a Map, we use the Type of that schema
	// as the type for each element in the Map.
	if s, ok := sch.Elem.(*schema.Schema); ok {
		return s.Type, nil
	}

	if _, ok := sch.Elem.(*Resource); ok {
		// TODO: We don't actually support this (yet)
		// but silently pass the validation, until we decide
		// how to handle nested structures in maps
		return schema.TypeString, nil
	}
	return 0, fmt.Errorf("%s: unexpected map value type: %#v", k, schema.Elem)
}
