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

const workloadIdentityPoolProviderIdRegexp = `^[0-9a-z-]+$`

func validateWorkloadIdentityPoolProviderId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "gcp-") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with \"gcp-\"", k, value))
	}

	if !regexp.MustCompile(workloadIdentityPoolProviderIdRegexp).MatchString(value) {
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

const IAMBetaWorkloadIdentityPoolProviderAssetType string = "iam.googleapis.com/WorkloadIdentityPoolProvider"

func resourceConverterIAMBetaWorkloadIdentityPoolProvider() ResourceConverter {
	return ResourceConverter{
		AssetType: IAMBetaWorkloadIdentityPoolProviderAssetType,
		Convert:   GetIAMBetaWorkloadIdentityPoolProviderCaiObject,
	}
}

func GetIAMBetaWorkloadIdentityPoolProviderCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetIAMBetaWorkloadIdentityPoolProviderApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: IAMBetaWorkloadIdentityPoolProviderAssetType,
			Resource: &AssetResource{
				Version:              "v",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/iam/v/rest",
				DiscoveryName:        "WorkloadIdentityPoolProvider",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetIAMBetaWorkloadIdentityPoolProviderApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	displayNameProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMBetaWorkloadIdentityPoolProviderDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	disabledProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	attributeMappingProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeMapping(d.Get("attribute_mapping"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("attribute_mapping"); !isEmptyValue(reflect.ValueOf(attributeMappingProp)) && (ok || !reflect.DeepEqual(v, attributeMappingProp)) {
		obj["attributeMapping"] = attributeMappingProp
	}
	attributeConditionProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeCondition(d.Get("attribute_condition"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("attribute_condition"); !isEmptyValue(reflect.ValueOf(attributeConditionProp)) && (ok || !reflect.DeepEqual(v, attributeConditionProp)) {
		obj["attributeCondition"] = attributeConditionProp
	}
	awsProp, err := expandIAMBetaWorkloadIdentityPoolProviderAws(d.Get("aws"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("aws"); !isEmptyValue(reflect.ValueOf(awsProp)) && (ok || !reflect.DeepEqual(v, awsProp)) {
		obj["aws"] = awsProp
	}
	oidcProp, err := expandIAMBetaWorkloadIdentityPoolProviderOidc(d.Get("oidc"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("oidc"); !isEmptyValue(reflect.ValueOf(oidcProp)) && (ok || !reflect.DeepEqual(v, oidcProp)) {
		obj["oidc"] = oidcProp
	}

	return obj, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAttributeMapping(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAttributeCondition(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAws(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAccountId, err := expandIAMBetaWorkloadIdentityPoolProviderAwsAccountId(original["account_id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccountId); val.IsValid() && !isEmptyValue(val) {
		transformed["accountId"] = transformedAccountId
	}

	return transformed, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAwsAccountId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidc(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowedAudiences, err := expandIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(original["allowed_audiences"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowedAudiences); val.IsValid() && !isEmptyValue(val) {
		transformed["allowedAudiences"] = transformedAllowedAudiences
	}

	transformedIssuerUri, err := expandIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(original["issuer_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIssuerUri); val.IsValid() && !isEmptyValue(val) {
		transformed["issuerUri"] = transformedIssuerUri
	}

	return transformed, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
