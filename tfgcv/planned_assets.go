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
	"os"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"github.com/golang/glog"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// ReadPlannedAssets extracts CAI assets from a terraform plan file.
// It ignores non-supported resources.
func ReadPlannedAssets(path, project string) ([]google.Asset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening plan file")
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, errors.Wrap(err, "reading terraform plan")
	}

	// TODO: Pull project from tfplan instead.
	// i.e. terraform.Plan.Module.Config().ProviderConfigs[].RawConfig
	// The complication with pulling from the above config is with uninterpolated
	// terraform variables/locals/etc.
	converter, err := google.NewConverter(project, "")
	if err != nil {
		return nil, errors.Wrap(err, "building google converter")
	}

	for _, r := range tfplan.ComposeResources(plan, converter.Schemas()) {
		if err := converter.AddResource(&r); err != nil {
			if errors.Cause(err) == google.ErrDuplicateAsset {
				glog.Warningf("converting resource: %v", err)
			} else {
				return nil, errors.Wrapf(err, "converting resource %v", r.Kind())
			}
		}
	}

	return converter.Assets(), nil
}
