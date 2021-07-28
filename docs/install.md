# Install Terraform Validator

Terraform Validator is compatible with Terraform 0.12+.

The released binaries are available under the `gs://terraform-validator` Google
Cloud Storage bucket for Linux, Windows, and Mac. They are organized by release,
for example:

```
$ gsutil ls -r "gs://terraform-validator/releases/v*"
...
gs://terraform-validator/releases/v0.4.0/:
gs://terraform-validator/releases/v0.4.0/terraform-validator-darwin-amd64
gs://terraform-validator/releases/v0.4.0/terraform-validator-linux-amd64
gs://terraform-validator/releases/v0.4.0/terraform-validator-windows-amd64
```

To download the binary, you need to
[install](https://cloud.google.com/storage/docs/gsutil_install#install) the
`gsutil` tool first. The following command downloads the Linux version of
Terraform Validator from vX.X.X release to your local directory:

```
gsutil cp gs://terraform-validator/releases/vX.X.X/terraform-validator-linux-amd64 .
chmod 755 terraform-validator-linux-amd64
```

The full list of releases, with release notes, is available [on Github](https://github.com/GoogleCloudPlatform/terraform-validator/releases).
