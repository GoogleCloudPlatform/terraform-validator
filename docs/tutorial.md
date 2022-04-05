A gcloud integration for Terraform Validator is in [Public Preview](https://cloud.google.com/products#product-launch-stages) at [`gcloud beta terraform vet`](https://cloud.google.com/docs/terraform/policy_validation). We recommend [migrating](https://cloud.google.com/docs/terraform/policy_validation/migrate_from_terraform_validator) to that tool.

# Tutorial

In this tutorial, you will apply a constraint that enforces IAM policy member
domain restriction using [Cloud Shell](https://cloud.google.com/shell/).

First click on this
[link](https://console.cloud.google.com/cloudshell/open?cloudshell_image=gcr.io/graphite-cloud-shell-images/terraform:latest&cloudshell_git_repo=https://github.com/GoogleCloudPlatform/policy-library.git)
to open a new Cloud Shell session. The Cloud Shell session has Terraform
pre-installed and the Policy Library repository cloned. Once you have the
session open, the next step is to copy over the sample IAM domain restriction
constraint:

```
cp samples/iam_service_accounts_only.yaml policies/constraints
```

Let's take a look at this constraint:

```
apiVersion: constraints.gatekeeper.sh/v1alpha1
kind: GCPIAMAllowedPolicyMemberDomainsConstraintV1
metadata:
  name: service_accounts_only
spec:
  severity: high
  match:
    target: ["organizations/**"]
  parameters:
    domains:
      - gserviceaccount.com
```

It specifies that only members from gserviceaccount.com domain can be present in
an IAM policy. To verify that it works, let's attempt to create a project.
Create the following Terraform `main.tf` file:

```
provider "google" {
  version = "~> 1.20"
  project = "your-terraform-provider-project"
}

resource "random_id" "proj" {
  byte_length = 8
}

resource "google_project" "sample_project" {
  project_id      = "validator-${random_id.proj.hex}"
  name            = "config validator test project"
}

resource "google_project_iam_binding" "sample_iam_binding" {
  project = "${google_project.sample_project.project_id}"
  role    = "roles/owner"

  members = [
    "user:your-email@your-domain"
  ]
}

```

Make sure to specify your Terraform
[provider project](https://www.terraform.io/docs/providers/google/getting_started.html)
and email address. Then initialize Terraform and generate a Terraform plan:

```
terraform init
terraform plan -out=test.tfplan
terraform show -json ./test.tfplan > ./tfplan.json
```

Since your email address is in the IAM policy binding, the plan should result in
a violation. Let's try this out:

```
gsutil cp gs://terraform-validator/releases/v0.12.5/terraform-validator_linux_amd64-0.12.5.tar.gz .
tar -xzvf terraform-validator_linux_amd64-0.12.5.tar.gz
chmod 755 terraform-validator
./terraform-validator validate tfplan.json --policy-path=policy-library
```

The Terraform validator should return a violation. As a test, you can relax the
constraint to make the violation go away. Edit the
`policy-library/policies/constraints/iam_service_accounts_only.yaml` file and
append your email domain to the domains allowlist:

```
apiVersion: constraints.gatekeeper.sh/v1alpha1
kind: GCPIAMAllowedPolicyMemberDomainsConstraintV1
metadata:
  name: service_accounts_only
spec:
  severity: high
  match:
    target: ["organizations/**"]
  parameters:
    domains:
      - gserviceaccount.com
      - your-domain-here
```

Then run Terraform plan and validate the output again:

```
terraform plan -out=test.tfplan
terraform show -json ./test.tfplan > ./tfplan.json
./terraform-validator validate tfplan.json --policy-path=policy-library
```

The command above should result in no violations found.
