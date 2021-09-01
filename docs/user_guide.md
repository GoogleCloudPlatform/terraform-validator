# User Guide

## Basic setup

These instructions assume you have a working set of terraform files and have already [installed Terraform Validator](./install.md) and [have an organizational policy library](./policy_library.md) available on the same filesystem.

Terraform Validator takes [terraform plan JSON](https://www.terraform.io/docs/internals/json-format.html) as its input. You can generate this file by running the following in your terraform directory:

```
terraform plan -out=tfplan.tfplan
terraform show -json ./tfplan.tfplan > ./tfplan.json
```

## Auth

`terraform-validator` supports the same environment variables for authentication used by the [`google` provider for terraform](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#authentication).

In particular, you can use the following environment variables (in order of precedence) to provide a [service account key file](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#full-reference):


- `GOOGLE_CREDENTIALS`
- `GOOGLE_CLOUD_KEYFILE_JSON`
- `GOOGLE_KEYFILE_JSON`

Using Terraform-Validator-specific [service accounts](https://cloud.google.com/docs/authentication/getting-started) is the recommended practice when using Terraform Validator.

You can also authenticate using an [OAuth 2.0 access token](https://developers.google.com/identity/protocols/OAuth2), which can be provided via the `GOOGLE_OAUTH_ACCESS_TOKEN` environment variable.

For local development, you can also use [Google Application Default Credentials](https://cloud.google.com/docs/authentication/production) by providing the path to your application default credentials file via the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

```
gcloud auth application-default login  # local development only
GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json
```

### Service account impersonation

You can specify a [service account to impersonate](https://cloud.google.com/iam/docs/impersonating-service-accounts) for all Google API calls with the `GOOGLE_IMPERSONATE_SERVICE_ACCOUNT` environment variable.

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
