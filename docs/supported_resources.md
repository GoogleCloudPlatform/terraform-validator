Terraform Validator is archived.

To enforce policy compliance with Constraint Framework policies as part of a CI/CD pipeline, [migrate to `gcloud beta terraform vet`](https://cloud.google.com/docs/terraform/policy-validation/migrate-from-terraform-validator).

For a library that converts terraform plan data to CAI Asset data, use https://github.com/GoogleCloudPlatform/terraform-google-conversion.

# Supported Resources

*Note: This may not reflect the resources supported by your binary. For OSS builds (available for versions up to 0.13.0), run `terraform-validator list-supported-resources` to get an up-to-date list for your binary.*

The full list of supported resources can be found in our [GCP Terraform docs](https://cloud.google.com/docs/terraform/policy-validation/create-cai-constraints#supported_resources).
