// +build tf_0_11

package tfgcv

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// parseProviderProject attempts to parse hardcoded "project" configuration
// from the "google" provider block. It is lazy and fails if the job involves
// interpolation.
// TODO: Replicate/incorporate terraform interpolation (or is that a good idea?).
func parseProviderProject(plan *terraform.Plan) (string, error) {
	for _, cfg := range plan.Module.Config().ProviderConfigs {
		if cfg.Name == "google" {
			inf, ok := cfg.RawConfig.Raw["project"]
			if !ok {
				continue
			}
			prj := inf.(string)

			// If the provider has a hardcoded project string, return it.
			if !strings.Contains(prj, "${") {
				return prj, nil
			}

			return "", ErrParsingProviderProject
		}
	}

	// If we have reached this point, there was no provider-level project that
	// was specified in this plan. This means the plan should be viable based
	// on resource-level project fields being set.
	return "", nil
}
