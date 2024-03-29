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

import "reflect"

const ApigeeEnvReferencesAssetType string = "apigee.googleapis.com/EnvReferences"

func resourceConverterApigeeEnvReferences() ResourceConverter {
	return ResourceConverter{
		AssetType: ApigeeEnvReferencesAssetType,
		Convert:   GetApigeeEnvReferencesCaiObject,
	}
}

func GetApigeeEnvReferencesCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//apigee.googleapis.com/{{env_id}}/references/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetApigeeEnvReferencesApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ApigeeEnvReferencesAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/apigee/v1/rest",
				DiscoveryName:        "EnvReferences",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetApigeeEnvReferencesApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandApigeeEnvReferencesName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandApigeeEnvReferencesDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	resourceTypeProp, err := expandApigeeEnvReferencesResourceType(d.Get("resource_type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("resource_type"); !isEmptyValue(reflect.ValueOf(resourceTypeProp)) && (ok || !reflect.DeepEqual(v, resourceTypeProp)) {
		obj["resourceType"] = resourceTypeProp
	}
	refersProp, err := expandApigeeEnvReferencesRefers(d.Get("refers"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("refers"); !isEmptyValue(reflect.ValueOf(refersProp)) && (ok || !reflect.DeepEqual(v, refersProp)) {
		obj["refers"] = refersProp
	}

	return obj, nil
}

func expandApigeeEnvReferencesName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvReferencesDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvReferencesResourceType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvReferencesRefers(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
