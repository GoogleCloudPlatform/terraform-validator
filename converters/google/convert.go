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
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/terraform-provider-google/helper/schema"
	"github.com/pkg/errors"
	provider "github.com/terraform-providers/terraform-provider-google/google"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
	"github.com/GoogleCloudPlatform/terraform-validator/ancestrymanager"
)

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

// Asset contains the resource data and metadata in the same format as
// Google CAI (Cloud Asset Inventory).
type Asset struct {
	Name      string         `json:"name"`
	Type      string         `json:"asset_type"`
	Ancestry  string         `json:"ancestry_path"`
	Resource  *AssetResource `json:"resource,omitempty"`
	IAMPolicy *IAMPolicy     `json:"iam_policy,omitempty"`
	OrgPolicy []*OrgPolicy   `json:"org_policy,omitempty"`
	// Store the converter's version of the asset to allow for merges which
	// operate on this type. When matching json tags land in the conversions
	// library, this could be nested to avoid the duplication of fields.
	converterAsset converter.Asset
}

// IAMPolicy is the representation of a Cloud IAM policy set on a cloud resource.
type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
}

// IAMBinding binds a role to a set of members.
type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

// AssetResource is nested within the Asset type.
type AssetResource struct {
	Version              string                 `json:"version"`
	DiscoveryDocumentURI string                 `json:"discovery_document_uri"`
	DiscoveryName        string                 `json:"discovery_name"`
	Parent               string                 `json:"parent"`
	Data                 map[string]interface{} `json:"data"`
}

//OrgPolicy is for managing organization policies.
type OrgPolicy struct {
	Constraint     string          `json:"constraint,omitempty"`
	ListPolicy     *ListPolicy     `json:"list_policy,omitempty"`
	BooleanPolicy  *BooleanPolicy  `json:"boolean_policy,omitempty"`
	RestoreDefault *RestoreDefault `json:"restore_default,omitempty"`
	UpdateTime     *Timestamp      `json:"update_time,omitempty"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nanos   int64 `json:"nanos,omitempty"`
}

// ListPolicyAllValues is used to set `Policies` that apply to all possible
// configuration values rather than specific values in `allowed_values` or
// `denied_values`.
type ListPolicyAllValues int32

// ListPolicy can define specific values and subtrees of Cloud Resource
// Manager resource hierarchy (`Organizations`, `Folders`, `Projects`) that
// are allowed or denied by setting the `allowed_values` and `denied_values`
// fields.
type ListPolicy struct {
	AllowedValues     []string            `json:"allowed_values,omitempty"`
	DeniedValues      []string            `json:"denied_values,omitempty"`
	AllValues         ListPolicyAllValues `json:"all_values,omitempty"`
	SuggestedValue    string              `json:"suggested_value,omitempty"`
	InheritFromParent bool                `json:"inherit_from_parent,omitempty"`
}

// BooleanPolicy If `true`, then the `Policy` is enforced. If `false`,
// then any configuration is acceptable.
type BooleanPolicy struct {
	Enforced bool `json:"enforced,omitempty"`
}

//RestoreDefault determines if the default values of the `Constraints` are active for the
// resources.
type RestoreDefault struct {
}

// NewConverter is a factory function for Converter.
func NewConverter(ctx context.Context, ancestryManager ancestrymanager.AncestryManager, project, credentials string, offline bool) (*Converter, error) {
	cfg := &converter.Config{
		Project:     project,
		Credentials: credentials,
	}
	if !offline {
		converter.ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
	}

	p := provider.Provider().(*schema.Provider)
	return &Converter{
		schema:          p,
		mapperFuncs:     mappers(),
		cfg:             cfg,
		ancestryManager: ancestryManager,
		assets:          make(map[string]Asset),
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

	// ancestryManager provides a manager to find the ancestry information for a project.
	ancestryManager ancestrymanager.AncestryManager

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
				converted = mapper.merge(existing.converterAsset, converted)
			} else {
				// If a merge function does not exist ignore the asset and return
				// a checkable error.
				return errors.Wrapf(ErrDuplicateAsset, "asset type %s: asset name %s",
					converted.Type, converted.Name)
			}
		}

		augmented, err := c.augmentAsset(data, c.cfg, converted)
		if err != nil {
			return errors.Wrap(err, "augmenting asset")
		}
		c.assets[key] = augmented
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

// augmentAsset adds data to an asset that is not set by the conversion library.
func (c *Converter) augmentAsset(tfData converter.TerraformResourceData, cfg *converter.Config, cai converter.Asset) (Asset, error) {
	project, err := getProject(tfData, cfg, cai)
	if err != nil {
		return Asset{}, err
	}
	ancestry, err := c.ancestryManager.GetAncestryWithResource(project, tfData, cai)
	if err != nil {
		return Asset{}, errors.Wrapf(err, "getting resource ancestry: project %v", project)
	}
	var resource *AssetResource
	if cai.Resource != nil {
		resource = &AssetResource{
			Version:              cai.Resource.Version,
			DiscoveryDocumentURI: cai.Resource.DiscoveryDocumentURI,
			DiscoveryName:        cai.Resource.DiscoveryName,
			Parent:               fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%v", project),
			Data:                 cai.Resource.Data,
		}
	}

	var policy *IAMPolicy
	if cai.IAMPolicy != nil {
		policy = &IAMPolicy{}
		for _, b := range cai.IAMPolicy.Bindings {
			policy.Bindings = append(policy.Bindings, IAMBinding{
				Role:    b.Role,
				Members: b.Members,
			})
		}
	}

	var orgPolicy []*OrgPolicy
	if cai.OrgPolicy != nil {
		for _, o := range cai.OrgPolicy {
			var listPolicy *ListPolicy
			var booleanPolicy *BooleanPolicy
			var restoreDefault *RestoreDefault
			if o.ListPolicy != nil {
				listPolicy = &ListPolicy{
					AllowedValues:     o.ListPolicy.AllowedValues,
					AllValues:         ListPolicyAllValues(o.ListPolicy.AllValues),
					DeniedValues:      o.ListPolicy.DeniedValues,
					SuggestedValue:    o.ListPolicy.SuggestedValue,
					InheritFromParent: o.ListPolicy.InheritFromParent,
				}
			}
			if o.BooleanPolicy != nil {
				booleanPolicy = &BooleanPolicy{
					Enforced: o.BooleanPolicy.Enforced,
				}
			}
			if o.RestoreDefault != nil {
				restoreDefault = &RestoreDefault{}
			}
			//As time is not information in terraform resource data, time is rounded for testing purposes
			currentTime := time.Now()
			currentTime = currentTime.Round(time.Minute)
			orgPolicy = append(orgPolicy, &OrgPolicy{
				Constraint:     o.Constraint,
				ListPolicy:     listPolicy,
				BooleanPolicy:  booleanPolicy,
				RestoreDefault: restoreDefault,
				UpdateTime: &Timestamp{
					Seconds: int64(currentTime.Unix()),
					Nanos:   int64(currentTime.UnixNano()),
				},
			})
		}
	}

	return Asset{
		Name:           cai.Name,
		Type:           cai.Type,
		Ancestry:       ancestry,
		Resource:       resource,
		IAMPolicy:      policy,
		OrgPolicy:      orgPolicy,
		converterAsset: cai,
	}, nil
}
