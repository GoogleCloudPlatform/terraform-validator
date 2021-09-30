package google

import (
	_ "embed"
	tfjson "github.com/hashicorp/terraform-json"
)

//go:embed provider_schema.json
var providerSchemaBytes []byte

func getProviderSchema() *tfjson.ProviderSchema {
	schemas := &tfjson.ProviderSchemas{}
	err := schemas.UnmarshalJSON(providerSchemaBytes)
	if err != nil {
		// We control this file. If it can't be loaded as json, we truly cannot and
		// should not proceed.
		panic(err)
	}

	return schemas.Schemas["registry.terraform.io/hashicorp/google"]
}
