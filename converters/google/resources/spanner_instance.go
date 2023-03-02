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
	"log"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deleteSpannerBackups(d *schema.ResourceData, config *Config, res map[string]interface{}, userAgent string, billingProject string) error {
	var v interface{}
	var ok bool

	v, ok = res["backups"]
	if !ok || v == nil {
		return nil
	}

	// Iterate over the list and delete each backup.
	for _, itemRaw := range v.([]interface{}) {
		if itemRaw == nil {
			continue
		}
		item := itemRaw.(map[string]interface{})

		backupName := item["name"].(string)

		log.Printf("[DEBUG] Found backups for resource %q: %#v)", d.Id(), item)

		path := "{{SpannerBasePath}}" + backupName

		url, err := replaceVars(d, config, path)
		if err != nil {
			return err
		}

		_, err = SendRequest(config, "DELETE", billingProject, url, userAgent, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceSpannerInstanceVirtualUpdate(d *schema.ResourceData, resourceSchema map[string]*schema.Schema) bool {
	// force_destroy is the only virtual field
	if d.HasChange("force_destroy") {
		for field := range resourceSchema {
			if field == "force_destroy" {
				continue
			}
			if d.HasChange(field) {
				return false
			}
		}
		return true
	}
	return false
}

const SpannerInstanceAssetType string = "spanner.googleapis.com/Instance"

func resourceConverterSpannerInstance() ResourceConverter {
	return ResourceConverter{
		AssetType: SpannerInstanceAssetType,
		Convert:   GetSpannerInstanceCaiObject,
	}
}

func GetSpannerInstanceCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetSpannerInstanceApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: SpannerInstanceAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/spanner/v1/rest",
				DiscoveryName:        "Instance",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetSpannerInstanceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandSpannerInstanceName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	configProp, err := expandSpannerInstanceConfig(d.Get("config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("config"); !isEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	displayNameProp, err := expandSpannerInstanceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	nodeCountProp, err := expandSpannerInstanceNumNodes(d.Get("num_nodes"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("num_nodes"); !isEmptyValue(reflect.ValueOf(nodeCountProp)) && (ok || !reflect.DeepEqual(v, nodeCountProp)) {
		obj["nodeCount"] = nodeCountProp
	}
	processingUnitsProp, err := expandSpannerInstanceProcessingUnits(d.Get("processing_units"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("processing_units"); !isEmptyValue(reflect.ValueOf(processingUnitsProp)) && (ok || !reflect.DeepEqual(v, processingUnitsProp)) {
		obj["processingUnits"] = processingUnitsProp
	}
	labelsProp, err := expandSpannerInstanceLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	return resourceSpannerInstanceEncoder(d, config, obj)
}

func resourceSpannerInstanceEncoder(d TerraformResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// Temp Logic to accommodate processing_units and num_nodes
	if obj["processingUnits"] == nil && obj["nodeCount"] == nil {
		obj["nodeCount"] = 1
	}
	newObj := make(map[string]interface{})
	newObj["instance"] = obj
	if obj["name"] == nil {
		if err := d.Set("name", resource.PrefixedUniqueId("tfgen-spanid-")[:30]); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
		newObj["instanceId"] = d.Get("name").(string)
	} else {
		newObj["instanceId"] = obj["name"]
	}
	delete(obj, "name")
	return newObj, nil
}

func expandSpannerInstanceName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	r := regexp.MustCompile("projects/(.+)/instanceConfigs/(.+)")
	if r.MatchString(v.(string)) {
		return v.(string), nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("projects/%s/instanceConfigs/%s", project, v.(string)), nil
}

func expandSpannerInstanceDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceNumNodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceProcessingUnits(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
