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

import "fmt"

// Provide a separate asset type constant so we don't have to worry about name conflicts between IAM and non-IAM converter files
const GKEBackupBackupPlanIAMAssetType string = "gkebackup.googleapis.com/BackupPlan"

func resourceConverterGKEBackupBackupPlanIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         GKEBackupBackupPlanIAMAssetType,
		Convert:           GetGKEBackupBackupPlanIamPolicyCaiObject,
		MergeCreateUpdate: MergeGKEBackupBackupPlanIamPolicy,
	}
}

func resourceConverterGKEBackupBackupPlanIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         GKEBackupBackupPlanIAMAssetType,
		Convert:           GetGKEBackupBackupPlanIamBindingCaiObject,
		FetchFullResource: FetchGKEBackupBackupPlanIamPolicy,
		MergeCreateUpdate: MergeGKEBackupBackupPlanIamBinding,
		MergeDelete:       MergeGKEBackupBackupPlanIamBindingDelete,
	}
}

func resourceConverterGKEBackupBackupPlanIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         GKEBackupBackupPlanIAMAssetType,
		Convert:           GetGKEBackupBackupPlanIamMemberCaiObject,
		FetchFullResource: FetchGKEBackupBackupPlanIamPolicy,
		MergeCreateUpdate: MergeGKEBackupBackupPlanIamMember,
		MergeDelete:       MergeGKEBackupBackupPlanIamMemberDelete,
	}
}

func GetGKEBackupBackupPlanIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newGKEBackupBackupPlanIamAsset(d, config, expandIamPolicyBindings)
}

func GetGKEBackupBackupPlanIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newGKEBackupBackupPlanIamAsset(d, config, expandIamRoleBindings)
}

func GetGKEBackupBackupPlanIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newGKEBackupBackupPlanIamAsset(d, config, expandIamMemberBindings)
}

func MergeGKEBackupBackupPlanIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeGKEBackupBackupPlanIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeGKEBackupBackupPlanIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeGKEBackupBackupPlanIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeGKEBackupBackupPlanIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newGKEBackupBackupPlanIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//gkebackup.googleapis.com/projects/{{project}}/locations/{{location}}/backupPlans/{{name}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: GKEBackupBackupPlanIAMAssetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchGKEBackupBackupPlanIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("location"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}
	if _, ok := d.GetOk("name"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		GKEBackupBackupPlanIamUpdaterProducer,
		d,
		config,
		"//gkebackup.googleapis.com/projects/{{project}}/locations/{{location}}/backupPlans/{{name}}",
		GKEBackupBackupPlanIAMAssetType,
	)
}
