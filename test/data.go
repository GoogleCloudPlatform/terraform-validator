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

package test

import (
	computev1 "google.golang.org/api/compute/v1"
)

// data represents the full dataset that is used for templating terraform
// configs. It contains Google API resources that are expected to be returned
// after converting the terraform plan.
type data struct {
	Provider map[string]string
	computev1.Disk
}

func newData(project, credentials string) data {
	return data{
		Provider: map[string]string{
			"project":     project,
			"credentials": credentials,
		},
		Disk: computev1.Disk{
			Name:        "my-disk",
			Type:        "https://www.googleapis.com/compute/v1/projects/" + project + "/zones/us-central1-a/diskTypes/pd-ssd",
			Zone:        "projects/" + project + "/global/zones/us-central1-a",
			SourceImage: "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
	}
}
