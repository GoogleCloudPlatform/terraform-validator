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
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// newFieldGetter accepts a resource schema map (field name --> Schema),
// a state instance, and a diff and returns a type with field getter methods.
func newFieldGetter(sch map[string]*schema.Schema, state *terraform.InstanceState, diff *terraform.InstanceDiff) *fieldGetter {
	if state == nil && diff == nil {
		panic("both state and diff should not be nil")
	}

	var baseState map[string]string
	if state != nil {
		baseState = state.Attributes
	} else {
		baseState = make(map[string]string)
	}

	stateReader := &schema.MapFieldReader{
		Map:    schema.BasicMapReader(baseState),
		Schema: sch,
	}

	var rdr schema.FieldReader
	if diff == nil {
		rdr = stateReader
	} else {
		rdr = &schema.DiffFieldReader{
			Diff:   diff,
			Source: stateReader,
			Schema: sch,
		}
	}

	return &fieldGetter{
		rdr:    rdr,
		schema: sch,
		state:  state,
	}
}

// fieldGetter exposes the ability to get terraform resource instance
// field values.
type fieldGetter struct {
	rdr    schema.FieldReader
	schema map[string]*schema.Schema
	state  *terraform.InstanceState
}

// Id returns the ID of the resource from state.
func (g *fieldGetter) Id() string {
	if g.state != nil {
		return g.state.ID
	}
	return ""
}

// Get reads a single field by key.
func (g *fieldGetter) Get(name string) interface{} {
	val, _ := g.GetOk(name)
	return val
}

// Get reads a single field by key and returns a boolean indicating
// whether the field exists.
func (g *fieldGetter) GetOk(name string) (interface{}, bool) {
	res, err := g.rdr.ReadField(strings.Split(name, "."))
	if err != nil {
		return nil, false
	}

	addr := strings.Split(name, ".")
	schemaPath := addrToSchema(addr, g.schema)
	if len(schemaPath) == 0 {
		return nil, false
	}

	return res.ValueOrZero(schemaPath[len(schemaPath)-1]), res.Exists && !res.Computed
}

// GetOkExists currently just calls GetOk but should be updated to support
// schema.ResourceData.GetOkExists functionality.
func (g *fieldGetter) GetOkExists(name string) (interface{}, bool) {
	return g.GetOk(name)
}

// addrToSchema finds the final element schema for the given address
// and the given schema. It returns all the schemas that led to the final
// schema. These are in order of the address (out to in).
// NOTE: This function was copied from the terraform library:
// github.com/hashicorp/terraform/helper/schema/field_reader.go
func addrToSchema(addr []string, schemaMap map[string]*schema.Schema) []*schema.Schema {
	const typeObject = 999

	current := &schema.Schema{
		Type: typeObject,
		Elem: schemaMap,
	}

	// If we aren't given an address, then the user is requesting the
	// full object, so we return the special value which is the full object.
	if len(addr) == 0 {
		return []*schema.Schema{current}
	}

	result := make([]*schema.Schema, 0, len(addr))
	for len(addr) > 0 {
		k := addr[0]
		addr = addr[1:]

	REPEAT:
		// We want to trim off the first "typeObject" since its not a
		// real lookup that people do. i.e. []string{"foo"} in a structure
		// isn't {typeObject, typeString}, its just a {typeString}.
		if len(result) > 0 || current.Type != typeObject {
			result = append(result, current)
		}

		switch t := current.Type; t {
		case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
			if len(addr) > 0 {
				return nil
			}
		case schema.TypeList, schema.TypeSet:
			isIndex := len(addr) > 0 && addr[0] == "#"

			switch v := current.Elem.(type) {
			case *schema.Resource:
				current = &schema.Schema{
					Type: typeObject,
					Elem: v.Schema,
				}
			case *schema.Schema:
				current = v
			case schema.ValueType:
				current = &schema.Schema{Type: v}
			default:
				// we may not know the Elem type and are just looking for the
				// index
				if isIndex {
					break
				}

				if len(addr) == 0 {
					// we've processed the address, so return what we've
					// collected
					return result
				}

				if len(addr) == 1 {
					if _, err := strconv.Atoi(addr[0]); err == nil {
						// we're indexing a value without a schema. This can
						// happen if the list is nested in another schema type.
						// Default to a TypeString like we do with a map
						current = &schema.Schema{Type: schema.TypeString}
						break
					}
				}

				return nil
			}

			// If we only have one more thing and the next thing
			// is a #, then we're accessing the index which is always
			// an int.
			if isIndex {
				current = &schema.Schema{Type: schema.TypeInt}
				break
			}

		case schema.TypeMap:
			if len(addr) > 0 {
				switch v := current.Elem.(type) {
				case schema.ValueType:
					current = &schema.Schema{Type: v}
				default:
					// maps default to string values. This is all we can have
					// if this is nested in another list or map.
					current = &schema.Schema{Type: schema.TypeString}
				}
			}
		case typeObject:
			// If we're already in the object, then we want to handle Sets
			// and Lists specially. Basically, their next key is the lookup
			// key (the set value or the list element). For these scenarios,
			// we just want to skip it and move to the next element if there
			// is one.
			if len(result) > 0 {
				lastType := result[len(result)-2].Type
				if lastType == schema.TypeSet || lastType == schema.TypeList {
					if len(addr) == 0 {
						break
					}

					k = addr[0]
					addr = addr[1:]
				}
			}

			m := current.Elem.(map[string]*schema.Schema)
			val, ok := m[k]
			if !ok {
				return nil
			}

			current = val
			goto REPEAT
		}
	}

	return result
}
