// Copyright 2021 Google LLC
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

package cnvconfig

import (
	"context"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/pkg/errors"
)

func GetConfig(ctx context.Context, project string, offline bool) (*resources.Config, error) {
	cfg := &resources.Config{
		Project: project,
	}

	// Search for default credentials
	cfg.Credentials = multiEnvSearch([]string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
	})

	cfg.AccessToken = multiEnvSearch([]string{
		"GOOGLE_OAUTH_ACCESS_TOKEN",
	})

	cfg.ImpersonateServiceAccount = multiEnvSearch([]string{
		"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	})

	if !offline {
		resources.ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
	}

	return cfg, nil
}
