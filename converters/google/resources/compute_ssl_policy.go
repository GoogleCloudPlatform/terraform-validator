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
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sslPolicyCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	profile := diff.Get("profile")
	customFeaturesCount := diff.Get("custom_features.#")

	// Validate that policy configs aren't incompatible during all phases
	// CUSTOM profile demands non-zero custom_features, and other profiles (i.e., not CUSTOM) demand zero custom_features
	if diff.HasChange("profile") || diff.HasChange("custom_features") {
		if profile.(string) == "CUSTOM" {
			if customFeaturesCount.(int) == 0 {
				return fmt.Errorf("Error in SSL Policy %s: the profile is set to %s but no custom_features are set.", diff.Get("name"), profile.(string))
			}
		} else {
			if customFeaturesCount != 0 {
				return fmt.Errorf("Error in SSL Policy %s: the profile is set to %s but using custom_features requires the profile to be CUSTOM.", diff.Get("name"), profile.(string))
			}
		}
		return nil
	}
	return nil
}

const ComputeSslPolicyAssetType string = "compute.googleapis.com/SslPolicy"

func resourceConverterComputeSslPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType: ComputeSslPolicyAssetType,
		Convert:   GetComputeSslPolicyCaiObject,
	}
}

func GetComputeSslPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//compute.googleapis.com/projects/{{project}}/global/sslPolicies/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetComputeSslPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ComputeSslPolicyAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "SslPolicy",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetComputeSslPolicyApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeSslPolicyDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeSslPolicyName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	profileProp, err := expandComputeSslPolicyProfile(d.Get("profile"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("profile"); !isEmptyValue(reflect.ValueOf(profileProp)) && (ok || !reflect.DeepEqual(v, profileProp)) {
		obj["profile"] = profileProp
	}
	minTlsVersionProp, err := expandComputeSslPolicyMinTlsVersion(d.Get("min_tls_version"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("min_tls_version"); !isEmptyValue(reflect.ValueOf(minTlsVersionProp)) && (ok || !reflect.DeepEqual(v, minTlsVersionProp)) {
		obj["minTlsVersion"] = minTlsVersionProp
	}
	customFeaturesProp, err := expandComputeSslPolicyCustomFeatures(d.Get("custom_features"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("custom_features"); !isEmptyValue(reflect.ValueOf(customFeaturesProp)) && (ok || !reflect.DeepEqual(v, customFeaturesProp)) {
		obj["customFeatures"] = customFeaturesProp
	}

	return obj, nil
}

func expandComputeSslPolicyDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSslPolicyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSslPolicyProfile(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSslPolicyMinTlsVersion(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSslPolicyCustomFeatures(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}
