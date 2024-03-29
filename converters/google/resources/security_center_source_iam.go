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
const SecurityCenterSourceIAMAssetType string = "securitycenter.googleapis.com/Source"

func resourceConverterSecurityCenterSourceIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         SecurityCenterSourceIAMAssetType,
		Convert:           GetSecurityCenterSourceIamPolicyCaiObject,
		MergeCreateUpdate: MergeSecurityCenterSourceIamPolicy,
	}
}

func resourceConverterSecurityCenterSourceIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         SecurityCenterSourceIAMAssetType,
		Convert:           GetSecurityCenterSourceIamBindingCaiObject,
		FetchFullResource: FetchSecurityCenterSourceIamPolicy,
		MergeCreateUpdate: MergeSecurityCenterSourceIamBinding,
		MergeDelete:       MergeSecurityCenterSourceIamBindingDelete,
	}
}

func resourceConverterSecurityCenterSourceIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         SecurityCenterSourceIAMAssetType,
		Convert:           GetSecurityCenterSourceIamMemberCaiObject,
		FetchFullResource: FetchSecurityCenterSourceIamPolicy,
		MergeCreateUpdate: MergeSecurityCenterSourceIamMember,
		MergeDelete:       MergeSecurityCenterSourceIamMemberDelete,
	}
}

func GetSecurityCenterSourceIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSecurityCenterSourceIamAsset(d, config, expandIamPolicyBindings)
}

func GetSecurityCenterSourceIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSecurityCenterSourceIamAsset(d, config, expandIamRoleBindings)
}

func GetSecurityCenterSourceIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSecurityCenterSourceIamAsset(d, config, expandIamMemberBindings)
}

func MergeSecurityCenterSourceIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSecurityCenterSourceIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeSecurityCenterSourceIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeSecurityCenterSourceIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeSecurityCenterSourceIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newSecurityCenterSourceIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//securitycenter.googleapis.com/organizations/{{organization}}/sources/{{source}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: SecurityCenterSourceIAMAssetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSecurityCenterSourceIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("organization"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}
	if _, ok := d.GetOk("source"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		SecurityCenterSourceIamUpdaterProducer,
		d,
		config,
		"//securitycenter.googleapis.com/organizations/{{organization}}/sources/{{source}}",
		SecurityCenterSourceIAMAssetType,
	)
}
