# Terraform Validator

This tool is used to validate terraform plans before they are applied. Validations are based on policies from the Config Validator [Policy Library](https://github.com/forseti-security/policy-library).

**Note**: Using Terraform Validator does _not_ require an active installation of Forseti. Terraform Validator is a self-contained binary.

Note: this tool supports Terraform v0.12+.

## Getting Started

To get started with Terraform Validator, please follow the [user guide](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-use-terraform-validator).

## Example Usage

See the [Auth](#Auth) section first.


### Terraform 0.12+ Usage

```
# The example/ directory contains a basic Terraform config for testing the validator.
cd example/

# Set default credentials.
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/credentials.json

# Set a project and org to test with
export TF_VAR_project_id=my-project-id
export TF_VAR_org_id=93392932

# Set the local forseti-config-policies repository path.
export POLICY_PATH=/path/to/your/forseti-config-policies/repo

# Generate a terraform plan.
terraform plan --out=terraform.tfplan

# Plan JSON representation.
terraform show -json ./terraform.tfplan > ./terraform.tfplan.json

# Validate the google resources the plan would create.
terraform-validator validate --policy-path=${POLICY_PATH} ./terraform.tfplan.json

# Apply the validated plan.
terraform apply ./terraform.tfplan
```

## Resources
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
google_spanner_instance
google_sql_database_instance
google_storage_bucket
google_storage_bucket_iam_binding
google_storage_bucket_iam_member
google_storage_bucket_iam_policy
```

## Testing

### Unit

```
make test
```

### Integration

#### Non-docker
```
gcloud auth application-default login
export TEST_PROJECT=my-project-id
export TEST_CREDENTIALS=~/.config/.config/gcloud/application_default_credentials.json
make test-integration
```

#### Docker
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

## Auth

The `terraform` and the `terraform-validator` commands need to be able to authenticate to Google Cloud APIs. This can be done by generating a `credentials.json` file:

https://cloud.google.com/docs/authentication/production

Once you have a credentials file on your local machine, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to the credentials file.

## Disclaimer

This is not an officially supported Google product.
