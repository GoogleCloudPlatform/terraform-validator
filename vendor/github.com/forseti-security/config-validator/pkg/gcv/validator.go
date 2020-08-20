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

// Package gcv provides a library and a RPC service for Forseti Config Validator.
package gcv

import (
	"context"
	"encoding/json"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	asset2 "github.com/forseti-security/config-validator/pkg/asset"
	"github.com/forseti-security/config-validator/pkg/gcptarget"
	"github.com/forseti-security/config-validator/pkg/gcv/configs"
	"github.com/forseti-security/config-validator/pkg/multierror"
	"github.com/golang/glog"
	cfclient "github.com/open-policy-agent/frameworks/constraint/pkg/client"
	"github.com/open-policy-agent/frameworks/constraint/pkg/client/drivers/local"
	cftemplates "github.com/open-policy-agent/frameworks/constraint/pkg/core/templates"
	k8starget "github.com/open-policy-agent/gatekeeper/pkg/target"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	logRequestsVerboseLevel = 2
	// The JSON object key for ancestry path
	ancestryPathKey = "ancestry_path"
	// The JSON object key for ancestors list
	ancestorSliceKey = "ancestors"
)

type ConfigValidator interface {
	ReviewAsset(ctx context.Context, asset *validator.Asset) ([]*validator.Violation, error)
}

// Validator checks GCP resource metadata for constraint violation.
//
// Expected usage pattern:
//   - call NewValidator to create a new Validator
//   - call AddData one or more times to add the GCP resource metadata to check
//   - call Audit to validate the GCP resource metadata that has been added so far
//   - call Reset to delete existing data
//   - call AddData to add a new set of GCP resource metadata to check
//   - call Reset to delete existing data
//
// Any data added in AddData stays in the underlying rule evaluation engine's memory.
// To avoid out of memory errors, callers can invoke Reset to delete existing data.
type Validator struct {
	// policyPaths is a list of paths where the constraints and constraint templates are stored as yaml files.
	// Each path can refer to a directory or file.
	policyPaths []string
	// policy dependencies directory points to rego files that provide supporting code for templates.
	// These rego dependencies should be packaged with the GCV deployment.
	// Right now expected to be set to point to "//policies/validator/lib" folder
	policyLibraryDir string
	gcpCFClient      *cfclient.Client
	k8sCFClient      *cfclient.Client
}

// NewValidatorConfig returns a new ValidatorConfig.
// By default it will initialize the underlying query evaluation engine by loading supporting library, constraints, and constraint templates.
// We may want to make this initialization behavior configurable in the future.
func NewValidatorConfig(policyPaths []string, policyLibraryPath string) (*configs.Configuration, error) {
	if len(policyPaths) == 0 {
		return nil, errors.Errorf("No policy path set, provide an option to set the policy path gcv.PolicyPath")
	}
	if policyLibraryPath == "" {
		return nil, errors.Errorf("No policy library set")
	}
	glog.V(logRequestsVerboseLevel).Infof("loading policy dir: %v lib dir: %s", policyPaths, policyLibraryPath)
	return configs.NewConfiguration(policyPaths, policyLibraryPath)
}

func newCFClient(
	targetHandler cfclient.TargetHandler,
	templates []*cftemplates.ConstraintTemplate,
	constraints []*unstructured.Unstructured) (
	*cfclient.Client, error) {
	driver := local.New(local.Tracing(false))
	backend, err := cfclient.NewBackend(cfclient.Driver(driver))
	if err != nil {
		return nil, errors.Wrap(err, "unable to set up Constraint Framework backend")
	}
	cfClient, err := backend.NewClient(cfclient.Targets(targetHandler))
	if err != nil {
		return nil, errors.Wrap(err, "unable to set up Constraint Framework client")
	}

	ctx := context.Background()
	var errs multierror.Errors
	for _, template := range templates {
		if _, err := cfClient.AddTemplate(ctx, template); err != nil {
			errs.Add(errors.Wrapf(err, "failed to add template %v", template))
		}
	}
	if !errs.Empty() {
		return nil, errs.ToError()
	}

	for _, constraint := range constraints {
		if _, err := cfClient.AddConstraint(ctx, constraint); err != nil {
			errs.Add(errors.Wrapf(err, "failed to add constraint %s", constraint))
		}
	}
	if !errs.Empty() {
		return nil, errs.ToError()
	}
	return cfClient, nil
}

// NewValidatorFromConfig creates the validator from a config.
func NewValidatorFromConfig(config *configs.Configuration) (*Validator, error) {
	gcpCFClient, err := newCFClient(gcptarget.New(), config.GCPTemplates, config.GCPConstraints)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set up GCP Constraint Framework client")
	}

	k8sCFClient, err := newCFClient(&k8starget.K8sValidationTarget{}, config.K8STemplates, config.K8SConstraints)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set up K8S Constraint Framework client")
	}

	ret := &Validator{
		gcpCFClient: gcpCFClient,
		k8sCFClient: k8sCFClient,
	}
	return ret, nil
}

// NewValidator returns a new Validator.
// By default it will initialize the underlying query evaluation engine by loading supporting library, constraints, and constraint templates.
// We may want to make this initialization behavior configurable in the future.
func NewValidator(policyPaths []string, policyLibraryPath string) (*Validator, error) {
	config, err := NewValidatorConfig(policyPaths, policyLibraryPath)
	if err != nil {
		return nil, err
	}
	return NewValidatorFromConfig(config)
}

// ReviewAsset reviews a single asset.
func (v *Validator) ReviewAsset(ctx context.Context, asset *validator.Asset) ([]*validator.Violation, error) {
	if err := asset2.ValidateAsset(asset); err != nil {
		return nil, err
	}

	if err := asset2.SanitizeAncestryPath(asset); err != nil {
		return nil, err
	}

	assetInterface, err := asset2.ConvertResourceViaJSONToInterface(asset)
	if err != nil {
		return nil, err
	}

	assetMapInterface := assetInterface.(map[string]interface{})
	result, err := v.ReviewUnmarshalledJSON(ctx, assetMapInterface)
	if err != nil {
		return nil, err
	}

	return result.ToViolations()
}

// fixAncestry will try to use the ancestors array to create the ancestorPath
// value if it is not present.
func (v *Validator) fixAncestry(input map[string]interface{}) error {
	ancestors, found, err := unstructured.NestedStringSlice(input, ancestorSliceKey)
	if found && err == nil {
		input[ancestryPathKey] = asset2.AncestryPath(ancestors)
		return nil
	}

	ancestry, found, err := unstructured.NestedString(input, ancestryPathKey)
	if found && err == nil {
		input[ancestryPathKey] = configs.NormalizeAncestry(ancestry)
		return nil
	}
	return errors.Errorf("asset missing ancestry information: %v", input)
}

// ReviewJSON reviews the content of a JSON string
func (v *Validator) ReviewJSON(ctx context.Context, data string) (*Result, error) {
	asset := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &asset); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal json")
	}
	return v.ReviewUnmarshalledJSON(ctx, asset)
}

// ReviewJSON evaluates a single asset without any threading in the background.
func (v *Validator) ReviewUnmarshalledJSON(ctx context.Context, asset map[string]interface{}) (*Result, error) {
	if err := v.fixAncestry(asset); err != nil {
		return nil, err
	}

	if asset2.IsK8S(asset) {
		return v.reviewK8SResource(ctx, asset)
	}
	return v.reviewGCPResource(ctx, asset)
}

// reviewK8SResource will unwrap k8s resources then pass them to the cf client with the gatekeeper target.
func (v *Validator) reviewK8SResource(ctx context.Context, asset map[string]interface{}) (*Result, error) {
	k8sResource, err := asset2.UnwrapCAIResource(asset)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert asset to admission request")
	}
	responses, err := v.k8sCFClient.Review(ctx, k8sResource)
	if err != nil {
		return nil, errors.Wrapf(err, "K8S target Constraint Framework review call failed")
	}
	return NewResult(configs.K8STargetName, asset, k8sResource.Object, responses)
}

// reviewGCPResource will unwrap k8s resources then pass them to the cf client with the gatekeeper target.
func (v *Validator) reviewGCPResource(ctx context.Context, asset map[string]interface{}) (*Result, error) {
	responses, err := v.gcpCFClient.Review(ctx, asset)
	if err != nil {
		return nil, errors.Wrapf(err, "GCP target Constraint Framework review call failed")
	}
	return NewResult(gcptarget.Name, asset, asset, responses)
}
