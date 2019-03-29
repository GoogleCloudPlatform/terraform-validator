package tfgcv

import (
	"testing"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func TestProtoViaJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    interface{}
		expected *validator.Asset
	}{
		{
			name:     "Nil",
			input:    nil,
			expected: &validator.Asset{},
		},
		{
			name:     "EmptyAssetMap",
			input:    map[string]interface{}{},
			expected: &validator.Asset{},
		},
		{
			name: "EmptyResourceMap",
			input: map[string]interface{}{
				"resource": map[string]interface{}{},
			},
			expected: &validator.Asset{
				Resource: &asset.Resource{},
			},
		},
		{
			name: "EmptyResourceDataMap",
			input: map[string]interface{}{
				"resource": map[string]interface{}{
					"data": map[string]interface{}{},
				},
			},
			expected: &validator.Asset{
				Resource: &asset.Resource{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{},
					},
				},
			},
		},
		{
			name: "ResourceMapEmptyValue",
			input: map[string]interface{}{
				"resource": map[string]interface{}{
					"data": map[string]interface{}{
						"abc": nil,
					},
				},
			},
			expected: &validator.Asset{
				Resource: &asset.Resource{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"abc": {
								Kind: &structpb.Value_NullValue{},
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			msg := &validator.Asset{}
			if err := protoViaJSON(c.input, msg); err != nil {
				t.Error(err)
			}
			require.EqualValues(t, c.expected, msg)
		})
	}
}
