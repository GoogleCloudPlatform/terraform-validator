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
	"sort"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

func conversionFuncs() map[string]convertFunc {
	return map[string]convertFunc{
		// TODO: Use a generated mapping once it lands in the conversion library.
		"google_compute_disk":          converter.GetComputeDiskCaiObject,
		"google_compute_instance":      converter.GetComputeInstanceCaiObject,
		"google_storage_bucket":        converter.GetStorageBucketCaiObject,
		"google_sql_database_instance": converter.GetSQLDatabaseInstanceCaiObject,

		"google_organization_iam_policy":  converter.GetOrganizationIamPolicyCaiObject,
		"google_organization_iam_binding": converter.GetOrganizationIamBindingCaiObject,
		"google_organization_iam_member":  converter.GetOrganizationIamMemberCaiObject,
		"google_folder_iam_policy":        converter.GetFolderIamPolicyCaiObject,
		"google_folder_iam_binding":       converter.GetFolderIamBindingCaiObject,
		"google_folder_iam_member":        converter.GetFolderIamMemberCaiObject,
		"google_project_iam_policy":       converter.GetProjectIamPolicyCaiObject,
		"google_project_iam_binding":      converter.GetProjectIamBindingCaiObject,
		"google_project_iam_member":       converter.GetProjectIamMemberCaiObject,
	}
}

func mergeFuncs() map[string]mergeFunc {
	return map[string]mergeFunc{
		// TODO: Use a generated mapping once it lands in the conversion library.
		"google_organization_iam_policy":  converter.MergeOrganizationIamPolicy,
		"google_organization_iam_binding": converter.MergeOrganizationIamBinding,
		"google_organization_iam_member":  converter.MergeOrganizationIamMember,
		"google_folder_iam_policy":        converter.MergeFolderIamPolicy,
		"google_folder_iam_binding":       converter.MergeFolderIamBinding,
		"google_folder_iam_member":        converter.MergeFolderIamMember,
		"google_project_iam_policy":       converter.MergeProjectIamPolicy,
		"google_project_iam_binding":      converter.MergeProjectIamBinding,
		"google_project_iam_member":       converter.MergeProjectIamMember,
	}
}

// SupportedResources returns a sorted list of terraform resource names.
func SupportedTerraformResources() []string {
	fns := conversionFuncs()
	list := make([]string, 0, len(fns))
	for k := range fns {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
