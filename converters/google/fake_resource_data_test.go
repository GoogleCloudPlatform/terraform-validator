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
package google

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func getComputeDiskSchema() *tfjson.Schema {
	schema := getProviderSchema()
	return schema.ResourceSchemas["google_compute_disk"]
}

func TestFakeResourceData_kind(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	assert.Equal(t, d.Kind(), "google_compute_disk")
}

func TestFakeResourceData_id(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	assert.Equal(t, d.Id(), "")
}

func TestFakeResourceData_get(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	assert.Equal(t, d.Get("name"), "test-disk")
}

func TestFakeResourceData_getOkOk(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	res, ok := d.GetOk("name")
	assert.Equal(t, res, "test-disk")
	assert.True(t, ok)
}

func TestFakeResourceData_getOkNotOk(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	res, ok := d.GetOk("incorrect")
	assert.Nil(t, res)
	assert.False(t, ok)
}

func TestFakeResourceData_getOk_missingList(t *testing.T) {
	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		getComputeDiskSchema(),
		values,
	)
	res, ok := d.GetOk("source_image_encryption_key")
	assert.Equal(t, []interface{}{}, res)
	assert.False(t, ok)
}

func TestAddressToTypes(t *testing.T) {
	cases := map[string]struct {
		address      []string
		block        *tfjson.SchemaBlock
		expectedType interface{}
	}{
		"bool": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Bool,
					},
				},
			},
			expectedType: cty.Bool,
		},
		"number": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Number,
					},
				},
			},
			expectedType: cty.Number,
		},
		"string": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.String,
					},
				},
			},
			expectedType: cty.String,
		},
		"missing primitive": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{},
			},
			expectedType: nil,
		},
		"address beyond primitive": {
			address: []string{"value", "subfield"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Bool,
					},
				},
			},
			expectedType: nil,
		},
		"list": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.List(cty.String),
					},
				},
			},
			expectedType: cty.List(cty.String),
		},
		"list index field returns number": {
			address: []string{"value", "#"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.List(cty.String),
					},
				},
			},
			expectedType: cty.Number,
		},
		"list at index returns element type": {
			address: []string{"value", "1"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.List(cty.String),
					},
				},
			},
			expectedType: cty.String,
		},
		"set": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Set(cty.String),
					},
				},
			},
			expectedType: cty.Set(cty.String),
		},
		"set index field returns number": {
			address: []string{"value", "#"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Set(cty.String),
					},
				},
			},
			expectedType: cty.Number,
		},
		"set at index returns element type": {
			address: []string{"value", "1"},
			block: &tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"value": &tfjson.SchemaAttribute{
						AttributeType: cty.Set(cty.String),
					},
				},
			},
			expectedType: cty.String,
		},
		"nested block: list": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.List(cty.DynamicPseudoType),
		},
		"nested block: list index returns number": {
			address: []string{"value", "#"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Number,
		},
		"nested block: list at index returns element type": {
			address: []string{"value", "1"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Map(cty.DynamicPseudoType),
		},
		"nested block: list element attribute type": {
			address: []string{"value", "1", "value"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Bool,
		},
		"nested block: set": {
			address: []string{"value"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeSet,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Set(cty.DynamicPseudoType),
		},
		"nested block: set index returns number": {
			address: []string{"value", "#"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeSet,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Number,
		},
		"nested block: set at index returns element type": {
			address: []string{"value", "1"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeSet,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Map(cty.DynamicPseudoType),
		},
		"nested block: set element attribute type": {
			address: []string{"value", "1", "value"},
			block: &tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"value": &tfjson.SchemaBlockType{
						NestingMode: tfjson.SchemaNestingModeSet,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"value": &tfjson.SchemaAttribute{
									AttributeType: cty.Bool,
								},
							},
						},
					},
				},
			},
			expectedType: cty.Bool,
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			valueType := addressToType(tc.address, tc.block)
			if tc.expectedType == nil {
				assert.Nil(t, valueType)
			} else {
				require.NotNil(t, valueType)
				assert.Equal(t, tc.expectedType, *valueType)
			}
		})
	}
}
