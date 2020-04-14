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
	"path/filepath"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/pkg/errors"
)

// To be set by Go build tools.
var buildVersion string

// BuildVersion returns the build version of Terraform Validator.
func BuildVersion() string {
	return buildVersion
}

// ValidateAssets instantiates GCV and audits CAI assets using "policies"
// and "lib" folder under policyRootPath.
func ValidateAssets(ctx context.Context, assets []google.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return ValidateAssetsWithLibrary(ctx, assets,
		[]string{filepath.Join(policyRootPath, "policies")},
		filepath.Join(policyRootPath, "lib"))
}

// ValidateAssetsWithLibrary instantiates GCV and audits CAI assets.
func ValidateAssetsWithLibrary(ctx context.Context, assets []google.Asset, policyPaths []string, policyLibraryDir string) (*validator.AuditResponse, error) {
	valid, err := gcv.NewValidator(policyPaths, policyLibraryDir)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}

	pbAssets := make([]*validator.Asset, len(assets))
	for i := range assets {
		pbAssets[i] = &validator.Asset{}
		if err := protoViaJSON(assets[i], pbAssets[i]); err != nil {
			return nil, errors.Wrapf(err, "converting asset %s to proto", assets[i].Name)
		}
	}

	auditResult := &validator.AuditResponse{}
	for _, asset := range pbAssets {
		violations, err := valid.ReviewAsset(context.Background(), asset)
		if err != nil {
			return nil, errors.Wrapf(err, "auditing %s", asset)
		}
		auditResult.Violations = append(auditResult.Violations, violations...)
	}

	return auditResult, nil
}
