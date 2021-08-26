# Creating a policy library

Your policy library contains all the constraints that you want to be validated.
A single policy library repository can be used for all tools that are compatible with
Constraint Framework, allowing you to enforce policies consistently across all the
stages of the deployment lifecycle.

## Table of Contents

- [Struture](#structure)
- [Quickstart](#quickstart)
- [Advanced usage](#advanced-usage)
- [Periodic updates](#periodic-updates)

## Structure

A Policy Library repository contains the following directories:

*   `policies`
    *   `constraints`: This is initially empty. You should place your constraint
        files here.
    *   `templates`: This directory contains pre-defined constraint templates.
*   `validator`: This directory contains the `.rego` files and their associated
    unit tests. You do not need to touch this directory unless you intend to
    modify existing constraint templates or create new ones. Running `make
    build` will inline the Rego content in the corresponding constraint template
    files.

## Quickstart

### 1. Duplicate Policy Library Repository

Google provides a [sample repository](https://github.com/GoogleCloudPlatform/policy-library)
with a set of pre-defined constraint templates. You can duplicate this repository
into a private repository. First you should create a new **private** git repository.
For example, if you use GitHub then you can use the [GitHub UI](https://github.com/new).
Then follow the steps below to get everything setup.

This policy library can also be made public, but it is not recommended. By
making your policy library public, it would be allowing others to see what you
are and __ARE NOT__ scanning for.

To run the following commands, you will need to configure git to connect
securely. It is recommended to connect with SSH. [Here is a helpful resource](https://help.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh) for learning about how
this works, including steps to set this up for GitHub repositories; other
providers offer this feature as well.

```
export GIT_REPO_ADDR="git@github.com:${YOUR_GITHUB_USERNAME}/policy-library.git"
git clone --bare https://github.com/GoogleCloudPlatform/policy-library.git
cd policy-library.git
git push --mirror ${GIT_REPO_ADDR}
cd ..
rm -rf policy-library.git
git clone ${GIT_REPO_ADDR}
```

### 2. Setup Constraints

Then you need to examine the available constraint templates inside the
`templates` directory. Pick the constraint templates that you wish to use,
create constraint YAML files corresponding to those templates, and place them
under `policies/constraints`. Commit the newly created constraint files to
**your** Git repository. For example, assuming you have created a Git repository
named "policy-library" under your GitHub account, you can use the following
commands to perform the initial commit:

```
cd policy-library
# Add new constraints...
git add --all
git commit -m "Initial commit of policy library constraints"
git push -u origin master
```

### 3. Instantiate constraints

The constraint template library only contains templates. Templates specify the
constraint logic, and you must create constraints based on those templates in
order to enforce them. Constraint parameters are defined as YAML files in the
following format:

```
apiVersion: constraints.gatekeeper.sh/v1alpha1
kind: # place constraint template kind here
metadata:
  name: # place constraint name here
spec:
  severity: # low, medium, or high
  match:
    target: [] # put the constraint application target here
    exclude: [] # optional, default is no exclusions
  parameters: # put the parameters defined in constraint template here
```

The <code><em>target</em></code> field is specified in a path-like format. It
specifies where in the GCP resources hierarchy the constraint is to be applied.
For example:

<table>
  <tr>
   <td>Target
   </td>
   <td>Description
   </td>
  </tr>
  <tr>
   <td>organizations/**
   </td>
   <td>All organizations
   </td>
  </tr>
  <tr>
   <td>organizations/123/**
   </td>
   <td>Everything in organization 123
   </td>
  </tr>
  <tr>
   <td>organizations/123/folders/**
   </td>
   <td>Everything in organization 123 that is under a folder
   </td>
  </tr>
  <tr>
   <td>organizations/123/folders/456
   </td>
   <td>Everything in folder 456 in organization 123
   </td>
  </tr>
  <tr>
   <td>organizations/123/folders/456/projects/789
   </td>
   <td>Everything in project 789 in folder 456 in organization 123
   </td>
  </tr>
</table>

The <code><em>exclude</em></code> field follows the same pattern and has
precedence over the <code><em>target</em></code> field. If a resource is in
both, it will be excluded.

The schema of the <code><em>parameters</em></code> field is defined in the
constraint template, using the
[OpenAPI V3](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#schemaObject)
schema. This is the same validation schema in Kubernetes's custom resource
definition. Every template contains a <code><em>validation</em></code> section
that looks like the following:

```
validation:
  openAPIV3Schema:
    properties:
      mode:
        type: string
      instances:
        type: array
        items: string
```

According to the template above, the parameter field in the constraint file
should contain a string named `mode` and a string array named
<code><em>instances</em></code>. For example:

```
parameters:
  mode: allowlist
  instances:
    - //compute.googleapis.com/projects/test-project/zones/us-east1-b/instances/one
    - //compute.googleapis.com/projects/test-project/zones/us-east1-b/instances/two
```

These parameters specify that two VM instances may have external IP addresses.
The are exempt from the constraint since they are allowlisted.

Here is a complete example of a sample external IP address constraint file:

```
apiVersion: constraints.gatekeeper.sh/v1alpha1
kind: GCPExternalIpAccessConstraintV1
metadata:
  name: forbid-external-ip-allowlist
spec:
  severity: high
  match:
    target: ["organizations/**"]
  parameters:
    mode: "allowlist"
    instances:
    - //compute.googleapis.com/projects/test-project/zones/us-east1-b/instances/one
    - //compute.googleapis.com/projects/test-project/zones/us-east1-b/instances/two
```

## Advanced usage

If the existing constraint templates don't meet your needs, you can also [build your own constraint templates](https://github.com/GoogleCloudPlatform/policy-library/blob/master/docs/constraint_template_authoring.md).

## Periodic updates

Periodically you should pull any changes from the public repository, which might
contain new templates and Rego files.

```
git remote add public https://github.com/GoogleCloudPlatform/policy-library.git
git pull public master
git push origin master
```
