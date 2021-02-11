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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/golang/glog"
	provider "github.com/hashicorp/terraform-provider-google/v3/google"
	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
	"github.com/GoogleCloudPlatform/terraform-validator/ancestrymanager"
)

var ErrDuplicateAsset = errors.New("duplicate asset")

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

	return &Converter{
		schema:          provider.Provider(),
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

// Compatibility shim: maintain support for ComposeTF12Resources -> AddResource pipeline
func (c *Converter) AddResource(rc *tfplan.ResourceChange) error {
	return c.AddResourceChanges([]tfplan.ResourceChange{*rc})
}

// AddResourceChange processes the resource changes in two stages:
// 1. Process deletions (fetching canonical resources from GCP as necessary)
// 2. Process creates and updates (fetching canonical resources from GCP as necessary)
// This will give us a deterministic end result even in cases where for example
// an IAM Binding and Member conflict with each other, but one is replacing the
// other.
func (c *Converter) AddResourceChanges(changes []tfplan.ResourceChange) error {
	var createOrUpdates []tfplan.ResourceChange
	for _, rc := range changes {
		// skip unknown resources
		if _, ok := c.schema.ResourcesMap[rc.Type]; !ok {
			glog.Infof("unknown resource: %s", rc.Type)
			continue
		}

		// Skip unsupported resources
		if _, ok := c.mapperFuncs[rc.Type]; !ok {
			glog.Infof("unsupported resource: %s", rc.Type)
			continue
		}

		if rc.IsCreate() || rc.IsUpdate() || rc.IsDeleteCreate() {
			createOrUpdates = append(createOrUpdates, rc)
		} else if rc.IsDelete() {
			if err := c.addDelete(rc); err != nil {
				return errors.Wrapf(err, "adding resource deletion")
			}
		}
	}

	for _, rc := range createOrUpdates {
		if err := c.addCreateOrUpdate(rc); err != nil {
			if errors.Cause(err) == ErrDuplicateAsset {
				glog.Warningf("adding resource change: %v", err)
			} else {
				return errors.Wrapf(err, "adding resource create or update")
			}
		}
	}

	return nil
}


// For deletions, we only need to handle mappers that support
// both fetch and mergeDelete. Supporting just one doesn't
// make sense, and supporting neither means that the deletion
// can just happen without needing to be merged.
func (c *Converter) addDelete(rc tfplan.ResourceChange) error {
	resource, _ := c.schema.ResourcesMap[rc.Type]
	rd := NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.Before,
	)
	for _, mapper := range c.mapperFuncs[rd.Kind()] {
		if mapper.fetch == nil || mapper.mergeDelete == nil {
			continue
		}
		converted, err := mapper.convert(&rd, c.cfg)
		if err != nil {
			if errors.Cause(err) == converter.ErrNoConversion {
				continue
			}
			return errors.Wrap(err, "converting asset")
		}

		key := converted.Type + converted.Name
		existing, exists := c.assets[key]
		var existingConverterAsset converter.Asset
		if !exists{
			existingConverterAsset, err = mapper.fetch(&rd, c.cfg)
			if err != nil {
				return errors.Wrap(err, "fetching asset")
			}
		} else {
			existingConverterAsset = existing.converterAsset
		}

		converted = mapper.mergeDelete(existingConverterAsset, converted)
		augmented, err := c.augmentAsset(&rd, c.cfg, converted)
		if err != nil {
			return errors.Wrap(err, "augmenting asset")
		}
		c.assets[key] = augmented
	}

	return nil
}


// For create/update, we need to handle both the case of no merging,
// and the case of merging. If merging, we expect both fetch and mergeCreateUpdate
// to be present.
func (c *Converter) addCreateOrUpdate(rc tfplan.ResourceChange) error {
	resource, _ := c.schema.ResourcesMap[rc.Type]
	rd := NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.After,
	)

	for _, mapper := range c.mapperFuncs[rd.Kind()] {
		converted, err := mapper.convert(&rd, c.cfg)
		if err != nil {
			if errors.Cause(err) == converter.ErrNoConversion {
				continue
			}
			return errors.Wrap(err, "converting asset")
		}

		key := converted.Type + converted.Name

		existing, exists := c.assets[key]
		var existingConverterAsset *converter.Asset
		if !exists && mapper.fetch != nil{
			asset, err := mapper.fetch(&rd, c.cfg)
			existingConverterAsset = &asset
			if err != nil {
				return errors.Wrap(err, "fetching asset")
			}
		} else if exists {
			existingConverterAsset = &existing.converterAsset
		}

		if existingConverterAsset != nil {
			if mapper.mergeCreateUpdate == nil {
				// If a merge function does not exist ignore the asset and return
				// a checkable error.
				return errors.Wrapf(ErrDuplicateAsset, "asset type %s: asset name %s",
					converted.Type, converted.Name)
			}
			converted = mapper.mergeCreateUpdate(*existingConverterAsset, converted)
		}

		augmented, err := c.augmentAsset(&rd, c.cfg, converted)
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
