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
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const workloadIdentityPoolIdRegexp = `^[0-9a-z-]+$`

func validateWorkloadIdentityPoolId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "gcp-") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with \"gcp-\"", k, value))
	}

	if !regexp.MustCompile(workloadIdentityPoolIdRegexp).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must contain only lowercase letters (a-z), numbers (0-9), or dashes (-)", k))
	}

	if len(value) < 4 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be smaller than 4 characters", k))
	}

	if len(value) > 32 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be greater than 32 characters", k))
	}

	return
}

const IAMBetaWorkloadIdentityPoolAssetType string = "iam.googleapis.com/WorkloadIdentityPool"

func resourceConverterIAMBetaWorkloadIdentityPool() ResourceConverter {
	return ResourceConverter{
		AssetType: IAMBetaWorkloadIdentityPoolAssetType,
		Convert:   GetIAMBetaWorkloadIdentityPoolCaiObject,
	}
}

func GetIAMBetaWorkloadIdentityPoolCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetIAMBetaWorkloadIdentityPoolApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: IAMBetaWorkloadIdentityPoolAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/iam/v1/rest",
				DiscoveryName:        "WorkloadIdentityPool",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetIAMBetaWorkloadIdentityPoolApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	displayNameProp, err := expandIAMBetaWorkloadIdentityPoolDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMBetaWorkloadIdentityPoolDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	disabledProp, err := expandIAMBetaWorkloadIdentityPoolDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	return obj, nil
}

func expandIAMBetaWorkloadIdentityPoolDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
