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

const StorageHmacKeyAssetType string = "storage.googleapis.com/HmacKey"

func resourceConverterStorageHmacKey() ResourceConverter {
	return ResourceConverter{
		AssetType: StorageHmacKeyAssetType,
		Convert:   GetStorageHmacKeyCaiObject,
	}
}

func GetStorageHmacKeyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//storage.googleapis.com/projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetStorageHmacKeyApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: StorageHmacKeyAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/storage/v1/rest",
				DiscoveryName:        "HmacKey",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetStorageHmacKeyApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	serviceAccountEmailProp, err := expandStorageHmacKeyServiceAccountEmail(d.Get("service_account_email"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service_account_email"); !isEmptyValue(reflect.ValueOf(serviceAccountEmailProp)) && (ok || !reflect.DeepEqual(v, serviceAccountEmailProp)) {
		obj["serviceAccountEmail"] = serviceAccountEmailProp
	}
	stateProp, err := expandStorageHmacKeyState(d.Get("state"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("state"); !isEmptyValue(reflect.ValueOf(stateProp)) && (ok || !reflect.DeepEqual(v, stateProp)) {
		obj["state"] = stateProp
	}

	return obj, nil
}

func expandStorageHmacKeyServiceAccountEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageHmacKeyState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
