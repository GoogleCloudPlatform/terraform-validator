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

// In order to interact with terraform-google-conversion, we need to be able to create
// "terraform resource data" that supports a very limited subset of the API actually
// used during the conversion process.
package google

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

// Compare to https://github.com/hashicorp/terraform-plugin-sdk/blob/97b4465/helper/schema/resource_data.go#L15
type FakeResourceData struct {
	kind   string
	schema *tfjson.Schema
	values map[string]interface{}
}

// Underlying implementation of getting a value from the given values and converting it
// to the real type. Returns two values: the value, and whether it exists.
// Note that the "value" will be an appropriate "zero value" if it does not exist.
// This does not currently take "computed"-ness into account.
func (d *FakeResourceData) get(key string) (interface{}, bool) {
	// try to get the value from our map
	value, exists := d.values[key]
	if !exists {
		value = nil
	}

	if value == nil {
		// Get the expected type for this value. We match the logic
		// from Terraform: use the most nested type that we can find.
		// We use this type to return a correct zero value.
		// https://github.com/hashicorp/terraform-plugin-sdk/blob/4d35f0240509689df3aedf16cd1ba90fb88a00aa/helper/schema/resource_data.go#L556
		address := strings.Split(key, ".")
		valueType := addressToType(address, d.schema.Block)
		if valueType == nil {
			// This means we couldn't find any type for this field.
			return nil, false
		}
		// Be consistent with zero values for Terraform
		// https://github.com/hashicorp/terraform-plugin-sdk/blob/e91cd2e5f2e60aa7d69f6011383f915216072c1b/helper/schema/schema.go#L2126
		switch *valueType {
		case cty.Bool:
			return false, exists
		case cty.Number:
			return 0.0, exists
		case cty.String:
			return "", exists
		}

		switch {
		case valueType.IsListType():
			return []interface{}{}, exists
		case valueType.IsSetType():
			return &schema.Set{}, exists
		case valueType.IsMapType():
			return map[string]interface{}{}, exists
		}

		panic(fmt.Sprintf("No zero value for type: %#v", valueType))
	}

	return value, exists
}

// Kind returns the type of resource (i.e. "google_storage_bucket").
func (d *FakeResourceData) Kind() string {
	return d.kind
}

// Id returns the ID of the resource from state.
func (d *FakeResourceData) Id() string {
	return ""
}

// Get reads a single field by key.
func (d *FakeResourceData) Get(name string) interface{} {
	val, _ := d.get(name)
	return val
}

// Get reads a single field by key and returns a boolean indicating
// whether the field exists.
func (d *FakeResourceData) GetOk(name string) (interface{}, bool) {
	return d.get(name)
}

// GetOkExists currently just calls GetOk but should be updated to support
// schema.ResourceData.GetOkExists functionality.
func (d *FakeResourceData) GetOkExists(name string) (interface{}, bool) {
	return d.GetOk(name)
}

// These methods are required by some mappers but we don't actually have (or need)
// implementations for them.
func (d *FakeResourceData) HasChange(string) bool             { return false }
func (d *FakeResourceData) Set(string, interface{}) error     { return nil }
func (d *FakeResourceData) SetId(string)                      {}
func (d *FakeResourceData) GetProviderMeta(interface{}) error { return nil }
func (d *FakeResourceData) Timeout(key string) time.Duration  { return time.Duration(1) }

func NewFakeResourceData(kind string, resourceSchema *tfjson.Schema, values map[string]interface{}) FakeResourceData {
	flattenedValues := map[string]interface{}{}
	flattenValues("", values, flattenedValues)
	return FakeResourceData{
		kind:   kind,
		schema: resourceSchema,
		values: flattenedValues,
	}
}

// Flatten resource values (which come from a json plan) to a one-level map.
// This simplifies getting the correct values out of the map later.
// For example:
//
//	{
//		"foo": {
//			"name": "value",
//          "should_do": false,
//		},
//	  "list": ["item1", "item2"],
//	}
//
// will be flattened to:
//
//	{
//		"foo.name": "value",
//      "should_do": false,
//		"list.#": 2,
//		"list.0": "item1",
//		"list.1": "item2",
//	}
func flattenValues(prefix string, src, dest map[string]interface{}) {
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenValues(prefix+k+".", child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
			dest[prefix+k+".#"] = len(child)
		default:
			dest[prefix+k] = v
		}
	}
}

// Given an address and a tfjson schema block, return the type that should
// be used for that address. This will generally be the deepest type found
// along the address.
func addressToType(address []string, block *tfjson.SchemaBlock) *cty.Type {
	currentBlock := block
	var currentType *cty.Type

	// Do this iteratively to avoid copying things
	for len(address) > 0 {
		currentAddress := address[0]
		address = address[1:]

		// Attributes include primitives and simple maps/lists
		if attribute, ok := currentBlock.Attributes[currentAddress]; ok {
			attributeType := attribute.AttributeType
			currentType = &attributeType

			// Primitives
			switch attributeType {
			case cty.Bool, cty.Number, cty.String:
				if len(address) > 0 {
					return nil
				}
				continue
			}

			// More complex types
			switch {
			case attributeType.IsSetType(), attributeType.IsListType():
				if len(address) == 0 {
					return currentType
				}

				if address[0] == "#" {
					return &cty.Number
				}

				elementType := attributeType.ElementType()
				return &elementType
			case attributeType.IsMapType():
				if len(address) == 0 {
					return currentType
				}

				// this may not be sufficient?
				elementType := attributeType.ElementType()
				return &elementType
			}

			// If we're not able to determine what type it is, that's an error
			// in this function.
			panic(fmt.Sprintf("Unhandled type (%s) for %s", attributeType.FriendlyName(), currentAddress))
		}

		// Nested blocks may result in updating currentBlock.
		// Note: the google provider currently only uses the following nesting modes:
		// - set
		// - list
		// - single
		// These all correspond to a list of values. We need to differentiate only
		// so far as calling functions (like expanders) might expect a specific data
		// type (such as *schema.Set).
		// https://github.com/hashicorp/terraform-json/blob/d1018bf93fd9c097133b0159ab8b3c0517a846c9/schemas.go#L131
		nestedBlock, ok := block.NestedBlocks[currentAddress]
		if !ok {
			break
		}
		if len(address) == 0 {
			var valueType cty.Type
			if nestedBlock.NestingMode == tfjson.SchemaNestingModeSet {
				valueType = cty.Set(cty.DynamicPseudoType)
			} else {
				valueType = cty.List(cty.DynamicPseudoType)
			}
			return &valueType
		}

		if address[0] == "#" {
			return &cty.Number
		}

		// We're going into the block. Set the current type to the
		// "element" type and advance to the next part of the address.
		mapType := cty.Map(cty.DynamicPseudoType)
		currentType = &mapType
		address = address[1:]

		currentBlock = nestedBlock.Block
	}

	// If the path doesn't exist, return nil
	return currentType
}
