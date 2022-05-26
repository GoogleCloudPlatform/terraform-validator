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
	"path/filepath"

	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/GoogleCloudPlatform/terraform-validator/ancestrymanager"
	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ReadPlannedAssetsFunc func(ctx context.Context, path, project string, ancestry map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger, userAgent string) ([]google.Asset, error)

// ReadPlannedAssets extracts CAI assets from a terraform plan file.
// If ancestry path is provided, it assumes the project is in that path rather
// than fetching the ancestry information using Google API. If convertUnchanged
// is set then resources that do not have any change from their deployed state
// are also reported in the output, otherwise only resources that are going to
// be changed are reported.
// It ignores non-supported resources.
func ReadPlannedAssets(ctx context.Context, path, project string, ancestry map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger, userAgent string) ([]google.Asset, error) {
	converter, err := newConverter(ctx, path, project, ancestry, offline, convertUnchanged, errorLogger, userAgent)
	if err != nil {
		return nil, err
	}

	data, err := readTF12Data(path)
	if err != nil {
		return nil, err
	}

	changes, err := tfplan.ReadResourceChanges(data)
	if err != nil {
		return nil, err
	}

	err = converter.AddResourceChanges(changes)
	if err != nil {
		return nil, err
	}

	return converter.Assets(), nil
}

func newConverter(ctx context.Context, path, project string, ancestry map[string]string, offline, convertUnchanged bool, errorLogger *zap.Logger, userAgent string) (*google.Converter, error) {
	cfg, err := resources.GetConfig(ctx, project, offline, userAgent)
	if err != nil {
		return nil, errors.Wrap(err, "building google configuration")
	}

	var resourceManager *cloudresourcemanager.Service
	if !offline {
		resourceManager = cfg.NewResourceManagerClient(cfg.UserAgent())
	}
	ancestryManager, err := ancestrymanager.New(resourceManager, ancestry, errorLogger)
	if err != nil {
		return nil, errors.Wrap(err, "building google ancestry manager")
	}
	converter := google.NewConverter(cfg, ancestryManager, offline, convertUnchanged, errorLogger)
	return converter, nil
}

func readTF12Data(path string) ([]byte, error) {
	if ".json" != filepath.Ext(path) {
		return nil, errors.New(fmt.Sprintf("Terraform 0.12 support plans only in JSON format, got: %s", filepath.Ext(path)))
	}
	// JSON format means Terraform 12
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening JSON plan file")
	}
	return data, nil
}
