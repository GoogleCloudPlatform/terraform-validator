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
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func Test_jsonResourceFieldReader_ReadField(t *testing.T) {
	type fields struct {
		Name  string
		Value interface{}
		Type  schema.ValueType
	}
	type args struct {
		address string
	}
	type want struct {
		Value interface{}
		Found bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			"string",
			fields{"name", "string_value", schema.TypeString},
			args{"name"},
			want{"string_value", true},
			false,
		},
		{
			"bool",
			fields{"name", "true", schema.TypeBool},
			args{"name"},
			want{"true", true},
			false,
		},
		{
			"set",
			fields{"name", []interface{}{"1", "2"}, schema.TypeSet},
			args{"name"},
			want{[]interface{}{"2", "1"}, true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := jsonResourceFieldReader{
				Source: createResource(tt.fields.Name, tt.fields.Value),
				Schema: createSchemaMap(tt.args.address, tt.fields.Type),
			}
			address := []string{tt.args.address}
			got, err := r.ReadField(address)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonResourceFieldReader.ReadField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			switch got.Value.(type) {
			case *schema.Set:
				tmpVal := got.Value.(*schema.Set)
				got.Value = tmpVal.List()
			}
			if !reflect.DeepEqual(got.Value, tt.want.Value) {
				t.Errorf("jsonResourceFieldReader.ReadField() = %v, want %v", got.Value, tt.want.Value)
			}

		})
	}
}

func createSchemaMap(address string, valueType schema.ValueType) map[string]*schema.Schema {
	result := map[string]*schema.Schema{}
	result[address] = &schema.Schema{Type: valueType}
	return result
}

func createResource(address string, value interface{}) jsonResource {
	values := map[string]interface{}{}
	values[address] = value
	return jsonResource{Values: values}
}
