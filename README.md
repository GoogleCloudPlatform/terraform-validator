# Terraform Validator

## Overview

As your business shifts towards an infrastructure-as-code workflow, security and
cloud administrators are concerned about misconfigurations that may cause
security and governance violations.

Cloud Administrators need to be able to put up guardrails that follow security
best practices and help drive the environment towards programmatic security and
governance while enabling developers to go fast.

Terraform Validator allows your administrators to enforce **constraints** on
developer machines and as part of your CI/CD pipeline, allowing you to check for
constraint violations and provide warnings or halt invalid deployments before
they reach production.

### One way to define constraints

Constraints are designed to be compatible with tools across the deployment
lifecycle. The same set of constraints that you use with Terraform Validator
can also be used with any other tool that supports them, either at deploy-time
or as an audit of deployed resources. These constraints live in your
organization's repository as the source of truth for your security and
governance requirements. You can obtain constraints from the
[Policy Library](./docs/policy_library.md), or
[build your own constraint templates](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/constraint_template_authoring.md).

## Table of Contents

- [Install Terraform Validator](./docs/install.md)
- [Tutorial](./docs/tutorial.md)
- [Creating a policy library](./docs/policy_library.md)
- [Using terraform validator](./docs/user_guide.md)
- [Supported resources](./docs/supported_resources.md)
- [Contributing](./docs/contributing/index.md)
  - [Add a new resource](./docs/contributing/add_new_resource.md)

## Disclaimer

This is not an officially supported Google product.
