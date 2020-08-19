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

package configs

import (
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

var (
	objectType = spec.StringOrArray{"object"}
)

func mustMergeDefs(defMaps ...map[string]spec.Schema) map[string]spec.Schema {
	merged := map[string]spec.Schema{}
	for _, defMap := range defMaps {
		for name, def := range defMap {
			if _, duplicate := merged[name]; duplicate {
				panic(fmt.Sprintf("duplicate def %s", name))
			}
			merged[name] = def
		}
	}
	return merged
}

var constraintDefinitions = map[string]spec.Schema{
	"metadata": {
		SchemaProps: spec.SchemaProps{
			Type:                 objectType,
			Required:             []string{"name"},
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Properties: map[string]spec.Schema{
				"name":        *spec.StringProperty(),
				"labels":      *spec.MapProperty(spec.StringProperty()),
				"annotations": *spec.MapProperty(spec.StringProperty()),
			},
		},
	},
	"speccrd": {
		SchemaProps: spec.SchemaProps{
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Required:             []string{"spec"},
			Properties: map[string]spec.Schema{
				"spec": {
					SchemaProps: spec.SchemaProps{
						AdditionalProperties: &spec.SchemaOrBool{Allows: false},
						Required:             []string{"names"},
						Properties: map[string]spec.Schema{
							"names": {
								SchemaProps: spec.SchemaProps{
									AdditionalProperties: &spec.SchemaOrBool{Allows: false},
									Required:             []string{"kind"},
									Properties: map[string]spec.Schema{
										"kind": *spec.StringProperty(),
									},
								},
							},
							"validation": {
								SchemaProps: spec.SchemaProps{
									AdditionalProperties: &spec.SchemaOrBool{Allows: false},
									Required:             []string{"openAPIV3Schema"},
									Properties: map[string]spec.Schema{
										"openAPIV3Schema": *refProperty("#/definitions/jsonschemaprops"),
									},
								},
							},
						},
					},
				},
			},
		},
	},
	"alphav1spec": {
		SchemaProps: spec.SchemaProps{
			Type:                 objectType,
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Required:             []string{"crd", "targets"},
			Properties: map[string]spec.Schema{
				"crd": *refProperty("#/definitions/speccrd"),
				"targets": *spec.MapProperty(&spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:                 objectType,
						AdditionalProperties: &spec.SchemaOrBool{Allows: false},
						Required:             []string{"rego"},
						Properties: map[string]spec.Schema{
							"rego": *spec.StringProperty(),
							"libs": *spec.ArrayProperty(spec.StringProperty()),
						},
					},
				}),
			},
		},
	},
	"betav1spec": {
		SchemaProps: spec.SchemaProps{
			Type:                 objectType,
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Required:             []string{"crd", "targets"},
			Properties: map[string]spec.Schema{
				"crd": *refProperty("#/definitions/speccrd"),
				// convert to array here.
				"targets": *spec.ArrayProperty(&spec.Schema{
					VendorExtensible: spec.VendorExtensible{},
					SchemaProps: spec.SchemaProps{
						Type:                 objectType,
						AdditionalProperties: &spec.SchemaOrBool{Allows: false},
						Required:             []string{"target", "rego"},
						Properties: map[string]spec.Schema{
							"target": *spec.StringProperty(),
							"rego":   *spec.StringProperty(),
							"libs":   *spec.ArrayProperty(spec.StringProperty()),
						},
					},
				}),
			},
		},
	},
}

// configValidatorV1Alpha1Schema is the legacy config validator schema for CF-like templates.  Note that there's
// a subtle difference between this where "targets" is a map rather than an array.
var configValidatorV1Alpha1Schema = spec.Schema{
	SchemaProps: spec.SchemaProps{
		Definitions:          mustMergeDefs(openAPISpecSchemaDefinitions, constraintDefinitions),
		AdditionalProperties: &spec.SchemaOrBool{Allows: false},
		Properties: map[string]spec.Schema{
			"apiVersion": *spec.StringProperty(),
			"kind":       *spec.StringProperty(),
			"metadata":   *refProperty("#/definitions/metadata"),
			"spec":       *refProperty("#/definitions/alphav1spec"),
		},
	},
}

var configValidatorV1Alpha1SchemaValidator = validate.NewSchemaValidator(
	&configValidatorV1Alpha1Schema, nil, "", strfmt.Default)

var configValidatorV1Beta1Schema = spec.Schema{
	SchemaProps: spec.SchemaProps{
		Definitions:          mustMergeDefs(openAPISpecSchemaDefinitions, constraintDefinitions),
		AdditionalProperties: &spec.SchemaOrBool{Allows: false},
		Properties: map[string]spec.Schema{
			"apiVersion": *spec.StringProperty(),
			"kind":       *spec.StringProperty(),
			"metadata":   *refProperty("#/definitions/metadata"),
			"spec":       *refProperty("#/definitions/betav1spec"),
		},
	},
}

var configValidatorV1Beta1SchemaValidator = validate.NewSchemaValidator(
	&configValidatorV1Beta1Schema, nil, "", strfmt.Default)
