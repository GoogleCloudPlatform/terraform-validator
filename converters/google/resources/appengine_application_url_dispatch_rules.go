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

const AppEngineApplicationUrlDispatchRulesAssetType string = "appengine.googleapis.com/ApplicationUrlDispatchRules"

func resourceConverterAppEngineApplicationUrlDispatchRules() ResourceConverter {
	return ResourceConverter{
		AssetType: AppEngineApplicationUrlDispatchRulesAssetType,
		Convert:   GetAppEngineApplicationUrlDispatchRulesCaiObject,
	}
}

func GetAppEngineApplicationUrlDispatchRulesCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//appengine.googleapis.com/apps/{{project}}/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetAppEngineApplicationUrlDispatchRulesApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: AppEngineApplicationUrlDispatchRulesAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/appengine/v1/rest",
				DiscoveryName:        "ApplicationUrlDispatchRules",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetAppEngineApplicationUrlDispatchRulesApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	dispatchRulesProp, err := expandAppEngineApplicationUrlDispatchRulesDispatchRules(d.Get("dispatch_rules"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("dispatch_rules"); !isEmptyValue(reflect.ValueOf(dispatchRulesProp)) && (ok || !reflect.DeepEqual(v, dispatchRulesProp)) {
		obj["dispatchRules"] = dispatchRulesProp
	}

	return obj, nil
}

func expandAppEngineApplicationUrlDispatchRulesDispatchRules(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDomain, err := expandAppEngineApplicationUrlDispatchRulesDispatchRulesDomain(original["domain"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDomain); val.IsValid() && !isEmptyValue(val) {
			transformed["domain"] = transformedDomain
		}

		transformedPath, err := expandAppEngineApplicationUrlDispatchRulesDispatchRulesPath(original["path"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPath); val.IsValid() && !isEmptyValue(val) {
			transformed["path"] = transformedPath
		}

		transformedService, err := expandAppEngineApplicationUrlDispatchRulesDispatchRulesService(original["service"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedService); val.IsValid() && !isEmptyValue(val) {
			transformed["service"] = transformedService
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandAppEngineApplicationUrlDispatchRulesDispatchRulesDomain(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationUrlDispatchRulesDispatchRulesPath(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationUrlDispatchRulesDispatchRulesService(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
