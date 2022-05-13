// Copyright 2022 Google LLC
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

package utils

import (
	"fmt"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"go.uber.org/zap"
)

// GetProjectFromResource reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func GetProjectFromResource(d resources.TerraformResourceData, config *resources.Config, cai resources.Asset, errorLogger *zap.Logger) (string, error) {

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
			errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from resource", cai.Name))
		}
	case "storage.googleapis.com/Bucket":
		if cai.Resource != nil {
			res, ok := cai.Resource.Data["project"]
			if ok {
				return res.(string), nil
			}
		}
		errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from cai resource", cai.Name))
	}

	return getProjectFromSchema("project", d, config)
}

func getProjectFromSchema(projectSchemaField string, d resources.TerraformResourceData, config *resources.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("required field '%s' is not set, you may use --project=my-project to provide a default project to resolve the issue", projectSchemaField)
}

// GetOrganizationFromResource reads org_id field from terraform data.
func GetOrganizationFromResource(tfData resources.TerraformResourceData) (string, bool) {
	orgID, ok := tfData.GetOk("org_id")
	if ok {
		return orgID.(string), ok
	}
	return "", false
}

// GetFolderFromResource reads folder_id or folder field from terraform data.
func GetFolderFromResource(tfData resources.TerraformResourceData) (string, bool) {
	folderID, ok := tfData.GetOk("folder_id")
	if ok {
		return folderID.(string), ok
	}
	folderID, ok = tfData.GetOk("folder")
	if ok {
		return folderID.(string), ok
	}
	return "", false
}
