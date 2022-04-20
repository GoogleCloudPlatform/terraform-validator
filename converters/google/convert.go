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

package google

import (
	errorssyslib "errors"
	"fmt"
	"sort"
	"strings"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	provider "github.com/hashicorp/terraform-provider-google/google"
	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/terraform-validator/ancestrymanager"
	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"go.uber.org/zap"
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
	converterAsset resources.Asset
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

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Unix(0, t.Nanos).UTC().Format(time.RFC3339Nano) + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	p, err := time.Parse(time.RFC3339Nano, strings.Trim(string(b), `"`))
	if err != nil {
		return fmt.Errorf("bad Timestamp: %v", err)
	}
	t.Seconds = p.Unix()
	t.Nanos = p.UnixNano()
	return nil
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
func NewConverter(cfg *resources.Config, ancestryManager ancestrymanager.AncestryManager, offline bool, convertUnchanged bool, errorLogger *zap.Logger) *Converter {
	return &Converter{
		schema:           provider.Provider(),
		converters:       resources.ResourceConverters(),
		offline:          offline,
		cfg:              cfg,
		ancestryManager:  ancestryManager,
		assets:           make(map[string]Asset),
		convertUnchanged: convertUnchanged,
		errorLogger:      errorLogger,
	}
}

// Converter knows how to convert terraform resources to their
// Google CAI (Cloud Asset Inventory) format (the Asset type).
type Converter struct {
	schema *schema.Provider

	// Map terraform resource kinds (i.e. "google_compute_instance")
	// to a ResourceConverter that can convert them to CAI assets.
	converters map[string][]resources.ResourceConverter

	offline bool
	cfg     *resources.Config

	// ancestryManager provides a manager to find the ancestry information for a project.
	ancestryManager ancestrymanager.AncestryManager

	// Map of converted assets (key = asset.Type + asset.Name)
	assets map[string]Asset

	// When set, Converter will convert ResourceChanges with no-op "actions".
	convertUnchanged bool

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

// AddResourceChange processes the resource changes in two stages:
// 1. Process deletions (fetching canonical resources from GCP as necessary)
// 2. Process creates, updates, and no-ops (fetching canonical resources from GCP as necessary)
// This will give us a deterministic end result even in cases where for example
// an IAM Binding and Member conflict with each other, but one is replacing the
// other.
func (c *Converter) AddResourceChanges(changes []*tfjson.ResourceChange) error {
	var createOrUpdateOrNoops []*tfjson.ResourceChange
	for _, rc := range changes {
		// skip unknown resources
		if _, ok := c.schema.ResourcesMap[rc.Type]; !ok {
			c.errorLogger.Info(fmt.Sprintf("unknown resource: %s", rc.Type))
			continue
		}

		// Skip unsupported resources
		if _, ok := c.converters[rc.Type]; !ok {
			c.errorLogger.Info(fmt.Sprintf("unsupported resource: %s", rc.Type))
			continue
		}

		if tfplan.IsCreate(rc) || tfplan.IsUpdate(rc) || tfplan.IsDeleteCreate(rc) || (c.convertUnchanged && tfplan.IsNoOp(rc)) {
			createOrUpdateOrNoops = append(createOrUpdateOrNoops, rc)
		} else if tfplan.IsDelete(rc) {
			if err := c.addDelete(rc); err != nil {
				return fmt.Errorf("adding resource deletion %w", err)
			}
		}
	}

	for _, rc := range createOrUpdateOrNoops {
		if err := c.addCreateOrUpdateOrNoop(rc); err != nil {
			if errorssyslib.Is(err, ErrDuplicateAsset) {
				c.errorLogger.Warn(fmt.Sprintf("adding resource change: %v", err))
			} else {
				return fmt.Errorf("adding resource create/update/no-op %w", err)
			}
		}
	}

	return nil
}

// For deletions, we only need to handle ResourceConverters that support
// both fetch and mergeDelete. Supporting just one doesn't
// make sense, and supporting neither means that the deletion
// can just happen without needing to be merged.
func (c *Converter) addDelete(rc *tfjson.ResourceChange) error {
	resource, _ := c.schema.ResourcesMap[rc.Type]
	rd := NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.Before.(map[string]interface{}),
	)
	for _, converter := range c.converters[rd.Kind()] {
		if converter.FetchFullResource == nil || converter.MergeDelete == nil {
			continue
		}
		convertedItems, err := converter.Convert(&rd, c.cfg)

		if err != nil {
			if errors.Cause(err) == resources.ErrNoConversion {
				continue
			}
			return errors.Wrap(err, "converting asset")
		}

		for _, converted := range convertedItems {

			key := converted.Type + converted.Name
			var existingConverterAsset *resources.Asset
			if existing, exists := c.assets[key]; exists {
				existingConverterAsset = &existing.converterAsset
			} else if !c.offline {
				asset, err := converter.FetchFullResource(&rd, c.cfg)
				if errors.Cause(err) == resources.ErrEmptyIdentityField {
					c.errorLogger.Warn(fmt.Sprintf("%s did not return a value for ID field. Skipping asset fetch.", key))
					existingConverterAsset = nil
				} else if errors.Cause(err) == resources.ErrResourceInaccessible {
					c.errorLogger.Warn(fmt.Sprintf("%s was not able to be fetched due to not existing or insufficient permission. Skipping asset fetch.", key))
					existingConverterAsset = nil
				} else if err != nil {
					return errors.Wrap(err, "fetching asset")
				} else {
					existingConverterAsset = &asset
				}
				if existingConverterAsset != nil {
					converted = converter.MergeDelete(*existingConverterAsset, converted)
					augmented, err := c.augmentAsset(&rd, c.cfg, converted)
					if err != nil {
						return errors.Wrap(err, "augmenting asset")
					}
					c.assets[key] = augmented
				}
			}
		}
	}

	return nil
}

// For create/update/no-op, we need to handle both the case of no merging,
// and the case of merging. If merging, we expect both fetch and mergeCreateUpdate
// to be present.
func (c *Converter) addCreateOrUpdateOrNoop(rc *tfjson.ResourceChange) error {
	resource, _ := c.schema.ResourcesMap[rc.Type]
	rd := NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.After.(map[string]interface{}),
	)

	for _, converter := range c.converters[rd.Kind()] {
		convertedAssets, err := converter.Convert(&rd, c.cfg)
		if err != nil {
			if errors.Cause(err) == resources.ErrNoConversion {
				continue
			}
			return errors.Wrap(err, "converting asset")
		}

		for _, converted := range convertedAssets {
			key := converted.Type + converted.Name

			var existingConverterAsset *resources.Asset
			if existing, exists := c.assets[key]; exists {
				existingConverterAsset = &existing.converterAsset
			} else if converter.FetchFullResource != nil && !c.offline {
				asset, err := converter.FetchFullResource(&rd, c.cfg)
				if errors.Cause(err) == resources.ErrEmptyIdentityField {
					c.errorLogger.Warn(fmt.Sprintf("%s did not return a value for ID field. Skipping asset fetch.", key))
					existingConverterAsset = nil
				} else if errors.Cause(err) == resources.ErrResourceInaccessible {
					c.errorLogger.Warn(fmt.Sprintf("%s was not able to be fetched due to not existing or insufficient permission. Skipping asset fetch.", key))
					existingConverterAsset = nil
				} else if err != nil {
					return errors.Wrap(err, "fetching asset")
				} else {
					existingConverterAsset = &asset
				}
			}

			if existingConverterAsset != nil {
				if converter.MergeCreateUpdate == nil {
					// If a merge function does not exist ignore the asset and return
					// a checkable error.
					return fmt.Errorf("asset type %s: asset name %s %w", converted.Type, converted.Name, ErrDuplicateAsset)
				}
				converted = converter.MergeCreateUpdate(*existingConverterAsset, converted)
			}

			augmented, err := c.augmentAsset(&rd, c.cfg, converted)
			if err != nil {
				return errors.Wrap(err, "augmenting asset")
			}
			c.assets[key] = augmented
		}
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
func (c *Converter) augmentAsset(tfData resources.TerraformResourceData, cfg *resources.Config, cai resources.Asset) (Asset, error) {
	project, err := getProject(tfData, cfg, cai, c.errorLogger)
	if err != nil {
		return Asset{}, fmt.Errorf("getting project for %v: %w", cai.Name, err)
	}
	ancestry, err := c.ancestryManager.GetAncestryWithResource(project, tfData, cai)
	if err != nil {
		return Asset{}, fmt.Errorf("getting resource ancestry for project %v: %w", project, err)
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
			//As time is not information in terraform resource data, time is fixed for testing purposes
			fixedTime := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
			orgPolicy = append(orgPolicy, &OrgPolicy{
				Constraint:     o.Constraint,
				ListPolicy:     listPolicy,
				BooleanPolicy:  booleanPolicy,
				RestoreDefault: restoreDefault,
				UpdateTime: &Timestamp{
					Seconds: int64(fixedTime.Unix()),
					Nanos:   int64(fixedTime.UnixNano()),
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
