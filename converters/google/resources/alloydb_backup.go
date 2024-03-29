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

const AlloydbBackupAssetType string = "alloydb.googleapis.com/Backup"

func resourceConverterAlloydbBackup() ResourceConverter {
	return ResourceConverter{
		AssetType: AlloydbBackupAssetType,
		Convert:   GetAlloydbBackupCaiObject,
	}
}

func GetAlloydbBackupCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//alloydb.googleapis.com/projects/{{project}}/locations/{{location}}/backups/{{backup_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetAlloydbBackupApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: AlloydbBackupAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/alloydb/v1/rest",
				DiscoveryName:        "Backup",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetAlloydbBackupApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	clusterNameProp, err := expandAlloydbBackupClusterName(d.Get("cluster_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("cluster_name"); !isEmptyValue(reflect.ValueOf(clusterNameProp)) && (ok || !reflect.DeepEqual(v, clusterNameProp)) {
		obj["clusterName"] = clusterNameProp
	}
	labelsProp, err := expandAlloydbBackupLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	descriptionProp, err := expandAlloydbBackupDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	return resourceAlloydbBackupEncoder(d, config, obj)
}

func resourceAlloydbBackupEncoder(d TerraformResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// The only other available type is AUTOMATED which cannot be set manually
	obj["type"] = "ON_DEMAND"
	return obj, nil
}

func expandAlloydbBackupClusterName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandAlloydbBackupLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandAlloydbBackupDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
