# Terraform Validator

Terraform Validator is a tool for validating compliance with organizational policies prior to applying a terraform plan.
It can be used either as a standalone tool or in conjunction with [Forseti](https://forsetisecurity.org/) or other policy enforcement tooling.
Terraform Validator relies on policies that are [compatible with Config Validator](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library). For examples, see the [Policy Library](https://github.com/forseti-security/policy-library).

Terraform Validator is compatible with Terraform 0.12+.

## Supported Terraform resources
The follow Terraform resources are supported for running validation checks:

```
google_bigquery_dataset
google_compute_disk
google_compute_firewall
google_compute_instance
google_container_cluster
google_container_node_pool
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
google_spanner_instance
google_sql_database_instance
google_storage_bucket
google_storage_bucket_iam_binding
google_storage_bucket_iam_member
google_storage_bucket_iam_policy
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

It can be useful to run the convert command separately to test conversion of terraform resources to CAI assets. After configuring the example project as described above, you can run:

```
terraform-validator convert ./tfplan.json
```

#### Validate command
Running the validate command requires setting up a local [policy library](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library).

```
# Set the local forseti-config-policies repository path.
export POLICY_PATH=/path/to/your/policy/library

# Validate the google resources the plan would create.
terraform-validator validate --policy-path=${POLICY_PATH} ./tfplan.json

# Apply the validated plan.
terraform apply ./tfplan.tfplan
```
## Testing

```
# Unit tests
make test

# Integration tests (interacts with real APIs)
gcloud auth application-default login
export TEST_PROJECT=my-project-id
export TEST_CREDENTIALS=~/.config/.config/gcloud/application_default_credentials.json
make test-integration
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

## Disclaimer

This is not an officially supported Google product.
