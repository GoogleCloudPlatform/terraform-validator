# Terraform Validator

This tool is used to validate terraform plans before they are applied. Validations are ran using Forseti Config Validator.

## Installation

```
go install .
```

## Example Usage

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
terraform-validator validate --project=${TF_VAR_project_id} --policy-path=${POLICY_PATH} ./terraform.tfplan

# Apply the validated plan.
terraform apply ./terraform.tfplan
```

## Testing

```
# Unit
make test

# Integration
cp <my-google-credentials-file> ./credentials.json
make test-e2e PROJECT=my-project
```

## Disclaimer
This is not an officially supported Google product.
