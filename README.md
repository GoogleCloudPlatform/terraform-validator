# Terraform Validator

Terraform Validator is a tool for validating compliance with organizational policies prior to applying a terraform plan.
It can be used either as a standalone tool or in conjunction with [Forseti](https://forsetisecurity.org/) or other policy enforcement tooling.
Terraform Validator relies on policies that are [compatible with Config Validator](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-set-up-constraints-with-policy-library). For examples, see the [Policy Library](https://github.com/forseti-security/policy-library).

Terraform Validator is compatible with Terraform 0.12+.

**Note**: Using Terraform Validator does _not_ require an active installation of Forseti. Terraform Validator is a self-contained binary.

## Supported Terraform resources
The follow Terraform resources are supported for running validation checks:

```
google_bigquery_dataset
google_bigtable_instance
google_compute_disk
google_compute_forwarding_rule
google_compute_global_forwarding_rule
google_compute_firewall
google_compute_instance
google_compute_network
google_container_cluster
google_container_node_pool
google_filestore_instance
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
google_pubsub_topic
google_pubsub_subscription
google_spanner_instance
google_sql_database_instance
google_storage_bucket
google_storage_bucket_iam_binding
google_storage_bucket_iam_member
google_storage_bucket_iam_policy
google_kms_crypto_key
google_kms_key_ring
```

If you want terraform validator to support an additional resource, please [open an enhancement request](https://github.com/GoogleCloudPlatform/terraform-validator/issues/new?assignees=&labels=enhancement&template=enhancement.md) or follow the instructions below to contribute code.

## Getting started

For instructions on downloading a binary for use on your development machine or CI/CD pipeline, please read the [user guide](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md#how-to-use-terraform-validator).

If you want to contribute to Terraform Validator, check out the [contribution guidelines](./CONTRIBUTING.md) and read the [Getting started docs](./docs/getting_started.md).

## Adding support for a new resource

See [Adding support for a new resource](./docs/add_new_resource.md).

## Disclaimer

This is not an officially supported Google product.
