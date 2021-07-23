# User Guide

## Basic setup

These instructions assume you have a working set of terraform files and have already [installed Terraform Validator](./install.md) and [have an organizational policy library](./policy_library.md) available on the same filesystem.

Terraform Validator takes [terraform plan JSON](https://www.terraform.io/docs/internals/json-format.html) as its input. You can generate this file by running the following in your terraform directory:

```
terraform plan -out=tfplan.tfplan
terraform show -json ./tfplan.tfplan > ./tfplan.json
```

## Auth

The `terraform` and `terraform-validator` commands need to be able to authenticate to Google Cloud APIs. This can be done by [generating a `credentials.json` file](https://cloud.google.com/docs/authentication/production). For local development, you can generate application default credentials. For production, use service account credentials instead.

Once you have a credentials file on your local machine, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to the credentials file.

```
gcloud auth application-default login  # local development only
GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json  # or path to service account credentials
```

## `terraform-validator validate`

This command allows you to validate your terraform plan JSON against a specific policy library.

Basic usage:

```
terraform-validator validate tfplan.json --policy-path=${POLICY_PATH}
```

### Flags

#### `--policy-path=${POLICY_PATH}`

The policy-path flag is set to the local clone of your Git repository that
contains your [organizational constraints and templates](./policy_library.md).

#### `--project=my-project` (optional)

Terraform Validator accepts an optional `--project` flag. This will be used as the default
project when building ancestry paths for any resource that doesn't have an explicit project set.

### Return value

If violations are found, `terraform-validator` will return exit code `2` and display a list
of violations:

```
Found Violations:

Constraint iam_domain_restriction on resource //cloudresourcemanager.googleapis.com/projects/12345678: IAM policy for //cloudresourcemanager.googleapis.com/projects/12345678 contains member from unexpected domain: user:foo@example.com

Constraint iam_domain_restriction on resource //cloudresourcemanager.googleapis.com/projects/12345678: IAM policy for //cloudresourcemanager.googleapis.com/projects/12345678 contains member from unexpected domain: group:bar@example.com
```

If all constraints are validated, the command will return exit code `0` and display
"`No violations found`."
