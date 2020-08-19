package configs

import (
	"github.com/go-openapi/spec"
)

func refProperty(refURI string) *spec.Schema {
	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef(refURI),
		},
	}
}

var openAPISpecSchemaDefinitions = map[string]spec.Schema{
	"jsonschemaprops": {
		SchemaProps: spec.SchemaProps{
			Type:                 objectType,
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Properties: map[string]spec.Schema{
				"id":                   *spec.StringProperty(),
				"schema":               *spec.StringProperty(),
				"ref":                  *spec.StringProperty(),
				"description":          *spec.StringProperty(),
				"type":                 *spec.StringProperty(),
				"format":               *spec.StringProperty(),
				"title":                *spec.StringProperty(),
				"default":              *refProperty("#/definitions/json"),
				"maximum":              *spec.Float64Property(),
				"exclusiveMaximum":     *spec.BooleanProperty(),
				"minimum":              *spec.Float64Property(),
				"exclusiveMinimum":     *spec.BooleanProperty(),
				"maxLength":            *spec.Int64Property(),
				"minLength":            *spec.Int64Property(),
				"pattern":              *spec.StringProperty(),
				"maxItems":             *spec.Int64Property(),
				"minItems":             *spec.Int64Property(),
				"uniqueItems":          *spec.BooleanProperty(),
				"multipleOf":           *spec.Float64Property(),
				"enum":                 *spec.ArrayProperty(refProperty("#/definitions/json")),
				"maxProperties":        *spec.Int64Property(),
				"minProperties":        *spec.Int64Property(),
				"required":             *spec.ArrayProperty(spec.StringProperty()),
				"items":                *refProperty("#/definitions/jsonschemaprops"),
				"allOf":                *spec.ArrayProperty(refProperty("#/definitions/jsonschemaprops")),
				"oneOf":                *spec.ArrayProperty(refProperty("#/definitions/jsonschemaprops")),
				"anyOf":                *spec.ArrayProperty(refProperty("#/definitions/jsonschemaprops")),
				"not":                  *refProperty("#/definitions/jsonschemaprops"),
				"properties":           *spec.MapProperty(refProperty("#/definitions/jsonschemaprops")),
				"additionalProperties": *refProperty("#/definitions/jsonschemapropsorbool"),
				"patternProperties":    *spec.MapProperty(refProperty("#/definitions/jsonschemaprops")),
				"dependencies":         *spec.MapProperty(refProperty("#/definitions/jsonschemapropsorstringarray")),
				"additionalItems":      *refProperty("#/definitions/jsonschemapropsorbool"),
				"externalDocs":         *refProperty("#/definitions/externaldocumentation"),
				"example":              *refProperty("#/definitions/json"),
				"nullable":             *spec.BooleanProperty(),
			},
		},
	},
	"json": {SchemaProps: spec.SchemaProps{ID: "#json"}},
	"externaldocumentation": {
		SchemaProps: spec.SchemaProps{
			Type:                 objectType,
			AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			Properties: map[string]spec.Schema{
				"description": *spec.StringProperty(),
				"url":         *spec.StringProperty(),
			},
		},
	},
	"jsonschemapropsorstringarray": {
		SchemaProps: spec.SchemaProps{
			OneOf: []spec.Schema{*spec.ArrayProperty(spec.StringProperty()), *refProperty("#/definitions/jsonschemaprops")},
		},
	},
	"jsonschemapropsorbool": {
		SchemaProps: spec.SchemaProps{
			OneOf: []spec.Schema{*spec.BooleanProperty(), *refProperty("#/definitions/jsonschemaprops")},
		},
	},
}
