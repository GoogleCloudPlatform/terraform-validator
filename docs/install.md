Terraform Validator is not an officially supported Google product; it is a library for conversion of Terraform plan data to CAI Assets.

[`gcloud beta terraform vet`](https://cloud.google.com/docs/terraform/policy-validation) is a tool for enforcing policy compliance as part of an infrastructure CI/CD pipeline. It is in [Public Preview](https://cloud.google.com/products#product-launch-stages). It uses Terraform Validator and provides additional functionality.


_If you have been using terraform-validator directly in the past, **we recommend [migrating](https://cloud.google.com/docs/terraform/policy_validation/migrate_from_terraform_validator) to `gcloud beta terraform vet`.**_

# Install Terraform Validator (Legacy)

Terraform Validator is compatible with Terraform 0.12+.

The released binaries are available under the `gs://terraform-validator` Google
Cloud Storage bucket for Linux, Windows, and Mac. They are organized by release,
for example:

```
$ gsutil ls -r "gs://terraform-validator/releases/v*"
...
gs://terraform-validator/releases/v0.13.0/:
gs://terraform-validator/releases/v0.13.0/terraform-validator_darwin_amd64-0.13.0.tar.gz
gs://terraform-validator/releases/v0.13.0/terraform-validator_darwin_arm64-0.13.0.tar.gz
gs://terraform-validator/releases/v0.13.0/terraform-validator_linux_amd64-0.13.0.tar.gz
gs://terraform-validator/releases/v0.13.0/terraform-validator_linux_arm64-0.13.0.tar.gz
gs://terraform-validator/releases/v0.13.0/terraform-validator_windows_amd64-0.13.0.tar.gz
```

To download the binary, you need to
[install](https://cloud.google.com/storage/docs/gsutil_install#install) the
`gsutil` tool first. The following commands download and uncompress the Linux AMD64
version of Terraform Validator from v0.13.0 release to your local directory:

```
gsutil cp gs://terraform-validator/releases/v0.13.0/terraform-validator_linux_amd64-0.13.0.tar.gz .
tar -xzvf terraform-validator_linux_amd64-0.13.0.tar.gz
chmod 755 terraform-validator
```

The full list of releases, with release notes, is available [on Github](https://github.com/GoogleCloudPlatform/terraform-validator/releases).

Binary builds are only available for versions up to 0.13.0. 0.14.0+ are only available via [`gcloud beta terraform vet`](https://cloud.google.com/docs/terraform/policy_validation)

## Disclaimer

This is not an officially supported Google product.
