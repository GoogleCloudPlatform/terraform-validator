# Contributing

If you want to contribute to Terraform Validator, check out the [contribution guidelines](../../CONTRIBUTING.md)

## Table of Contents

- [Example project](#example-project)
- [Convert command](#convert-command)
- [Testing](#testing)
  - [Docker](#docker)
- [Add a new resource](#add-a-new-resource)

## Example project

The `example/` directory contains a basic Terraform config for testing the validator. Fully running the validator will require setting up a local [policy library](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library); however, this is not required to test conversion of terraform resources to CAI Assets.

```

cd example/

# Set default credentials.
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/credentials.json

# Set a project and org to test with
export TF_VAR_project_id=my-project-id
export TF_VAR_org_id=12345678

# Generate a terraform plan.
terraform init
terraform plan --out=tfplan.tfplan

# Plan JSON representation.
terraform show -json ./tfplan.tfplan > ./tfplan.json
```

## Convert command

It can be useful to run the convert command separately to test conversion of terraform resources to CAI assets. After configuring the example project as described above, you can run (from the repository root):

```
make build
bin/terraform-validator convert example/tfplan.json
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
```

### Add a new resource

See [Add a new resource](./add_new_resource.md)
