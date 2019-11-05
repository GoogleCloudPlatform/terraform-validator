// +build !tf_0_11

package tfgcv

import "github.com/hashicorp/terraform-plugin-sdk/terraform"

// parseProviderProject is no longer supported in v0.12+.
func parseProviderProject(plan *terraform.Plan) (string, error) {
	return "", ErrParsingProviderProject
}
