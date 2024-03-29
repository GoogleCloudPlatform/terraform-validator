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

const ApigeeEnvgroupAssetType string = "apigee.googleapis.com/Envgroup"

func resourceConverterApigeeEnvgroup() ResourceConverter {
	return ResourceConverter{
		AssetType: ApigeeEnvgroupAssetType,
		Convert:   GetApigeeEnvgroupCaiObject,
	}
}

func GetApigeeEnvgroupCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//apigee.googleapis.com/{{org_id}}/envgroups/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetApigeeEnvgroupApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ApigeeEnvgroupAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/apigee/v1/rest",
				DiscoveryName:        "Envgroup",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetApigeeEnvgroupApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandApigeeEnvgroupName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	hostnamesProp, err := expandApigeeEnvgroupHostnames(d.Get("hostnames"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("hostnames"); !isEmptyValue(reflect.ValueOf(hostnamesProp)) && (ok || !reflect.DeepEqual(v, hostnamesProp)) {
		obj["hostnames"] = hostnamesProp
	}

	return obj, nil
}

func expandApigeeEnvgroupName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvgroupHostnames(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
