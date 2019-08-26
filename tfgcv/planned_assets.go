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
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"github.com/golang/glog"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// ReadPlannedAssets extracts CAI assets from a terraform plan file.
// If ancestry path is provided, it assumes the project is in that path rather
// than fetching the ancestry information using Google API.
// It ignores non-supported resources.
func ReadPlannedAssets(path, project, ancestry, tfVersion string) ([]google.Asset, error) {

	// Add User Agent string to indicate Terraform Validator usage.
	// Do *NOT* change the "config-validator-tf/" prefix, or else it will
	// break usage tracking.
	ua := option.WithUserAgent(fmt.Sprintf("config-validator-tf/%s", BuildVersion()))
	resourceManager, err := cloudresourcemanager.NewService(context.Background(), ua)
	if err != nil {
		return nil, errors.Wrap(err, "constructing resource manager client")
	}
	converter, err := google.NewConverter(resourceManager, project, ancestry, "")
	if err != nil {
		return nil, errors.Wrap(err, "building google converter")
	}

	var resources []tfplan.Resource

	switch tfVersion {
	case tfplan.TF12:
		if ".json" != filepath.Ext(path) {
			return nil, errors.New(fmt.Sprintf("Terraform 0.12 support plans only in JSON format, got: %s", filepath.Ext(path)))
		}
		// JSON format means Terraform 12
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "opening JSON plan file")
		}

		resources, err = tfplan.ComposeTF12Resources(data, converter.Schemas())
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal from JSON and composing terraform plan")
		}
	case tfplan.TF11:
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrap(err, "opening plan file")
		}
		defer f.Close()

		plan, err := terraform.ReadPlan(f)
		if err != nil {
			return nil, errors.Wrap(err, "reading terraform plan")
		}

		if project == "" {
			return nil, errors.Wrap(err, "please use --project flag")
		}

		resources = tfplan.ComposeResources(plan, converter.Schemas())
	default:
		return nil, errors.New(fmt.Sprintf("Possible values for --tf-version flag are [%s, %s], got: %s", tfplan.TF11, tfplan.TF12, tfVersion))
	}

	for _, r := range resources {
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

