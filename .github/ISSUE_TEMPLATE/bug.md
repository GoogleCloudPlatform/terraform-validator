---
name: Bug
labels: bug
about: For when something is there, but doesn't work how it should.

---

<!--- Please keep this note for the community --->

### Community Note

* Please vote on this issue by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original issue to help the community and maintainers prioritize this request
* Please do not leave _+1_ or _me too_ comments; they generate extra noise for issue followers and do not help prioritize the request
* If you are interested in working on this issue or have submitted a pull request, please leave a comment.
* If the issue is assigned to a user, that user is claiming responsibility for the issue.

<!--- Thank you for keeping this note for the community --->

### Terraform Validator version

<!--- This is the version of terraform-validator you downloaded, or the SHA if you are building from source yourself --->

terraform-validator: vX.X.X

### Affected Resource(s)

<!--- Please list the affected Terraform resources --->

* google_XXXXX

### Terraform Plan JSON

<!--- Information about code formatting: https://help.github.com/articles/basic-writing-and-formatting-syntax/#quoting-code --->

```json
# Copy-paste your Terraform plan JSON here
#
# Ideally this would be a minimal plan that reproduces your issue.
# For large plan files, please use a service like Dropbox and share a link to the ZIP file.
```

### Debug Output

<!---
Please provide a link to a GitHub Gist containing the complete debug output. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

To obtain the debug output, run your terraform-validator command with the `--verbose` option.
--->

### Expected Behavior

<!--- What should have happened? --->

### Actual Behavior

<!--- What actually happened? --->

### Steps to Reproduce

<!--- Please list the steps required to reproduce the issue. --->

1. `terraform-validator convert tfplan.json`

### Important Factoids

<!--- Are there anything atypical about your use case that we should know? --->

### References

<!---
Information about referencing Github Issues: https://help.github.com/articles/basic-writing-and-formatting-syntax/#referencing-issues-and-pull-requests

Are there any other GitHub issues (open or closed) or pull requests that should be linked here? Vendor documentation?
--->

* #0000

<!---
Note Google Cloud customers who are working with a dedicated Technical Account Manager / Customer Engineer: to expedite the investigation and resolution of this issue, please refer to these instructions: https://github.com/hashicorp/terraform-provider-google/wiki/Customer-Contact#raising-gcp-internal-issues-with-the-provider-development-team
--->
