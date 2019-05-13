# Terraform Validator

This tool is used to validate terraform plans before they are applied. Validations are ran using Forseti Config Validator.

## Installation

```
go install .
```

## Example Usage

See the [Auth](#Auth) section first.

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

# Validate the google resources the plan would create.
terraform-validator validate --policy-path=${POLICY_PATH} ./terraform.tfplan

# Apply the validated plan.
terraform apply ./terraform.tfplan
```

## Resources
The follow Terraform resources are supported for running validation checks:

- `google_compute_disk`
- `google_compute_instance`
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

See the [Auth](#Auth) section for obtaining a credentials file.

```
make test-e2e PROJECT=my-project-id CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS
```

## Auth

The `terraform` and the `terraform-validator` commands need to be able to authenticate to Google Cloud APIs. This can be done by generating a `credentials.json` file:

https://cloud.google.com/docs/authentication/production

Once you have a credentials file on your local machine, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to the credentials file.

## Disclaimer

This is not an officially supported Google product.
