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

package tfgcv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/config-validator/pkg/api/validator"
	"github.com/GoogleCloudPlatform/config-validator/pkg/gcv"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
)

type ValidateAssetsFunc func(ctx context.Context, assets []google.Asset, policyRootPath string) ([]*validator.Violation, error)

// ValidateAssets instantiates GCV and audits CAI assets using "policies"
// and "lib" folder under policyRootPath.
func ValidateAssets(ctx context.Context, assets []google.Asset, policyRootPath string) ([]*validator.Violation, error) {
	policiesPath := filepath.Join(policyRootPath, "policies")
	libPath := filepath.Join(policyRootPath, "lib")
	_, err := os.Stat(policiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read files in %s", policiesPath)
	}
	_, err = os.Stat(libPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read files in %s", libPath)
	}
	return ValidateAssetsWithLibrary(ctx, assets,
		[]string{policiesPath},
		libPath)
}

// ValidateAssetsWithLibrary instantiates GCV and audits CAI assets.
func ValidateAssetsWithLibrary(ctx context.Context, assets []google.Asset, policyPaths []string, policyLibraryDir string) ([]*validator.Violation, error) {
	valid, err := gcv.NewValidator(policyPaths, policyLibraryDir)
	if err != nil {
		return nil, fmt.Errorf("initializing gcv validator: %w", err)
	}

	pbAssets := make([]*validator.Asset, len(assets))
	for i := range assets {
		pbAssets[i] = &validator.Asset{}
		if err := protoViaJSON(assets[i], pbAssets[i]); err != nil {
			return nil, fmt.Errorf("converting asset %s to proto: %w", assets[i].Name, err)
		}
	}

	pbSplitAssets := splitAssets(pbAssets)

	// Make an empty slice, not a nil slice, so that this
	// can be properly serialized to JSON.
	violations := []*validator.Violation{}
	for _, asset := range pbSplitAssets {
		newViolations, err := valid.ReviewAsset(context.Background(), asset)

		if err != nil {
			return nil, fmt.Errorf("reviewing asset %s: %w", asset, err)
		}
		violations = append(violations, newViolations...)
	}

	return violations, nil
}

// splitAssets split assets because for the GCP target Constraint
// Framework ReviewAsset call an asset must have only one of:
// resource, iam policy, org policy or access context policy
func splitAssets(pbAssets []*validator.Asset) []*validator.Asset {

	pbSplitAssets := make([]*validator.Asset, 0)

	for _, asset := range pbAssets {
		if asset.Resource != nil {
			splitAsset := *asset
			splitAsset.IamPolicy = nil
			splitAsset.OrgPolicy = nil
			splitAsset.AccessContextPolicy = nil
			splitAsset.OrgPolicyPolicy = nil
			pbSplitAssets = append(pbSplitAssets, &splitAsset)
		}
		if asset.IamPolicy != nil {
			splitAsset := *asset
			splitAsset.Resource = nil
			splitAsset.OrgPolicy = nil
			splitAsset.AccessContextPolicy = nil
			splitAsset.OrgPolicyPolicy = nil
			pbSplitAssets = append(pbSplitAssets, &splitAsset)
		}
		if asset.OrgPolicy != nil {
			splitAsset := *asset
			splitAsset.Resource = nil
			splitAsset.IamPolicy = nil
			splitAsset.AccessContextPolicy = nil
			splitAsset.OrgPolicyPolicy = nil
			pbSplitAssets = append(pbSplitAssets, &splitAsset)
		}
		if asset.AccessContextPolicy != nil {
			splitAsset := *asset
			splitAsset.Resource = nil
			splitAsset.IamPolicy = nil
			splitAsset.OrgPolicy = nil
			splitAsset.OrgPolicyPolicy = nil
			pbSplitAssets = append(pbSplitAssets, &splitAsset)
		}
		if asset.OrgPolicyPolicy != nil {
			splitAsset := *asset
			splitAsset.Resource = nil
			splitAsset.IamPolicy = nil
			splitAsset.AccessContextPolicy = nil
			splitAsset.OrgPolicy = nil
			pbSplitAssets = append(pbSplitAssets, &splitAsset)
		}
	}
	return pbSplitAssets
}
