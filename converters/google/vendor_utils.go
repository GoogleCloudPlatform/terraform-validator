// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"fmt"
	"log"
	"os"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

// NOTE: These functions were pulled from github.com/terraform-providers/terraform-provider-google. They can go away when the functionality they are providing is implemented in the future github.com/GoogleCloudPlatform/terraform-converters package.

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d converter.TerraformResourceData, config *converter.Config, cai converter.Asset) (string, error) {
	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Project",
		"cloudbilling.googleapis.com/ProjectBillingInfo":
		res, ok := d.GetOk("number")
		if ok {
			return res.(string), nil
		}
		// Fall back to project_id if number is not available.
		res, ok = d.GetOk("project_id")
		if ok {
			return res.(string), nil
		} else {
			log.Printf("[WARN] Failed to retrieve project_id for %s from resource", cai.Name)
		}
	case "storage.googleapis.com/Bucket":
		if cai.Resource != nil {
			res, ok := cai.Resource.Data["project"]
			if ok {
				return res.(string), nil
			}
		}
		log.Printf("[WARN] Failed to retrieve project_id for %s from cai resource", cai.Name)
	}

	return getProjectFromSchema("project", d, config)
}

func getProjectFromSchema(projectSchemaField string, d converter.TerraformResourceData, config *converter.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("required field '%s' is not set", projectSchemaField)
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
