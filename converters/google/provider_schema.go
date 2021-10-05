package google

import (
	_ "embed"
	tfjson "github.com/hashicorp/terraform-json"
)

//go:embed google-3.86.0.json
var google3_86Bytes []byte

func getProviderSchema() *tfjson.ProviderSchema {
	schemas := &tfjson.ProviderSchemas{}
	err := schemas.UnmarshalJSON(google3_86Bytes)
	if err != nil {
		// We control this file. If it can't be loaded as json, we truly cannot and
		// should not proceed.
		panic(err)
	}

	return schemas.Schemas["registry.terraform.io/hashicorp/google"]
}
