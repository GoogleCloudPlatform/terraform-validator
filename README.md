# Terraform Validator

Terraform Validator is a tool for validating compliance with organizational policies prior to applying a terraform plan.
It can be used either as a standalone tool or in conjunction with [Forseti](https://forsetisecurity.org/) or other policy enforcement tooling.
Terraform Validator relies on policies that are [compatible with Config Validator](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library). For examples, see the [Policy Library](https://github.com/forseti-security/policy-library).

Terraform Validator is compatible with Terraform 0.12+.

**Note**: Using Terraform Validator does _not_ require an active installation of Forseti. Terraform Validator is a self-contained binary.

## Supported Terraform resources
The follow Terraform resources are supported for running validation checks:

```
google_bigquery_dataset
google_bigtable_instance
google_compute_disk
google_compute_forwarding_rule
google_compute_global_forwarding_rule
google_compute_firewall
google_compute_instance
google_container_cluster
google_container_node_pool
google_filestore_instance
google_folder_iam_binding
google_folder_iam_member
google_folder_iam_policy
google_organization_iam_binding
google_organization_iam_member
google_organization_iam_policy
google_project
google_project_iam_binding
google_project_iam_member
google_project_iam_policy
google_project_organization_policy
google_project_service
google_pubsub_topic
google_pubsub_subscription
google_spanner_instance
google_sql_database_instance
google_storage_bucket
google_storage_bucket_iam_binding
google_storage_bucket_iam_member
google_storage_bucket_iam_policy
google_kms_crypto_key
google_kms_key_ring
```

## Getting Started

For instructions on downloading a binary for use on your development machine or CI/CD pipeline, please read the [user guide](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-use-terraform-validator). The instructions in this README are aimed at developers working on Terraform Validator itself.

### Auth

The `terraform` and `terraform-validator` commands need to be able to authenticate to Google Cloud APIs. This can be done by [generating a `credentials.json` file](https://cloud.google.com/docs/authentication/production). For local development, you can generate application default credentials. For production, use service account credentials instead.

Once you have a credentials file on your local machine, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to the credentials file.

```
gcloud auth application-default login  # local development only
GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json  # or path to service account credentials
```

### Example project

The `example/` directory contains a basic Terraform config for testing the validator. Fully running the validator will require setting up a local [policy library](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library); however, this is not required to test conversion of terraform resources to CAI Assets.

```

cd example/

# Set default credentials.
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/credentials.json

# Set a project and org to test with
export TF_VAR_project_id=my-project-id
export TF_VAR_org_id=93392932

# Generate a terraform plan.
terraform init
terraform plan --out=tfplan.tfplan

# Plan JSON representation.
terraform show -json ./tfplan.tfplan > ./tfplan.json
```

#### Convert command

It can be useful to run the convert command separately to test conversion of terraform resources to CAI assets. After configuring the example project as described above, you can run (from the repository root):

```
make build
bin/terraform-validator convert example/tfplan.json
```

#### Validate command
Running the validate command requires setting up a local [policy library](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library).

```
# Set the local policy library repository path.
export POLICY_PATH=/path/to/your/policy/library

# Build the binary
make build

# Validate the google resources the plan would create.
bin/terraform-validator validate --policy-path=${POLICY_PATH} example/tfplan.json

# Apply the validated plan.
terraform apply example/tfplan.tfplan
```
## Testing

```
# Unit tests
make test

# Integration tests (interacts with real APIs)
gcloud auth application-default login
export TEST_PROJECT=my-project-id
export TEST_CREDENTIALS=~/.config/gcloud/application_default_credentials.json
make test-integration

# Specific integration test
go test -v -run=<test name or prefix> ./test
```

### Docker
It is also possible to run the integration tests inside a Docker container to match the CI/CD pipeline.
First, build the Docker container:

```
make build-docker
```

See the [Auth](#Auth) section for obtaining a credentials file, then start the Docker container:

```
export PROJECT_ID=my-project-id
export GOOGLE_APPLICATION_CREDENTIALS=$(pwd)/credentials.json
make run-docker
```

Finally, run the integration tests inside the container:
```
make test-integration
````

## Adding support for a new resource

We are using code generation tool called [Magic Modules](https://github.com/googleCloudPlatform/magic-modules/) that uses a shared code base to generate the [google](https://github.com/hashicorp/terraform-provider-google) and [google-beta](https://github.com/hashicorp/terraform-provider-google-beta) Terraform providers as well as a library called [terraform-google-conversion](https://github.com/GoogleCloudPlatform/terraform-google-conversion). terraform-google-conversion is what Terraform Validator uses to convert Terraform resources to CAI Assets for validation.

Some Terraform resources are fully generated, whereas some resources are hand written and located in [the third_party/validator/ folder in magic modules](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/third_party/validator/resources). Compilation and copying of files into terraform-google-conversion happens in [mmv1/provider/terraform_object_library.rb](https://github.com/GoogleCloudPlatform/magic-modules/blob/100ba410e1db645a6ae0e6351f87e82e897eade7/mmv1/provider/terraform_object_library.rb).

Adding support for a new resource follows these steps:

0. If the resource is already [supported by terraform-google-conversion](https://github.com/GoogleCloudPlatform/terraform-google-conversion/tree/master/google) you can skip to step 6. You may be need to upgrade the terraform-google-conversion dependency in terraform-validator's [go.mod file](https://github.com/GoogleCloudPlatform/terraform-validator/blob/master/go.mod).
1. Add support for the resource in Magic Modules (preferably auto-generated; hand-written resources are harder to maintain.)
2. [Generate the terraform-google-conversion code](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/README.md#generating-downstream-tools).
3. Run `make test` inside the terraform-google-conversion repository.
4. Create a PR against Magic Modules. This will be reviewed by a core contributor.
5. Once that is merged, go to terraform-validator and run `go get github.com/GoogleCloudPlatform/terraform-google-conversion` to update the version of terraform-google-conversion in use
6. Add one or more [mappers](https://github.com/GoogleCloudPlatform/terraform-validator/blob/86e5a59ce0dbf4089db8484a482dbac4f48dc93a/converters/google/mappers.go#L42) for the new resource.
7. Add test cases to [test/cli_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/cli_test.go#L52) and [test/read_test.go](https://github.com/GoogleCloudPlatform/terraform-validator/blob/c1295c541897e1357eb3e4d93a88d7083ff41c90/test/read_test.go#L24)
8. Run tests.
9. Create a PR against terraform-validator.

### mappers

"mappers" are the glue that connects a specific terraform resource type (like `google_compute_disk`) to specific terraform-google-conversion functions necessary to convert it to a CAI Asset. A mapper can have the following functions:

- `convert`: Required. This function does basic conversion of a Terraform resource to a CAI Asset, including converting nested structures and specifying what the [CAI Asset Type](https://cloud.google.com/asset-inventory/docs/supported-asset-types) is.
- `fetch`, `mergeCreateUpdate`, `mergeDelete`: Optional. Some assets, like [IAM Members and Bindings](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam), have to be merged with remote data prior to validation in order to properly check whether policies are being followed.

## Disclaimer

This is not an officially supported Google product.
