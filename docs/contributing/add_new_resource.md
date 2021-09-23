# Adding support for a new resource

## terraform-validator vs config-validator

At its core, terraform-validator is a thin layer on top of [config-validator](https://github.com/GoogleCloudPlatform/config-validator), a shared library that takes in a [policy library](https://github.com/GoogleCloudPlatform/policy-library) and a set of [CAI assets](https://cloud.google.com/asset-inventory/docs/overview) and reports back any violations of the specified policies.

terraform-validator consumes a Terraform plan and uses it to build CAI Assets, which then get run through config-validator. These built Assets only exist locally, in memory.

### Adding a new constraint template

If an existing [bundle](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/index.md#policy-bundles) (for example, [CIS v1.1](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/bundles/cis-v1.1.md)) doesn't support a check you need, please consider contributing a [new constraint template](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/constraint_template_authoring.md) to the policy-library repository.

### Getting a terraform resource name from a GCP resource name

The first step in determining if a GCP resource is supported is to figure out the name of the corresponding Terraform resource. You can often do this by searching for the GCP resource name in the [Terraform google provider documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

## How to add support for a new resource

A resource is "supported" by terraform-validator if it has an entry in [mappers.go](https://github.com/GoogleCloudPlatform/terraform-google-conversion/blob/master/google/mappers.go). For example, you could search mappers.go for [`google_compute_disk`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_disk) to see if that resource is supported.

Adding support for a resource has four steps:

1. Make changes to [Magic Modules](https://github.com/GoogleCloudPlatform/magic-modules) to add any necessary code to terraform-google-conversion.
2. Add tests for the new resource to terraform-validator, and run them against a [locally-generated](https://github.com/GoogleCloudPlatform/magic-modules/#generating-terraform-google-conversion) copy of terraform-google-conversion.
3. Make PRs for Magic Modules & terraform-validator with your changes. The reviewer will double-check that your code works and then merge the Magic Modules PR.
4. Once the Magic Modules PR is merged, it will automatically update terraform-google-conversion with your changes. Update your terraform-validator PR to use the new version of terraform-google.conversion.

Each of these is discussed in more detail below.

**Note**: terraform-validator can only support resources that are supported by the GA terraform provider, not beta resources.

### 1. Magic Modules

Magic Modules uses a shared code base to generate terraform-google-conversion and the [google](https://github.com/hashicorp/terraform-provider-google) and [google-beta](https://github.com/hashicorp/terraform-provider-google-beta) Terraform providers.
Most Terraform resources are represented as [yaml files which are grouped by product](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products).
Each product has an `api.yaml` file (which defines the basic API schema) and a `terraform.yaml` file (which defines any terraform-specific overrides.)
A `terraform.yaml` file can specify `exclude_validator: true` on a resource to skip terraform-google-conversion autogeneration, or `exclude_resource: true` to skip autogeneration for both terraform-google-conversion and the providers.

Auto-generating terraform-google-conversion code based on yaml files is strongly preferred; however, if that is not possible, you can instead place a handwritten file in the [`magic-modules/mmv1/third_party/validator` folder](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/third_party/validator).
Most resources will only need a conversion func, which should look something like:

```golang
func GetWhateverResourceCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	// get the correct format from https://cloud.google.com/asset-inventory/docs/supported-asset-types
	// The asset path (name) will substitute in variables from the Terraform resource
	name, err := assetName(d, config, "//whatever.googleapis.com/projects/{{project}}/whatevers/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetWhateverResourceApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			// The type also comes from https://cloud.google.com/asset-inventory/docs/supported-asset-types
			Type: "whatever.googleapis.com/Whatever",
			Resource: &AssetResource{
				Version:              "v1",  // or whatever the correct version is
				DiscoveryDocumentURI: "https://www.googleapis.com/path/to/rest/api/docs",
				DiscoveryName:        "Whatever",  // The term used to refer to this resource by the official documentation
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetWhateverResourceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	// copy values from the terraform resource to obj
	// return any errors encountered
	// ...

	return obj, nil
}

```

For handwritten conversion code, you will also need to add an entry to [`mappers.go.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/templates/validator/mappers/mappers.go.erb), which is used to generate the mappers.go file in terraform-google-conversion. Each entry in `mappers.go.erb` has the following keys:

- `convert`: Required. This function does basic conversion of a Terraform resource to a CAI Asset, including converting nested structures and specifying what the [CAI Asset Type](https://cloud.google.com/asset-inventory/docs/supported-asset-types) is.
- `fetch`, `mergeCreateUpdate`, `mergeDelete`: Optional. Some assets, like [IAM Members and Bindings](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam), have to be merged with remote data prior to validation in order to properly check whether policies are being followed. If you're not sure whether you need this, you probably don't.

To generate terraform-google-conversion code locally, run the following from the root of the `magic-modules` repository:

```
make validator OUTPUT_PATH="/path/to/your/terraform-google-conversion"
```

You can then run `make test` inside your terraform-google-conversion repository to make sure those tests pass.

### 2. Terraform Validator tests

Now that you have a local copy of terraform-google-conversion that has been [generated from Magic Modules](https://github.com/GoogleCloudPlatform/magic-modules/#generating-terraform-google-conversion), you need to make Terraform Validator use it for local testing. You can do this with a [`replace` directive](https://golang.org/ref/mod#go-mod-file-replace):

```
replace github.com/GoogleCloudPlatform/terraform-google-conversion => /path/to/your/terraform-google-conversion
```

You can now build the binary (with `make build`) and test it. One way to do this would be to create a test project following the instructions in the [policy library user guide](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/user_guide.md#for-local-development-environments) (but using the binary you just built.) It's easiest to use a [GCPAlwaysViolatesConstraintV1](https://github.com/GoogleCloudPlatform/terraform-validator/blob/master/testdata/sample_policies/always_violate/policies/constraints/always_violates.yaml) constraint for testing new resources; this is what the tests do. `terraform-validator convert tfplan.json` can show you what terraform-validator thinks the converted Asset looks like.

Be sure to add test cases to [test/cli_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/cli_test.go#L52) and [test/read_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/read_test.go#L24). The test names refer to files in [testdata/templates](https://github.com/GoogleCloudPlatform/terraform-validator/tree/master/testdata/templates). You will generally need to add the following files:
   - A .tf file.
   - A .tfplan.json file.
   - A .json file (representing the output of `terraform-validator convert`)

See [Getting started](./getting_started.md) for details on running tests.

Try to get your tests passing locally before proceeding. (But you can also go ahead and open PRs if you're running into issues you can't figure out how to resolve.)

### 3. Make PRs

Now that you have your code working locally, open PRs for [Magic Modules](https://github.com/GoogleCloudPlatform/magic-modules) and terraform-validator. The reviewer will make sure your code works as expected.

For the Magic Modules PR, the most important check is `terraform-google-conversion-test` - as long as that's passing, you're probably fine. If it is failing, go back to step 1 and try running `make test` for terraform-google-conversion to reproduce & fix the failure.

For terraform-validator, the CI tests will not pass at this point, because the terraform-google-conversion dependency has not yet been updated. As long as the tests are passing locally for you, it should be fine.

### 4. Update terraform-google-conversion dependency

After the Magic Modules PR is merged, and the terraform-google-conversion repository contains your changes, update the terraform-google-conversion dependency in your terraform-validator PR. This command will make the necessary changes:

```bash
go get github.com/GoogleCloudPlatform/terraform-google-conversion
```

If the CI tests are still failing after you make this change, double-check that you're able to run `make test` locally inside the terraform-validator repository using the updated (not replaced) dependency.
