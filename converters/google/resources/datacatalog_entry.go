// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"encoding/json"
	"reflect"
)

const DataCatalogEntryAssetType string = "datacatalog.googleapis.com/Entry"

func resourceConverterDataCatalogEntry() ResourceConverter {
	return ResourceConverter{
		AssetType: DataCatalogEntryAssetType,
		Convert:   GetDataCatalogEntryCaiObject,
	}
}

func GetDataCatalogEntryCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//datacatalog.googleapis.com/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetDataCatalogEntryApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: DataCatalogEntryAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/datacatalog/v1/rest",
				DiscoveryName:        "Entry",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetDataCatalogEntryApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	linkedResourceProp, err := expandDataCatalogEntryLinkedResource(d.Get("linked_resource"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("linked_resource"); !isEmptyValue(reflect.ValueOf(linkedResourceProp)) && (ok || !reflect.DeepEqual(v, linkedResourceProp)) {
		obj["linkedResource"] = linkedResourceProp
	}
	displayNameProp, err := expandDataCatalogEntryDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandDataCatalogEntryDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	schemaProp, err := expandDataCatalogEntrySchema(d.Get("schema"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("schema"); !isEmptyValue(reflect.ValueOf(schemaProp)) && (ok || !reflect.DeepEqual(v, schemaProp)) {
		obj["schema"] = schemaProp
	}
	typeProp, err := expandDataCatalogEntryType(d.Get("type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("type"); !isEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	userSpecifiedTypeProp, err := expandDataCatalogEntryUserSpecifiedType(d.Get("user_specified_type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("user_specified_type"); !isEmptyValue(reflect.ValueOf(userSpecifiedTypeProp)) && (ok || !reflect.DeepEqual(v, userSpecifiedTypeProp)) {
		obj["userSpecifiedType"] = userSpecifiedTypeProp
	}
	userSpecifiedSystemProp, err := expandDataCatalogEntryUserSpecifiedSystem(d.Get("user_specified_system"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("user_specified_system"); !isEmptyValue(reflect.ValueOf(userSpecifiedSystemProp)) && (ok || !reflect.DeepEqual(v, userSpecifiedSystemProp)) {
		obj["userSpecifiedSystem"] = userSpecifiedSystemProp
	}
	gcsFilesetSpecProp, err := expandDataCatalogEntryGcsFilesetSpec(d.Get("gcs_fileset_spec"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("gcs_fileset_spec"); !isEmptyValue(reflect.ValueOf(gcsFilesetSpecProp)) && (ok || !reflect.DeepEqual(v, gcsFilesetSpecProp)) {
		obj["gcsFilesetSpec"] = gcsFilesetSpecProp
	}

	return obj, nil
}

func expandDataCatalogEntryLinkedResource(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntrySchema(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	b := []byte(v.(string))
	if len(b) == 0 {
		return nil, nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func expandDataCatalogEntryType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryUserSpecifiedType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryUserSpecifiedSystem(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryGcsFilesetSpec(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFilePatterns, err := expandDataCatalogEntryGcsFilesetSpecFilePatterns(original["file_patterns"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFilePatterns); val.IsValid() && !isEmptyValue(val) {
		transformed["filePatterns"] = transformedFilePatterns
	}

	transformedSampleGcsFileSpecs, err := expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecs(original["sample_gcs_file_specs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSampleGcsFileSpecs); val.IsValid() && !isEmptyValue(val) {
		transformed["sampleGcsFileSpecs"] = transformedSampleGcsFileSpecs
	}

	return transformed, nil
}

func expandDataCatalogEntryGcsFilesetSpecFilePatterns(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecs(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedFilePath, err := expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecsFilePath(original["file_path"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFilePath); val.IsValid() && !isEmptyValue(val) {
			transformed["filePath"] = transformedFilePath
		}

		transformedSizeBytes, err := expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecsSizeBytes(original["size_bytes"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSizeBytes); val.IsValid() && !isEmptyValue(val) {
			transformed["sizeBytes"] = transformedSizeBytes
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecsFilePath(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDataCatalogEntryGcsFilesetSpecSampleGcsFileSpecsSizeBytes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
