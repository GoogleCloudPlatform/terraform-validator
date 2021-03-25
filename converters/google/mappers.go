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

type convertFunc func(d converter.TerraformResourceData, config *converter.Config) ([]converter.Asset, error)

// fetchFunc allows initial data for a resource to be fetched from the API and merged
// with the planned changes. This is useful for resources that are only partially managed
// by Terraform, like IAM policies managed with member/binding resources.
type fetchFunc func(d converter.TerraformResourceData, config *converter.Config) (converter.Asset, error)

// mergeFunc combines multiple terraform resources into a single CAI asset.
// The incoming asset will either be an asset that was created/updated or deleted.
type mergeFunc func(existing, incoming converter.Asset) converter.Asset

// mapper pairs related conversion/merging functions.
type mapper struct {
	convert           convertFunc // required
	fetch             fetchFunc   // optional
	mergeCreateUpdate mergeFunc   // optional
	mergeDelete       mergeFunc   // optional
}

// mappers maps terraform resource types (i.e. `google_project`) into
// a slice of mapperFuncs.
//
// Modelling of relationships:
// terraform resources to CAI assets as []mapperFuncs:
// 1:1 = [mapper{convert: convertAbc}]                  (len=1)
// 1:N = [mapper{convert: convertAbc}, ...]             (len=N)
// N:1 = [mapper{convert: convertAbc, merge: mergeAbc}] (len=1)
func mappers() map[string][]mapper {
	return map[string][]mapper{
		// TODO: Use a generated mapping once it lands in the conversion library.
		"google_compute_firewall":               {{convert: converter.GetComputeFirewallCaiObject}},
		"google_compute_disk":                   {{convert: converter.GetComputeDiskCaiObject}},
		"google_compute_forwarding_rule":        {{convert: converter.GetComputeForwardingRuleCaiObject}},
		"google_compute_global_forwarding_rule": {{convert: converter.GetComputeGlobalForwardingRuleCaiObject}},
		"google_compute_instance":               {{convert: converter.GetComputeInstanceCaiObject}},
		"google_storage_bucket":                 {{convert: converter.GetStorageBucketCaiObject}},
		"google_sql_database_instance":          {{convert: converter.GetSQLDatabaseInstanceCaiObject}},
		"google_container_cluster":              {{convert: converter.GetContainerClusterCaiObject}},
		"google_container_node_pool":            {{convert: converter.GetContainerNodePoolCaiObject}},
		"google_bigquery_dataset":               {{convert: converter.GetBigQueryDatasetCaiObject}},
		"google_spanner_instance":               {{convert: converter.GetSpannerInstanceCaiObject}},
		"google_project_service":                {{convert: converter.GetServiceUsageCaiObject}},
		"google_pubsub_subscription":            {{convert: converter.GetPubsubSubscriptionCaiObject}},
		"google_pubsub_topic":                   {{convert: converter.GetPubsubTopicCaiObject}},
		"google_kms_crypto_key":                 {{convert: converter.GetKMSCryptoKeyCaiObject}},
		"google_kms_key_ring":                   {{convert: converter.GetKMSKeyRingCaiObject}},
		"google_filestore_instance":             {{convert: converter.GetFilestoreInstanceCaiObject}},

		// Terraform resources of type "google_project" have a 1:N relationship with CAI assets.
		"google_project": {
			{
				convert:           converter.GetProjectCaiObject,
				mergeCreateUpdate: mergeProject,
			},
			{convert: converter.GetProjectBillingInfoCaiObject},
		},

		// Terraform IAM policy resources have a N:1 relationship with CAI assets.
		"google_organization_iam_policy": {
			{
				convert:           converter.GetOrganizationIamPolicyCaiObject,
				mergeCreateUpdate: converter.MergeOrganizationIamPolicy,
			},
		},
		"google_project_organization_policy": {
			{
				convert:           converter.GetProjectOrgPolicyCaiObject,
				mergeCreateUpdate: converter.MergeProjectOrgPolicy,
			},
		},
		"google_organization_iam_binding": {
			{
				convert:           converter.GetOrganizationIamBindingCaiObject,
				mergeCreateUpdate: converter.MergeOrganizationIamBinding,
			},
		},
		"google_organization_iam_member": {
			{
				convert:           converter.GetOrganizationIamMemberCaiObject,
				mergeCreateUpdate: converter.MergeOrganizationIamMember,
			},
		},
		"google_folder_iam_policy": {
			{
				convert:           converter.GetFolderIamPolicyCaiObject,
				mergeCreateUpdate: converter.MergeFolderIamPolicy,
			},
		},
		"google_folder_iam_binding": {
			{
				convert:           converter.GetFolderIamBindingCaiObject,
				mergeCreateUpdate: converter.MergeFolderIamBinding,
			},
		},
		"google_folder_iam_member": {
			{
				convert:           converter.GetFolderIamMemberCaiObject,
				mergeCreateUpdate: converter.MergeFolderIamMember,
			},
		},
		"google_project_iam_policy": {
			{
				convert:           converter.GetProjectIamPolicyCaiObject,
				mergeCreateUpdate: converter.MergeProjectIamPolicy,
			},
		},
		"google_project_iam_binding": {
			{
				convert:           converter.GetProjectIamBindingCaiObject,
				mergeCreateUpdate: converter.MergeProjectIamBinding,
				mergeDelete:       converter.MergeProjectIamBindingDelete,
				fetch:             converter.FetchProjectIamPolicy,
			},
		},
		"google_project_iam_member": {
			{
				convert:           converter.GetProjectIamMemberCaiObject,
				mergeCreateUpdate: converter.MergeProjectIamMember,
				mergeDelete:       converter.MergeProjectIamMemberDelete,
				fetch:             converter.FetchProjectIamPolicy,
			},
		},
		"google_storage_bucket_iam_policy": {
			{
				convert:           converter.GetBucketIamPolicyCaiObject,
				mergeCreateUpdate: converter.MergeBucketIamPolicy,
			},
		},
		"google_storage_bucket_iam_binding": {
			{
				convert:           converter.GetBucketIamBindingCaiObject,
				mergeCreateUpdate: converter.MergeBucketIamBinding,
				mergeDelete:       converter.MergeBucketIamBindingDelete,
				fetch:             converter.FetchBucketIamPolicy,
			},
		},
		"google_storage_bucket_iam_member": {
			{
				convert:           converter.GetBucketIamMemberCaiObject,
				mergeCreateUpdate: converter.MergeBucketIamMember,
				mergeDelete:       converter.MergeBucketIamMemberDelete,
				fetch:             converter.FetchBucketIamPolicy,
			},
		},
	}
}

// SupportedResources returns a sorted list of terraform resource names.
func SupportedTerraformResources() []string {
	fns := mappers()
	list := make([]string, 0, len(fns))
	for k := range fns {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
