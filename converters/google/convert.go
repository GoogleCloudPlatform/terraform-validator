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
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	provider "github.com/terraform-providers/terraform-provider-google/google"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

type Asset = converter.Asset

var ErrDuplicateAsset = errors.New("duplicate asset")

// TerraformResource represents the required methods needed to convert a terraform
// resource into an Asset type.
type TerraformResource interface {
	Id() string
	Kind() string
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	GetOkExists(string) (interface{}, bool)
}

// nonImplementedResourceData represents the non-getter fields that are
// passed to terraform providers that are not implemented here but are
// required (and not used) by downstream mappers.
type nonImplementedResourceData struct{}

func (nonImplementedResourceData) HasChange(string) bool         { return false }
func (nonImplementedResourceData) Set(string, interface{}) error { return nil }
func (nonImplementedResourceData) SetId(string)                  {}

// NewConverter is a factory function for Converter.
func NewConverter(project, credentials string) (*Converter, error) {
	cfg := &converter.Config{
		Project:     project,
		Credentials: credentials,
	}
	if err := cfg.LoadAndValidate(); err != nil {
		return nil, errors.Wrap(err, "configuring")
	}

	p := provider.Provider().(*schema.Provider)
	return &Converter{
		schema:      p,
		mapperFuncs: mappers(),
		cfg:         cfg,
		assets:      make(map[string]Asset),
	}, nil
}

// Converter knows how to convert terraform resources to their
// Google CAI (Cloud Asset Inventory) format (the Asset type).
type Converter struct {
	schema *schema.Provider

	// Map terraform resource kinds (i.e. "google_compute_instance")
	// to their mapping/merging functions.
	mapperFuncs map[string][]mapper

	cfg *converter.Config

	// Map of converted assets (key = asset.Type + asset.Name)
	assets map[string]Asset
}

// Schemas exposes the schemas of resources this converter knows about.
func (c *Converter) Schemas() map[string]*schema.Resource {
	supported := make(map[string]*schema.Resource)
	for k := range c.schema.ResourcesMap {
		if _, ok := c.mapperFuncs[k]; ok {
			supported[k] = c.schema.ResourcesMap[k]
		}
	}
	return supported
}

// AddResource converts a terraform resource and stores the converted asset.
func (c *Converter) AddResource(r TerraformResource) error {
	for _, mapper := range c.mapperFuncs[r.Kind()] {
		data := struct {
			TerraformResource
			nonImplementedResourceData
		}{TerraformResource: r}

		converted, err := mapper.convert(data, c.cfg)
		if err != nil {
			if errors.Cause(err) == converter.ErrNoConversion {
				continue
			}
			return errors.Wrap(err, "converting asset")
		}

		key := converted.Type + converted.Name

		if existing, exists := c.assets[key]; exists {
			// The existance of a merge function signals that this tf resource maps to a
			// patching operation on an API resource.
			if mapper.merge != nil {
				converted = mapper.merge(existing, converted)
			} else {
				// If a merge function does not exist ignore the asset and return
				// a checkable error.
				return errors.Wrapf(ErrDuplicateAsset, "asset type %s: asset name %s",
					converted.Type, converted.Name)
			}
		}

		c.assets[key] = converted
	}

	return nil
}

type byName []Asset

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Assets lists all converted assets previously added by calls to AddResource.
func (c *Converter) Assets() []Asset {
	list := make([]Asset, 0, len(c.assets))
	for _, a := range c.assets {
		list = append(list, a)
	}
	sort.Sort(byName(list))
	return list
}
