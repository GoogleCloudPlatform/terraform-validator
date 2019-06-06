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
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"github.com/golang/glog"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// ReadPlannedAssets extracts CAI assets from a terraform plan file.
// If ancestry path is provided, the ancestry cache will be prewarmed for the
// project provided. If offline is set as true, it will avoid fetching resource
// using Google API.
// It ignores non-supported resources.
func ReadPlannedAssets(path, project, ancestry string, offline bool) ([]google.Asset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening plan file")
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, errors.Wrap(err, "reading terraform plan")
	}

	// Attempt to pull the project from the provider.
	if project == "" {
		project, err = parseProviderProject(plan)
		if err != nil {
			return nil, err
		}
	}

	var resourceManager *cloudresourcemanager.Service
	if !offline {
		// Add User Agent string to indicate Terraform Validator usage.
		// Do *NOT* change the "config-validator-tf/" prefix, or else it will
		// break usage tracking.
		ua := option.WithUserAgent(fmt.Sprintf("config-validator-tf/%s", BuildVersion()))
		if resourceManager, err = cloudresourcemanager.NewService(context.Background(), ua); err != nil {
			return nil, errors.Wrap(err, "constructing resource manager client")
		}
	}
	converter, err := google.NewConverter(resourceManager, project, ancestry, "", offline)
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

var ErrParsingProviderProject = errors.New("unable to parse provider project")

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
