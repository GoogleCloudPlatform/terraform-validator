// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package policy provides utilities for easy operations on IAM policy in terraform.
package policy

import (
	"context"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/pkg/errors"
)

// Overlay contains the policies that have changed in the plan.
// It is a map from the resource path to the new policy.
// Only resources with a change in policy are included.
type Overlay map[string]*google.IAMPolicy

// BuildOverlay creates the Overlay for the current plan.
func BuildOverlay(path, project, ancestry string, offline bool) (Overlay, error) {
	ctx := context.Background()

	before, err := tfgcv.ReadCurrentAssets(ctx, path, project, ancestry, offline)
	if err != nil {
		if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
			return nil, errors.New("unable to parse provider project, please use --project flag")
		}
		return nil, errors.Wrap(err, "converting tfplan to CAI assets")
	}

	after, err := tfgcv.ReadPlannedAssets(ctx, path, project, ancestry, offline)
	if err != nil {
		if errors.Cause(err) == tfgcv.ErrParsingProviderProject {
			return nil, errors.New("unable to parse provider project, please use --project flag")
		}
		return nil, errors.Wrap(err, "converting tfplan to CAI assets")
	}

	return extractOverlay(before, after), nil
}

// extractOverlay compares the current and planned assets to determine
// what policies have changed.
func extractOverlay(before, after []google.Asset) Overlay {
	emptyPolicy := func() *google.IAMPolicy {return &google.IAMPolicy{Bindings:[]google.IAMBinding{}}}

	afterIndex := make(map[string]google.Asset)
	for _, a := range after {
		afterIndex[a.Name] = a
	}
	beforeIndex := make(map[string]google.Asset)
	for _, b := range before {
		beforeIndex[b.Name] = b
	}

	// Cover updates to policy and deletions of policy or resource.
	overlay := Overlay{}
	for name, b := range beforeIndex {
		a, ok := afterIndex[name]
		switch {
		case !ok:
			// If it doesn't exist in the after, it means the entire resource, not just
			// the policy, is being deleted.
			if b.IAMPolicy != nil {
				// We don't treat this policy deletion any differently.  It's still
				// a policy deletion, even if more than the policy is being deleted.
				overlay[name] = emptyPolicy()
			}
		case a.IAMPolicy != b.IAMPolicy:
			if a.IAMPolicy != nil {
				overlay[name] = a.IAMPolicy
			} else {
				// If the policy exists in the before and not in the after, then the
				// policy has been deleted. The resource has not been deleted because
				// that case would have already been caught above.
				overlay[name] = emptyPolicy()
			}
		}
	}

	// Cover the case where new resources were created.
	for name, a := range afterIndex {
		_, ok := beforeIndex[name]
		if !ok && a.IAMPolicy != nil {
			overlay[name] = a.IAMPolicy
		}
	}

	return overlay
}
