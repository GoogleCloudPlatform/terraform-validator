# Adding support for a new resource

## terraform-validator vs config-validator

At its core, terraform-validator is a thin layer on top of [config-validator](https://github.com/forseti-security/config-validator), a shared library that takes in a [policy library](https://github.com/forseti-security/policy-library) and a set of [CAI assets](https://cloud.google.com/asset-inventory/docs/overview) and reports back any violations of the specified policies.

terraform-validator consumes a Terraform plan and uses it to build CAI Assets, which then get run through config-validator. These built Assets only exist locally, in memory.

**Note**: Although policy-library is a repository inside of the forseti-security organization, Terraform Validator does _not_ require an active installation of Forseti. Terraform Validator is a self-contained binary. Policy libraries are a configuration mechanism shared by a number of tools via config-validator.

### Adding a new constraint template

If an existing [bundle](https://github.com/forseti-security/policy-library/blob/master/docs/index.md#policy-bundles) (for example, [CIS v1.1](https://github.com/forseti-security/policy-library/blob/master/docs/bundles/cis-v1.1.md)) doesn't support a check you need, please consider contributing a [new constraint template](https://github.com/forseti-security/policy-library/blob/master/docs/constraint_template_authoring.md) to the policy-library repository.

### Getting a terraform resource name from a GCP resource name

The first step in determining if a GCP resource is supported is to figure out the name of the corresponding Terraform resource. You can often do this by searching for the GCP resource name in the [Terraform google provider documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

## How to add support for a new resource

A resource is "supported" by terraform-validator if it has an entry in [mappers.go](../converters/google/mappers.go). For example, you could search mappers.go for [`google_compute_disk`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_disk) to see if that resource is supported. Each entry in mappers.go has the following keys:

- `convert`: Required. This function does basic conversion of a Terraform resource to a CAI Asset, including converting nested structures and specifying what the [CAI Asset Type](https://cloud.google.com/asset-inventory/docs/supported-asset-types) is.
- `fetch`, `mergeCreateUpdate`, `mergeDelete`: Optional. Some assets, like [IAM Members and Bindings](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam), have to be merged with remote data prior to validation in order to properly check whether policies are being followed. If you're not sure whether you need this, you probably don't.

The code referenced by the mappers comes from [terraform-google-conversion](https://github.com/GoogleCloudPlatform/terraform-google-conversion/tree/master/google), but that repository is generated using [Magic Modules](https://github.com/googleCloudPlatform/magic-modules/). Magic Modules uses a shared code base to generate terraform-google-conversion and the [google](https://github.com/hashicorp/terraform-provider-google) and [google-beta](https://github.com/hashicorp/terraform-provider-google-beta) Terraform providers.

So, adding support for a resource means:

1. Make a PR for Magic Modules to add the necessary code to terraform-google-conversion. Once your PR is merged, the code will be automatically copied into terraform-google-conversion.
2. Make a PR for terraform-validator that updates the version of terraform-google-conversion and adds a new mapper to mappers.go for the new resource, as well as new tests.

Each of these steps is discussed in more detail below.

**Note**: terraform-validator can only support resources that are supported by the GA terraform provider, not beta resources.

### 1. Magic Modules

The goal of Magic Modules is to auto-generate code targeting specific API endpoints using [yaml files](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products) which are grouped by product.
By default, those yaml files will also be used to generate code for terraform-google-conversion.
A `terraform.yaml` file can specify `exclude_validator: true` on a resource to skip terraform-google-conversion autogeneration, or `exclude_resource: true` to skip autogeneration for both terraform-google-conversion and the providers.

If terraform-google-conversion code can't be auto-generated based on a yaml file, you can instead place a handwritten file in the [`magic-modules/mmv1/third_party/validator` folder](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/third_party/validator). Most resources will only need a conversion func, which should look something like:

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
			Type: "sql.googleapis.com/Database",
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

To generate terraform-google-conversion code locally, run the following from `magic-modules/mmv1`:

```
bundle exec compiler -a -e terraform -f validator -o "/path/to/your/terraform-google-conversion"
```

You can then run `make test` inside your terraform-google-conversion repository to make sure the tests pass prior to creating your PR.

### 2. Terraform Validator

1. Run `go get github.com/GoogleCloudPlatform/terraform-google-conversion` to update the version of terraform-google-conversion in use. (You can also use a [`replace` directive](https://golang.org/ref/mod#go-mod-file-replace) to use your local copy of the repository.)
2. Add a new entry to [mappers.go](../converters/google/mappers.go) for the new resource.

You can now build the binary (with `make build`) and test it. One way to do this would be to create a test project following the instructions in the [policy library user guide](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#for-local-development-environments) (but using the binary you just built.) It's easiest to use a [GCPAlwaysViolatesConstraintV1](https://github.com/GoogleCloudPlatform/terraform-validator/blob/master/testdata/sample_policies/always_violate/policies/constraints/always_violates.yaml) constraint for testing new resources; this is what the tests do. `terraform-validator convert tfplan.json` can show you what terraform-validator thinks the converted Asset looks like.

Be sure to add test cases to [test/cli_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/cli_test.go#L52) and [test/read_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/read_test.go#L24). The test names refer to files in [testdata/templates](https://github.com/GoogleCloudPlatform/terraform-validator/tree/master/testdata/templates). You will generally need to add the following files:
   - A .tf file.
   - A .tfplan.json file.
   - A .json file (representing the output of `terraform-validator convert`)

See [Getting started](./getting_started.md) for details on running tests.