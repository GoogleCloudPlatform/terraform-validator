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

const workforcePoolIdRegexp = `^[a-z][a-z0-9-]{4,61}[a-z0-9]$`

func validateWorkforcePoolId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "gcp-") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with \"gcp-\". "+
				"The prefix `gcp-` is reserved for use by Google, and may not be specified.", k, value))
	}

	if !regexp.MustCompile(workforcePoolIdRegexp).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must contain only lowercase letters [a-z], digits [0-9], and hyphens "+
				"[-]. The WorkforcePool ID must be between 6 and 63 characters, begin "+
				"with a letter, and cannot have a trailing hyphen.", k, value))
	}

	return
}

const IAMWorkforcePoolWorkforcePoolAssetType string = "iam.googleapis.com/WorkforcePool"

func resourceConverterIAMWorkforcePoolWorkforcePool() ResourceConverter {
	return ResourceConverter{
		AssetType: IAMWorkforcePoolWorkforcePoolAssetType,
		Convert:   GetIAMWorkforcePoolWorkforcePoolCaiObject,
	}
}

func GetIAMWorkforcePoolWorkforcePoolCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/locations/{{location}}/workforcePools/{{workforce_pool_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetIAMWorkforcePoolWorkforcePoolApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: IAMWorkforcePoolWorkforcePoolAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/iam/v1/rest",
				DiscoveryName:        "WorkforcePool",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetIAMWorkforcePoolWorkforcePoolApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	parentProp, err := expandIAMWorkforcePoolWorkforcePoolParent(d.Get("parent"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("parent"); !isEmptyValue(reflect.ValueOf(parentProp)) && (ok || !reflect.DeepEqual(v, parentProp)) {
		obj["parent"] = parentProp
	}
	displayNameProp, err := expandIAMWorkforcePoolWorkforcePoolDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMWorkforcePoolWorkforcePoolDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	disabledProp, err := expandIAMWorkforcePoolWorkforcePoolDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	sessionDurationProp, err := expandIAMWorkforcePoolWorkforcePoolSessionDuration(d.Get("session_duration"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("session_duration"); !isEmptyValue(reflect.ValueOf(sessionDurationProp)) && (ok || !reflect.DeepEqual(v, sessionDurationProp)) {
		obj["sessionDuration"] = sessionDurationProp
	}

	return obj, nil
}

func expandIAMWorkforcePoolWorkforcePoolParent(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolWorkforcePoolDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolWorkforcePoolDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolWorkforcePoolDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolWorkforcePoolSessionDuration(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
