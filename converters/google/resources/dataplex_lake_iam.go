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
const DataplexLakeIAMAssetType string = "dataplex.googleapis.com/Lake"

func resourceConverterDataplexLakeIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         DataplexLakeIAMAssetType,
		Convert:           GetDataplexLakeIamPolicyCaiObject,
		MergeCreateUpdate: MergeDataplexLakeIamPolicy,
	}
}

func resourceConverterDataplexLakeIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         DataplexLakeIAMAssetType,
		Convert:           GetDataplexLakeIamBindingCaiObject,
		FetchFullResource: FetchDataplexLakeIamPolicy,
		MergeCreateUpdate: MergeDataplexLakeIamBinding,
		MergeDelete:       MergeDataplexLakeIamBindingDelete,
	}
}

func resourceConverterDataplexLakeIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         DataplexLakeIAMAssetType,
		Convert:           GetDataplexLakeIamMemberCaiObject,
		FetchFullResource: FetchDataplexLakeIamPolicy,
		MergeCreateUpdate: MergeDataplexLakeIamMember,
		MergeDelete:       MergeDataplexLakeIamMemberDelete,
	}
}

func GetDataplexLakeIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newDataplexLakeIamAsset(d, config, expandIamPolicyBindings)
}

func GetDataplexLakeIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newDataplexLakeIamAsset(d, config, expandIamRoleBindings)
}

func GetDataplexLakeIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newDataplexLakeIamAsset(d, config, expandIamMemberBindings)
}

func MergeDataplexLakeIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeDataplexLakeIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeDataplexLakeIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeDataplexLakeIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeDataplexLakeIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newDataplexLakeIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//dataplex.googleapis.com/projects/{{project}}/locations/{{location}}/lakes/{{lake}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: DataplexLakeIAMAssetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchDataplexLakeIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("location"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}
	if _, ok := d.GetOk("lake"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		DataplexLakeIamUpdaterProducer,
		d,
		config,
		"//dataplex.googleapis.com/projects/{{project}}/locations/{{location}}/lakes/{{lake}}",
		DataplexLakeIAMAssetType,
	)
}
