# Terraform Validator

This tool is used to validate terraform plans before they are applied. Validations are ran using Forseti Config Validator.

Note: this tool supports Terraform v0.12 by default. To switch to use Terraform v0.11, please see the section [Terraform v0.11](#terraform-v011).

## Getting Started

To get started with Terraform Validator, please follow the [user guide](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-use-terraform-validator).

## Example Usage

See the [Auth](#Auth) section first.


### Steps similar both for Terraform v0.11 and v0.12 versions

```
# The example/ directory contains a basic Terraform config for testing the validator.
cd example/

# Set default credentials.
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/credentials.json

# Set a project to test with
export TF_VAR_project_id=my-project-id

# Set the local forseti-config-policies repository path.
export POLICY_PATH=/path/to/your/forseti-config-policies/repo

# Generate a terraform plan.
terraform plan --out=terraform.tfplan

```

### Terraform v0.11

```
# Switch to use Terraform v0.11 dependencies.
make prepare-v11

# Then run the make command as usual.
make build

# Validate the google resources the plan would create.
terraform-validator validate --policy-path=${POLICY_PATH} ./terraform.tfplan

# Apply the validated plan.
terraform apply ./terraform.tfplan
```

```
# Restore to use Terraform v0.12.
make prepare-v12 build
```

### Terraform v0.12

For 0.12 Terraform release validator required plan exported in JSON format

```
# Plan JSON representation. 
terraform show -json ./terraform.tfplan > ./terraform.tfplan.json

# Validate the google resources the plan would create.
terraform-validator validate --policy-path=${POLICY_PATH} ./terraform.tfplan.json
```

### Apply validated plan

```
# Apply the validated plan.
terraform apply ./terraform.tfplan
```

## Resources
The follow Terraform resources are supported for running validation checks:

- `google_compute_disk`
- `google_compute_instance`
- `google_compute_firewall`
- `google_storage_bucket`
- `google_sql_database_instance`
- `google_project`
- `google_organization_iam_policy`
- `google_organization_iam_binding`
- `google_organization_iam_member`
- `google_folder_iam_policy`
- `google_folder_iam_binding`
- `google_folder_iam_member`
- `google_project_iam_policy`
- `google_project_iam_binding`
- `google_project_iam_member`

## Testing

### Unit

```
make test
```

### Integration

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
