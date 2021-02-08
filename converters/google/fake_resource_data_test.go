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
package google

import (
	"testing"

	provider "github.com/hashicorp/terraform-provider-google/v3/google"
	"github.com/stretchr/testify/assert"
)

func TestFakeResourceData_kind(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":  "test-disk",
		"type":  "pd-ssd",
		"zone":  "us-central1-a",
		"image": "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
	)
	assert.Equal(t, d.Kind(), "google_compute_disk")
}


func TestFakeResourceData_id(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":  "test-disk",
		"type":  "pd-ssd",
		"zone":  "us-central1-a",
		"image": "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
	)
	assert.Equal(t, d.Id(), "")
}


func TestFakeResourceData_get(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":  "test-disk",
		"type":  "pd-ssd",
		"zone":  "us-central1-a",
		"image": "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
	)
	assert.Equal(t, d.Get("name"), "test-disk")
}


func TestFakeResourceData_getOkOk(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":  "test-disk",
		"type":  "pd-ssd",
		"zone":  "us-central1-a",
		"image": "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
	)
	res, ok := d.GetOk("name")
	assert.Equal(t, res, "test-disk")
	assert.True(t, ok)
}


func TestFakeResourceData_getOkNotOk(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":  "test-disk",
		"type":  "pd-ssd",
		"zone":  "us-central1-a",
		"image": "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceData(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
	)
	res, ok := d.GetOk("incorrect")
	assert.Nil(t, res)
	assert.False(t, ok)
}
