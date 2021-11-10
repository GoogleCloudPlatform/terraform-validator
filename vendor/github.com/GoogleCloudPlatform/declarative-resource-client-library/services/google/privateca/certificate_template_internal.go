// Copyright 2021 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package privateca

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *CertificateTemplate) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.PredefinedValues) {
		if err := r.PredefinedValues.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.IdentityConstraints) {
		if err := r.IdentityConstraints.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PassthroughExtensions) {
		if err := r.PassthroughExtensions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateTemplatePredefinedValues) validate() error {
	if !dcl.IsEmptyValueIndirect(r.KeyUsage) {
		if err := r.KeyUsage.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.CaOptions) {
		if err := r.CaOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateTemplatePredefinedValuesKeyUsage) validate() error {
	if !dcl.IsEmptyValueIndirect(r.BaseKeyUsage) {
		if err := r.BaseKeyUsage.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ExtendedKeyUsage) {
		if err := r.ExtendedKeyUsage.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateTemplatePredefinedValuesCaOptions) validate() error {
	return nil
}
func (r *CertificateTemplatePredefinedValuesPolicyIds) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateTemplatePredefinedValuesAdditionalExtensions) validate() error {
	if err := dcl.Required(r, "objectId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "value"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.ObjectId) {
		if err := r.ObjectId.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateTemplateIdentityConstraints) validate() error {
	if err := dcl.Required(r, "allowSubjectPassthrough"); err != nil {
		return err
	}
	if err := dcl.Required(r, "allowSubjectAltNamesPassthrough"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.CelExpression) {
		if err := r.CelExpression.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateTemplateIdentityConstraintsCelExpression) validate() error {
	return nil
}
func (r *CertificateTemplatePassthroughExtensions) validate() error {
	return nil
}
func (r *CertificateTemplatePassthroughExtensionsAdditionalExtensions) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateTemplate) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://privateca.googleapis.com/v1/", params)
}

func (r *CertificateTemplate) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *CertificateTemplate) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/certificateTemplates", nr.basePath(), userBasePath, params), nil

}

func (r *CertificateTemplate) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/certificateTemplates?certificateTemplateId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *CertificateTemplate) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}", nr.basePath(), userBasePath, params), nil
}

// certificateTemplateApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type certificateTemplateApiOperation interface {
	do(context.Context, *CertificateTemplate, *Client) error
}

// newUpdateCertificateTemplateUpdateCertificateTemplateRequest creates a request for an
// CertificateTemplate resource's UpdateCertificateTemplate update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateCertificateTemplateUpdateCertificateTemplateRequest(ctx context.Context, f *CertificateTemplate, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v, err := expandCertificateTemplatePredefinedValues(c, f.PredefinedValues); err != nil {
		return nil, fmt.Errorf("error expanding PredefinedValues into predefinedValues: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["predefinedValues"] = v
	}
	if v, err := expandCertificateTemplateIdentityConstraints(c, f.IdentityConstraints); err != nil {
		return nil, fmt.Errorf("error expanding IdentityConstraints into identityConstraints: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["identityConstraints"] = v
	}
	if v, err := expandCertificateTemplatePassthroughExtensions(c, f.PassthroughExtensions); err != nil {
		return nil, fmt.Errorf("error expanding PassthroughExtensions into passthroughExtensions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["passthroughExtensions"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/certificateTemplates/%s", *f.Project, *f.Location, *f.Name)

	return req, nil
}

// marshalUpdateCertificateTemplateUpdateCertificateTemplateRequest converts the update into
// the final JSON request body.
func marshalUpdateCertificateTemplateUpdateCertificateTemplateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateCertificateTemplateUpdateCertificateTemplateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateCertificateTemplateUpdateCertificateTemplateOperation) do(ctx context.Context, r *CertificateTemplate, c *Client) error {
	_, err := c.GetCertificateTemplate(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateCertificateTemplate")
	if err != nil {
		return err
	}
	mask := dcl.TopLevelUpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateCertificateTemplateUpdateCertificateTemplateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateCertificateTemplateUpdateCertificateTemplateRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listCertificateTemplateRaw(ctx context.Context, r *CertificateTemplate, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != CertificateTemplateMaxPage {
		m["pageSize"] = fmt.Sprintf("%v", pageSize)
	}

	u, err = dcl.AddQueryParams(u, m)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	return ioutil.ReadAll(resp.Response.Body)
}

type listCertificateTemplateOperation struct {
	CertificateTemplates []map[string]interface{} `json:"certificateTemplates"`
	Token                string                   `json:"nextPageToken"`
}

func (c *Client) listCertificateTemplate(ctx context.Context, r *CertificateTemplate, pageToken string, pageSize int32) ([]*CertificateTemplate, string, error) {
	b, err := c.listCertificateTemplateRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listCertificateTemplateOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*CertificateTemplate
	for _, v := range m.CertificateTemplates {
		res, err := unmarshalMapCertificateTemplate(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllCertificateTemplate(ctx context.Context, f func(*CertificateTemplate) bool, resources []*CertificateTemplate) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteCertificateTemplate(ctx, res)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", strings.Join(errors, "\n"))
	} else {
		return nil
	}
}

type deleteCertificateTemplateOperation struct{}

func (op *deleteCertificateTemplateOperation) do(ctx context.Context, r *CertificateTemplate, c *Client) error {
	r, err := c.GetCertificateTemplate(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "CertificateTemplate not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetCertificateTemplate checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		return err
	}

	// we saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// this is the reason we are adding retry to handle that case.
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		_, err = c.GetCertificateTemplate(ctx, r)
		if !dcl.IsNotFound(err) {
			if i == maxRetry {
				return dcl.NotDeletedError{ExistingResource: r}
			}
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createCertificateTemplateOperation struct {
	response map[string]interface{}
}

func (op *createCertificateTemplateOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createCertificateTemplateOperation) do(ctx context.Context, r *CertificateTemplate, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(req), c.Config.RetryProvider)
	if err != nil {
		return err
	}
	// wait for object to be created.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetCertificateTemplate(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getCertificateTemplateRaw(ctx context.Context, r *CertificateTemplate) ([]byte, error) {

	u, err := r.getURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	b, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) certificateTemplateDiffsForRawDesired(ctx context.Context, rawDesired *CertificateTemplate, opts ...dcl.ApplyOption) (initial, desired *CertificateTemplate, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *CertificateTemplate
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*CertificateTemplate); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected CertificateTemplate, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetCertificateTemplate(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a CertificateTemplate resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve CertificateTemplate resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that CertificateTemplate resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeCertificateTemplateDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for CertificateTemplate: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for CertificateTemplate: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeCertificateTemplateInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for CertificateTemplate: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeCertificateTemplateDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for CertificateTemplate: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffCertificateTemplate(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeCertificateTemplateInitialState(rawInitial, rawDesired *CertificateTemplate) (*CertificateTemplate, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeCertificateTemplateDesiredState(rawDesired, rawInitial *CertificateTemplate, opts ...dcl.ApplyOption) (*CertificateTemplate, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.PredefinedValues = canonicalizeCertificateTemplatePredefinedValues(rawDesired.PredefinedValues, nil, opts...)
		rawDesired.IdentityConstraints = canonicalizeCertificateTemplateIdentityConstraints(rawDesired.IdentityConstraints, nil, opts...)
		rawDesired.PassthroughExtensions = canonicalizeCertificateTemplatePassthroughExtensions(rawDesired.PassthroughExtensions, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &CertificateTemplate{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.PredefinedValues = canonicalizeCertificateTemplatePredefinedValues(rawDesired.PredefinedValues, rawInitial.PredefinedValues, opts...)
	canonicalDesired.IdentityConstraints = canonicalizeCertificateTemplateIdentityConstraints(rawDesired.IdentityConstraints, rawInitial.IdentityConstraints, opts...)
	canonicalDesired.PassthroughExtensions = canonicalizeCertificateTemplatePassthroughExtensions(rawDesired.PassthroughExtensions, rawInitial.PassthroughExtensions, opts...)
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.Labels) {
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}

	return canonicalDesired, nil
}

func canonicalizeCertificateTemplateNewState(c *Client, rawNew, rawDesired *CertificateTemplate) (*CertificateTemplate, error) {

	if dcl.IsNotReturnedByServer(rawNew.Name) && dcl.IsNotReturnedByServer(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.PredefinedValues) && dcl.IsNotReturnedByServer(rawDesired.PredefinedValues) {
		rawNew.PredefinedValues = rawDesired.PredefinedValues
	} else {
		rawNew.PredefinedValues = canonicalizeNewCertificateTemplatePredefinedValues(c, rawDesired.PredefinedValues, rawNew.PredefinedValues)
	}

	if dcl.IsNotReturnedByServer(rawNew.IdentityConstraints) && dcl.IsNotReturnedByServer(rawDesired.IdentityConstraints) {
		rawNew.IdentityConstraints = rawDesired.IdentityConstraints
	} else {
		rawNew.IdentityConstraints = canonicalizeNewCertificateTemplateIdentityConstraints(c, rawDesired.IdentityConstraints, rawNew.IdentityConstraints)
	}

	if dcl.IsNotReturnedByServer(rawNew.PassthroughExtensions) && dcl.IsNotReturnedByServer(rawDesired.PassthroughExtensions) {
		rawNew.PassthroughExtensions = rawDesired.PassthroughExtensions
	} else {
		rawNew.PassthroughExtensions = canonicalizeNewCertificateTemplatePassthroughExtensions(c, rawDesired.PassthroughExtensions, rawNew.PassthroughExtensions)
	}

	if dcl.IsNotReturnedByServer(rawNew.Description) && dcl.IsNotReturnedByServer(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.CreateTime) && dcl.IsNotReturnedByServer(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.UpdateTime) && dcl.IsNotReturnedByServer(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Labels) && dcl.IsNotReturnedByServer(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeCertificateTemplatePredefinedValues(des, initial *CertificateTemplatePredefinedValues, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValues {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValues{}

	cDes.KeyUsage = canonicalizeCertificateTemplatePredefinedValuesKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCertificateTemplatePredefinedValuesCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCertificateTemplatePredefinedValuesPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	if dcl.IsZeroValue(des.AiaOcspServers) {
		des.AiaOcspServers = initial.AiaOcspServers
	} else {
		cDes.AiaOcspServers = des.AiaOcspServers
	}
	cDes.AdditionalExtensions = canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesSlice(des, initial []CertificateTemplatePredefinedValues, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValues {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValues, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValues(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValues, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValues(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValues(c *Client, des, nw *CertificateTemplatePredefinedValues) *CertificateTemplatePredefinedValues {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValues while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCertificateTemplatePredefinedValuesKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCertificateTemplatePredefinedValuesCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCertificateTemplatePredefinedValuesPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	nw.AdditionalExtensions = canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesSet(c *Client, des, nw []CertificateTemplatePredefinedValues) []CertificateTemplatePredefinedValues {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValues
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesSlice(c *Client, des, nw []CertificateTemplatePredefinedValues) []CertificateTemplatePredefinedValues {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValues
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValues(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsage(des, initial *CertificateTemplatePredefinedValuesKeyUsage, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageSlice(des, initial []CertificateTemplatePredefinedValuesKeyUsage, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesKeyUsage {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsage(c *Client, des, nw *CertificateTemplatePredefinedValuesKeyUsage) *CertificateTemplatePredefinedValuesKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageSet(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsage) []CertificateTemplatePredefinedValuesKeyUsage {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesKeyUsage
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsage) []CertificateTemplatePredefinedValuesKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(des, initial *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}

	if dcl.BoolCanonicalize(des.DigitalSignature, initial.DigitalSignature) || dcl.IsZeroValue(des.DigitalSignature) {
		cDes.DigitalSignature = initial.DigitalSignature
	} else {
		cDes.DigitalSignature = des.DigitalSignature
	}
	if dcl.BoolCanonicalize(des.ContentCommitment, initial.ContentCommitment) || dcl.IsZeroValue(des.ContentCommitment) {
		cDes.ContentCommitment = initial.ContentCommitment
	} else {
		cDes.ContentCommitment = des.ContentCommitment
	}
	if dcl.BoolCanonicalize(des.KeyEncipherment, initial.KeyEncipherment) || dcl.IsZeroValue(des.KeyEncipherment) {
		cDes.KeyEncipherment = initial.KeyEncipherment
	} else {
		cDes.KeyEncipherment = des.KeyEncipherment
	}
	if dcl.BoolCanonicalize(des.DataEncipherment, initial.DataEncipherment) || dcl.IsZeroValue(des.DataEncipherment) {
		cDes.DataEncipherment = initial.DataEncipherment
	} else {
		cDes.DataEncipherment = des.DataEncipherment
	}
	if dcl.BoolCanonicalize(des.KeyAgreement, initial.KeyAgreement) || dcl.IsZeroValue(des.KeyAgreement) {
		cDes.KeyAgreement = initial.KeyAgreement
	} else {
		cDes.KeyAgreement = des.KeyAgreement
	}
	if dcl.BoolCanonicalize(des.CertSign, initial.CertSign) || dcl.IsZeroValue(des.CertSign) {
		cDes.CertSign = initial.CertSign
	} else {
		cDes.CertSign = des.CertSign
	}
	if dcl.BoolCanonicalize(des.CrlSign, initial.CrlSign) || dcl.IsZeroValue(des.CrlSign) {
		cDes.CrlSign = initial.CrlSign
	} else {
		cDes.CrlSign = des.CrlSign
	}
	if dcl.BoolCanonicalize(des.EncipherOnly, initial.EncipherOnly) || dcl.IsZeroValue(des.EncipherOnly) {
		cDes.EncipherOnly = initial.EncipherOnly
	} else {
		cDes.EncipherOnly = des.EncipherOnly
	}
	if dcl.BoolCanonicalize(des.DecipherOnly, initial.DecipherOnly) || dcl.IsZeroValue(des.DecipherOnly) {
		cDes.DecipherOnly = initial.DecipherOnly
	} else {
		cDes.DecipherOnly = des.DecipherOnly
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice(des, initial []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c *Client, des, nw *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.DigitalSignature, nw.DigitalSignature) {
		nw.DigitalSignature = des.DigitalSignature
	}
	if dcl.BoolCanonicalize(des.ContentCommitment, nw.ContentCommitment) {
		nw.ContentCommitment = des.ContentCommitment
	}
	if dcl.BoolCanonicalize(des.KeyEncipherment, nw.KeyEncipherment) {
		nw.KeyEncipherment = des.KeyEncipherment
	}
	if dcl.BoolCanonicalize(des.DataEncipherment, nw.DataEncipherment) {
		nw.DataEncipherment = des.DataEncipherment
	}
	if dcl.BoolCanonicalize(des.KeyAgreement, nw.KeyAgreement) {
		nw.KeyAgreement = des.KeyAgreement
	}
	if dcl.BoolCanonicalize(des.CertSign, nw.CertSign) {
		nw.CertSign = des.CertSign
	}
	if dcl.BoolCanonicalize(des.CrlSign, nw.CrlSign) {
		nw.CrlSign = des.CrlSign
	}
	if dcl.BoolCanonicalize(des.EncipherOnly, nw.EncipherOnly) {
		nw.EncipherOnly = des.EncipherOnly
	}
	if dcl.BoolCanonicalize(des.DecipherOnly, nw.DecipherOnly) {
		nw.DecipherOnly = des.DecipherOnly
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSet(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(des, initial *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}

	if dcl.BoolCanonicalize(des.ServerAuth, initial.ServerAuth) || dcl.IsZeroValue(des.ServerAuth) {
		cDes.ServerAuth = initial.ServerAuth
	} else {
		cDes.ServerAuth = des.ServerAuth
	}
	if dcl.BoolCanonicalize(des.ClientAuth, initial.ClientAuth) || dcl.IsZeroValue(des.ClientAuth) {
		cDes.ClientAuth = initial.ClientAuth
	} else {
		cDes.ClientAuth = des.ClientAuth
	}
	if dcl.BoolCanonicalize(des.CodeSigning, initial.CodeSigning) || dcl.IsZeroValue(des.CodeSigning) {
		cDes.CodeSigning = initial.CodeSigning
	} else {
		cDes.CodeSigning = des.CodeSigning
	}
	if dcl.BoolCanonicalize(des.EmailProtection, initial.EmailProtection) || dcl.IsZeroValue(des.EmailProtection) {
		cDes.EmailProtection = initial.EmailProtection
	} else {
		cDes.EmailProtection = des.EmailProtection
	}
	if dcl.BoolCanonicalize(des.TimeStamping, initial.TimeStamping) || dcl.IsZeroValue(des.TimeStamping) {
		cDes.TimeStamping = initial.TimeStamping
	} else {
		cDes.TimeStamping = des.TimeStamping
	}
	if dcl.BoolCanonicalize(des.OcspSigning, initial.OcspSigning) || dcl.IsZeroValue(des.OcspSigning) {
		cDes.OcspSigning = initial.OcspSigning
	} else {
		cDes.OcspSigning = des.OcspSigning
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice(des, initial []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c *Client, des, nw *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.ServerAuth, nw.ServerAuth) {
		nw.ServerAuth = des.ServerAuth
	}
	if dcl.BoolCanonicalize(des.ClientAuth, nw.ClientAuth) {
		nw.ClientAuth = des.ClientAuth
	}
	if dcl.BoolCanonicalize(des.CodeSigning, nw.CodeSigning) {
		nw.CodeSigning = des.CodeSigning
	}
	if dcl.BoolCanonicalize(des.EmailProtection, nw.EmailProtection) {
		nw.EmailProtection = des.EmailProtection
	}
	if dcl.BoolCanonicalize(des.TimeStamping, nw.TimeStamping) {
		nw.TimeStamping = des.TimeStamping
	}
	if dcl.BoolCanonicalize(des.OcspSigning, nw.OcspSigning) {
		nw.OcspSigning = des.OcspSigning
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(des, initial *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) {
		des.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesCaOptions(des, initial *CertificateTemplatePredefinedValuesCaOptions, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesCaOptions{}

	if dcl.BoolCanonicalize(des.IsCa, initial.IsCa) || dcl.IsZeroValue(des.IsCa) {
		cDes.IsCa = initial.IsCa
	} else {
		cDes.IsCa = des.IsCa
	}
	if dcl.IsZeroValue(des.MaxIssuerPathLength) {
		des.MaxIssuerPathLength = initial.MaxIssuerPathLength
	} else {
		cDes.MaxIssuerPathLength = des.MaxIssuerPathLength
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesCaOptionsSlice(des, initial []CertificateTemplatePredefinedValuesCaOptions, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesCaOptions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesCaOptions(c *Client, des, nw *CertificateTemplatePredefinedValuesCaOptions) *CertificateTemplatePredefinedValuesCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.IsCa, nw.IsCa) {
		nw.IsCa = des.IsCa
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesCaOptionsSet(c *Client, des, nw []CertificateTemplatePredefinedValuesCaOptions) []CertificateTemplatePredefinedValuesCaOptions {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesCaOptions
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesCaOptionsSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesCaOptions) []CertificateTemplatePredefinedValuesCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesPolicyIds(des, initial *CertificateTemplatePredefinedValuesPolicyIds, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) {
		des.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesPolicyIdsSlice(des, initial []CertificateTemplatePredefinedValuesPolicyIds, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesPolicyIds(c *Client, des, nw *CertificateTemplatePredefinedValuesPolicyIds) *CertificateTemplatePredefinedValuesPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesPolicyIdsSet(c *Client, des, nw []CertificateTemplatePredefinedValuesPolicyIds) []CertificateTemplatePredefinedValuesPolicyIds {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesPolicyIds
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesPolicyIdsSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesPolicyIds) []CertificateTemplatePredefinedValuesPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensions(des, initial *CertificateTemplatePredefinedValuesAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
	if dcl.BoolCanonicalize(des.Critical, initial.Critical) || dcl.IsZeroValue(des.Critical) {
		cDes.Critical = initial.Critical
	} else {
		cDes.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, initial.Value) || dcl.IsZeroValue(des.Value) {
		cDes.Value = initial.Value
	} else {
		cDes.Value = des.Value
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(des, initial []CertificateTemplatePredefinedValuesAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensions(c *Client, des, nw *CertificateTemplatePredefinedValuesAdditionalExtensions) *CertificateTemplatePredefinedValuesAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsSet(c *Client, des, nw []CertificateTemplatePredefinedValuesAdditionalExtensions) []CertificateTemplatePredefinedValuesAdditionalExtensions {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesAdditionalExtensions
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesAdditionalExtensions) []CertificateTemplatePredefinedValuesAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(des, initial *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) {
		des.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice(des, initial []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c *Client, des, nw *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSet(c *Client, des, nw []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplateIdentityConstraints(des, initial *CertificateTemplateIdentityConstraints, opts ...dcl.ApplyOption) *CertificateTemplateIdentityConstraints {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplateIdentityConstraints{}

	cDes.CelExpression = canonicalizeCertificateTemplateIdentityConstraintsCelExpression(des.CelExpression, initial.CelExpression, opts...)
	if dcl.BoolCanonicalize(des.AllowSubjectPassthrough, initial.AllowSubjectPassthrough) || dcl.IsZeroValue(des.AllowSubjectPassthrough) {
		cDes.AllowSubjectPassthrough = initial.AllowSubjectPassthrough
	} else {
		cDes.AllowSubjectPassthrough = des.AllowSubjectPassthrough
	}
	if dcl.BoolCanonicalize(des.AllowSubjectAltNamesPassthrough, initial.AllowSubjectAltNamesPassthrough) || dcl.IsZeroValue(des.AllowSubjectAltNamesPassthrough) {
		cDes.AllowSubjectAltNamesPassthrough = initial.AllowSubjectAltNamesPassthrough
	} else {
		cDes.AllowSubjectAltNamesPassthrough = des.AllowSubjectAltNamesPassthrough
	}

	return cDes
}

func canonicalizeCertificateTemplateIdentityConstraintsSlice(des, initial []CertificateTemplateIdentityConstraints, opts ...dcl.ApplyOption) []CertificateTemplateIdentityConstraints {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplateIdentityConstraints, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplateIdentityConstraints(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplateIdentityConstraints, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplateIdentityConstraints(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplateIdentityConstraints(c *Client, des, nw *CertificateTemplateIdentityConstraints) *CertificateTemplateIdentityConstraints {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplateIdentityConstraints while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.CelExpression = canonicalizeNewCertificateTemplateIdentityConstraintsCelExpression(c, des.CelExpression, nw.CelExpression)
	if dcl.BoolCanonicalize(des.AllowSubjectPassthrough, nw.AllowSubjectPassthrough) {
		nw.AllowSubjectPassthrough = des.AllowSubjectPassthrough
	}
	if dcl.BoolCanonicalize(des.AllowSubjectAltNamesPassthrough, nw.AllowSubjectAltNamesPassthrough) {
		nw.AllowSubjectAltNamesPassthrough = des.AllowSubjectAltNamesPassthrough
	}

	return nw
}

func canonicalizeNewCertificateTemplateIdentityConstraintsSet(c *Client, des, nw []CertificateTemplateIdentityConstraints) []CertificateTemplateIdentityConstraints {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplateIdentityConstraints
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplateIdentityConstraintsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplateIdentityConstraintsSlice(c *Client, des, nw []CertificateTemplateIdentityConstraints) []CertificateTemplateIdentityConstraints {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplateIdentityConstraints
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplateIdentityConstraints(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplateIdentityConstraintsCelExpression(des, initial *CertificateTemplateIdentityConstraintsCelExpression, opts ...dcl.ApplyOption) *CertificateTemplateIdentityConstraintsCelExpression {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplateIdentityConstraintsCelExpression{}

	if dcl.StringCanonicalize(des.Expression, initial.Expression) || dcl.IsZeroValue(des.Expression) {
		cDes.Expression = initial.Expression
	} else {
		cDes.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, initial.Title) || dcl.IsZeroValue(des.Title) {
		cDes.Title = initial.Title
	} else {
		cDes.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, initial.Location) || dcl.IsZeroValue(des.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}

	return cDes
}

func canonicalizeCertificateTemplateIdentityConstraintsCelExpressionSlice(des, initial []CertificateTemplateIdentityConstraintsCelExpression, opts ...dcl.ApplyOption) []CertificateTemplateIdentityConstraintsCelExpression {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplateIdentityConstraintsCelExpression, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplateIdentityConstraintsCelExpression(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplateIdentityConstraintsCelExpression, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplateIdentityConstraintsCelExpression(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplateIdentityConstraintsCelExpression(c *Client, des, nw *CertificateTemplateIdentityConstraintsCelExpression) *CertificateTemplateIdentityConstraintsCelExpression {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplateIdentityConstraintsCelExpression while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Expression, nw.Expression) {
		nw.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, nw.Title) {
		nw.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}

	return nw
}

func canonicalizeNewCertificateTemplateIdentityConstraintsCelExpressionSet(c *Client, des, nw []CertificateTemplateIdentityConstraintsCelExpression) []CertificateTemplateIdentityConstraintsCelExpression {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplateIdentityConstraintsCelExpression
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplateIdentityConstraintsCelExpressionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplateIdentityConstraintsCelExpressionSlice(c *Client, des, nw []CertificateTemplateIdentityConstraintsCelExpression) []CertificateTemplateIdentityConstraintsCelExpression {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplateIdentityConstraintsCelExpression
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplateIdentityConstraintsCelExpression(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePassthroughExtensions(des, initial *CertificateTemplatePassthroughExtensions, opts ...dcl.ApplyOption) *CertificateTemplatePassthroughExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePassthroughExtensions{}

	if dcl.IsZeroValue(des.KnownExtensions) {
		des.KnownExtensions = initial.KnownExtensions
	} else {
		cDes.KnownExtensions = des.KnownExtensions
	}
	cDes.AdditionalExtensions = canonicalizeCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateTemplatePassthroughExtensionsSlice(des, initial []CertificateTemplatePassthroughExtensions, opts ...dcl.ApplyOption) []CertificateTemplatePassthroughExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePassthroughExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePassthroughExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePassthroughExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePassthroughExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePassthroughExtensions(c *Client, des, nw *CertificateTemplatePassthroughExtensions) *CertificateTemplatePassthroughExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePassthroughExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AdditionalExtensions = canonicalizeNewCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateTemplatePassthroughExtensionsSet(c *Client, des, nw []CertificateTemplatePassthroughExtensions) []CertificateTemplatePassthroughExtensions {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePassthroughExtensions
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePassthroughExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePassthroughExtensionsSlice(c *Client, des, nw []CertificateTemplatePassthroughExtensions) []CertificateTemplatePassthroughExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePassthroughExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePassthroughExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateTemplatePassthroughExtensionsAdditionalExtensions(des, initial *CertificateTemplatePassthroughExtensionsAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateTemplatePassthroughExtensionsAdditionalExtensions{}

	if dcl.IsZeroValue(des.ObjectIdPath) {
		des.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(des, initial []CertificateTemplatePassthroughExtensionsAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateTemplatePassthroughExtensionsAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateTemplatePassthroughExtensionsAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateTemplatePassthroughExtensionsAdditionalExtensions(c *Client, des, nw *CertificateTemplatePassthroughExtensionsAdditionalExtensions) *CertificateTemplatePassthroughExtensionsAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateTemplatePassthroughExtensionsAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateTemplatePassthroughExtensionsAdditionalExtensionsSet(c *Client, des, nw []CertificateTemplatePassthroughExtensionsAdditionalExtensions) []CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return nw
	}
	var reorderedNew []CertificateTemplatePassthroughExtensionsAdditionalExtensions
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareCertificateTemplatePassthroughExtensionsAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c *Client, des, nw []CertificateTemplatePassthroughExtensionsAdditionalExtensions) []CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateTemplatePassthroughExtensionsAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateTemplatePassthroughExtensionsAdditionalExtensions(c, &d, &n))
	}

	return items
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffCertificateTemplate(c *Client, desired, actual *CertificateTemplate, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PredefinedValues, actual.PredefinedValues, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValues, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("PredefinedValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdentityConstraints, actual.IdentityConstraints, dcl.Info{ObjectFunction: compareCertificateTemplateIdentityConstraintsNewStyle, EmptyObject: EmptyCertificateTemplateIdentityConstraints, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("IdentityConstraints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PassthroughExtensions, actual.PassthroughExtensions, dcl.Info{ObjectFunction: compareCertificateTemplatePassthroughExtensionsNewStyle, EmptyObject: EmptyCertificateTemplatePassthroughExtensions, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("PassthroughExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	return newDiffs, nil
}
func compareCertificateTemplatePredefinedValuesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValues)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValues or *CertificateTemplatePredefinedValues", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValues)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValues", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesKeyUsageNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesCaOptionsNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesCaOptions, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesPolicyIdsNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesPolicyIds, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesAdditionalExtensions, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsage or *CertificateTemplatePredefinedValuesKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage or *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DigitalSignature, actual.DigitalSignature, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("DigitalSignature")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContentCommitment, actual.ContentCommitment, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ContentCommitment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyEncipherment, actual.KeyEncipherment, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("KeyEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataEncipherment, actual.DataEncipherment, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("DataEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyAgreement, actual.KeyAgreement, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("KeyAgreement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertSign, actual.CertSign, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("CertSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlSign, actual.CrlSign, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("CrlSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncipherOnly, actual.EncipherOnly, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("EncipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DecipherOnly, actual.DecipherOnly, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("DecipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage or *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ServerAuth, actual.ServerAuth, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ServerAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientAuth, actual.ClientAuth, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ClientAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CodeSigning, actual.CodeSigning, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("CodeSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailProtection, actual.EmailProtection, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("EmailProtection")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TimeStamping, actual.TimeStamping, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("TimeStamping")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OcspSigning, actual.OcspSigning, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("OcspSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages or *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesCaOptions or *CertificateTemplatePredefinedValuesCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesCaOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IsCa, actual.IsCa, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("IsCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxIssuerPathLength, actual.MaxIssuerPathLength, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("MaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesPolicyIds or *CertificateTemplatePredefinedValuesPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesPolicyIds", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesAdditionalExtensions or *CertificateTemplatePredefinedValuesAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.Info{ObjectFunction: compareCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId or *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplateIdentityConstraintsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplateIdentityConstraints)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplateIdentityConstraints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplateIdentityConstraints or *CertificateTemplateIdentityConstraints", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplateIdentityConstraints)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplateIdentityConstraints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplateIdentityConstraints", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CelExpression, actual.CelExpression, dcl.Info{ObjectFunction: compareCertificateTemplateIdentityConstraintsCelExpressionNewStyle, EmptyObject: EmptyCertificateTemplateIdentityConstraintsCelExpression, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("CelExpression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowSubjectPassthrough, actual.AllowSubjectPassthrough, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("AllowSubjectPassthrough")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowSubjectAltNamesPassthrough, actual.AllowSubjectAltNamesPassthrough, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("AllowSubjectAltNamesPassthrough")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplateIdentityConstraintsCelExpressionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplateIdentityConstraintsCelExpression)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplateIdentityConstraintsCelExpression)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplateIdentityConstraintsCelExpression or *CertificateTemplateIdentityConstraintsCelExpression", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplateIdentityConstraintsCelExpression)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplateIdentityConstraintsCelExpression)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplateIdentityConstraintsCelExpression", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Expression, actual.Expression, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Expression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Title, actual.Title, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Title")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePassthroughExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePassthroughExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePassthroughExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePassthroughExtensions or *CertificateTemplatePassthroughExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePassthroughExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePassthroughExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePassthroughExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KnownExtensions, actual.KnownExtensions, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("KnownExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.Info{ObjectFunction: compareCertificateTemplatePassthroughExtensionsAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateTemplatePassthroughExtensionsAdditionalExtensions, OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateTemplatePassthroughExtensionsAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateTemplatePassthroughExtensionsAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateTemplatePassthroughExtensionsAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePassthroughExtensionsAdditionalExtensions or *CertificateTemplatePassthroughExtensionsAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateTemplatePassthroughExtensionsAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateTemplatePassthroughExtensionsAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateTemplatePassthroughExtensionsAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateCertificateTemplateUpdateCertificateTemplateOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *CertificateTemplate) urlNormalized() *CertificateTemplate {
	normalized := dcl.Copy(*r).(CertificateTemplate)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *CertificateTemplate) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateCertificateTemplate" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the CertificateTemplate resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *CertificateTemplate) marshal(c *Client) ([]byte, error) {
	m, err := expandCertificateTemplate(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CertificateTemplate: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalCertificateTemplate decodes JSON responses into the CertificateTemplate resource schema.
func unmarshalCertificateTemplate(b []byte, c *Client) (*CertificateTemplate, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCertificateTemplate(m, c)
}

func unmarshalMapCertificateTemplate(m map[string]interface{}, c *Client) (*CertificateTemplate, error) {

	flattened := flattenCertificateTemplate(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCertificateTemplate expands CertificateTemplate into a JSON request object.
func expandCertificateTemplate(c *Client, f *CertificateTemplate) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s/locations/%s/certificateTemplates/%s", f.Name, f.Project, f.Location, f.Name); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if v != nil {
		m["name"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValues(c, f.PredefinedValues); err != nil {
		return nil, fmt.Errorf("error expanding PredefinedValues into predefinedValues: %w", err)
	} else if v != nil {
		m["predefinedValues"] = v
	}
	if v, err := expandCertificateTemplateIdentityConstraints(c, f.IdentityConstraints); err != nil {
		return nil, fmt.Errorf("error expanding IdentityConstraints into identityConstraints: %w", err)
	} else if v != nil {
		m["identityConstraints"] = v
	}
	if v, err := expandCertificateTemplatePassthroughExtensions(c, f.PassthroughExtensions); err != nil {
		return nil, fmt.Errorf("error expanding PassthroughExtensions into passthroughExtensions: %w", err)
	} else if v != nil {
		m["passthroughExtensions"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if v != nil {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if v != nil {
		m["location"] = v
	}

	return m, nil
}

// flattenCertificateTemplate flattens CertificateTemplate from a JSON request object into the
// CertificateTemplate type.
func flattenCertificateTemplate(c *Client, i interface{}) *CertificateTemplate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &CertificateTemplate{}
	res.Name = dcl.FlattenString(m["name"])
	res.PredefinedValues = flattenCertificateTemplatePredefinedValues(c, m["predefinedValues"])
	res.IdentityConstraints = flattenCertificateTemplateIdentityConstraints(c, m["identityConstraints"])
	res.PassthroughExtensions = flattenCertificateTemplatePassthroughExtensions(c, m["passthroughExtensions"])
	res.Description = dcl.FlattenString(m["description"])
	res.CreateTime = dcl.FlattenString(m["createTime"])
	res.UpdateTime = dcl.FlattenString(m["updateTime"])
	res.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	res.Project = dcl.FlattenString(m["project"])
	res.Location = dcl.FlattenString(m["location"])

	return res
}

// expandCertificateTemplatePredefinedValuesMap expands the contents of CertificateTemplatePredefinedValues into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesMap(c *Client, f map[string]CertificateTemplatePredefinedValues) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValues(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesSlice expands the contents of CertificateTemplatePredefinedValues into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesSlice(c *Client, f []CertificateTemplatePredefinedValues) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValues(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesMap flattens the contents of CertificateTemplatePredefinedValues from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValues{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValues{}
	}

	items := make(map[string]CertificateTemplatePredefinedValues)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValues(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesSlice flattens the contents of CertificateTemplatePredefinedValues from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValues {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValues{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValues{}
	}

	items := make([]CertificateTemplatePredefinedValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValues(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValues expands an instance of CertificateTemplatePredefinedValues into a JSON
// request object.
func expandCertificateTemplatePredefinedValues(c *Client, f *CertificateTemplatePredefinedValues) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateTemplatePredefinedValuesKeyUsage(c, f.KeyUsage); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValuesCaOptions(c, f.CaOptions); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValuesPolicyIdsSlice(c, f.PolicyIds); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["policyIds"] = v
	}
	if v := f.AiaOcspServers; v != nil {
		m["aiaOcspServers"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c, f.AdditionalExtensions); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValues flattens an instance of CertificateTemplatePredefinedValues from a JSON
// response object.
func flattenCertificateTemplatePredefinedValues(c *Client, i interface{}) *CertificateTemplatePredefinedValues {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValues{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValues
	}
	r.KeyUsage = flattenCertificateTemplatePredefinedValuesKeyUsage(c, m["keyUsage"])
	r.CaOptions = flattenCertificateTemplatePredefinedValuesCaOptions(c, m["caOptions"])
	r.PolicyIds = flattenCertificateTemplatePredefinedValuesPolicyIdsSlice(c, m["policyIds"])
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c, m["additionalExtensions"])

	return r
}

// expandCertificateTemplatePredefinedValuesKeyUsageMap expands the contents of CertificateTemplatePredefinedValuesKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageMap(c *Client, f map[string]CertificateTemplatePredefinedValuesKeyUsage) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesKeyUsageSlice expands the contents of CertificateTemplatePredefinedValuesKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageSlice(c *Client, f []CertificateTemplatePredefinedValuesKeyUsage) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageMap flattens the contents of CertificateTemplatePredefinedValuesKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesKeyUsage{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesKeyUsage(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesKeyUsageSlice flattens the contents of CertificateTemplatePredefinedValuesKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesKeyUsage{}
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesKeyUsage(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesKeyUsage expands an instance of CertificateTemplatePredefinedValuesKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsage(c *Client, f *CertificateTemplatePredefinedValuesKeyUsage) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, f.BaseKeyUsage); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsage flattens an instance of CertificateTemplatePredefinedValuesKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsage(c *Client, i interface{}) *CertificateTemplatePredefinedValuesKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsage
	}
	r.BaseKeyUsage = flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, m["baseKeyUsage"])
	r.ExtendedKeyUsage = flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"])
	r.UnknownExtendedKeyUsages = flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"])

	return r
}

// expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageMap expands the contents of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageMap(c *Client, f map[string]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice expands the contents of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice(c *Client, f []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageMap flattens the contents of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice flattens the contents of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage expands an instance of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c *Client, f *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DigitalSignature; !dcl.IsEmptyValueIndirect(v) {
		m["digitalSignature"] = v
	}
	if v := f.ContentCommitment; !dcl.IsEmptyValueIndirect(v) {
		m["contentCommitment"] = v
	}
	if v := f.KeyEncipherment; !dcl.IsEmptyValueIndirect(v) {
		m["keyEncipherment"] = v
	}
	if v := f.DataEncipherment; !dcl.IsEmptyValueIndirect(v) {
		m["dataEncipherment"] = v
	}
	if v := f.KeyAgreement; !dcl.IsEmptyValueIndirect(v) {
		m["keyAgreement"] = v
	}
	if v := f.CertSign; !dcl.IsEmptyValueIndirect(v) {
		m["certSign"] = v
	}
	if v := f.CrlSign; !dcl.IsEmptyValueIndirect(v) {
		m["crlSign"] = v
	}
	if v := f.EncipherOnly; !dcl.IsEmptyValueIndirect(v) {
		m["encipherOnly"] = v
	}
	if v := f.DecipherOnly; !dcl.IsEmptyValueIndirect(v) {
		m["decipherOnly"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage flattens an instance of CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(c *Client, i interface{}) *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	}
	r.DigitalSignature = dcl.FlattenBool(m["digitalSignature"])
	r.ContentCommitment = dcl.FlattenBool(m["contentCommitment"])
	r.KeyEncipherment = dcl.FlattenBool(m["keyEncipherment"])
	r.DataEncipherment = dcl.FlattenBool(m["dataEncipherment"])
	r.KeyAgreement = dcl.FlattenBool(m["keyAgreement"])
	r.CertSign = dcl.FlattenBool(m["certSign"])
	r.CrlSign = dcl.FlattenBool(m["crlSign"])
	r.EncipherOnly = dcl.FlattenBool(m["encipherOnly"])
	r.DecipherOnly = dcl.FlattenBool(m["decipherOnly"])

	return r
}

// expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageMap expands the contents of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice expands the contents of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice(c *Client, f []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageMap flattens the contents of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice flattens the contents of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage expands an instance of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c *Client, f *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ServerAuth; !dcl.IsEmptyValueIndirect(v) {
		m["serverAuth"] = v
	}
	if v := f.ClientAuth; !dcl.IsEmptyValueIndirect(v) {
		m["clientAuth"] = v
	}
	if v := f.CodeSigning; !dcl.IsEmptyValueIndirect(v) {
		m["codeSigning"] = v
	}
	if v := f.EmailProtection; !dcl.IsEmptyValueIndirect(v) {
		m["emailProtection"] = v
	}
	if v := f.TimeStamping; !dcl.IsEmptyValueIndirect(v) {
		m["timeStamping"] = v
	}
	if v := f.OcspSigning; !dcl.IsEmptyValueIndirect(v) {
		m["ocspSigning"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage flattens an instance of CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(c *Client, i interface{}) *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages expands an instance of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c *Client, f *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages flattens an instance of CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}) *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateTemplatePredefinedValuesCaOptionsMap expands the contents of CertificateTemplatePredefinedValuesCaOptions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesCaOptionsMap(c *Client, f map[string]CertificateTemplatePredefinedValuesCaOptions) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesCaOptions(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesCaOptionsSlice expands the contents of CertificateTemplatePredefinedValuesCaOptions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesCaOptionsSlice(c *Client, f []CertificateTemplatePredefinedValuesCaOptions) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesCaOptions(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesCaOptionsMap flattens the contents of CertificateTemplatePredefinedValuesCaOptions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesCaOptionsMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesCaOptions{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesCaOptions)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesCaOptions(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesCaOptionsSlice flattens the contents of CertificateTemplatePredefinedValuesCaOptions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesCaOptionsSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesCaOptions{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesCaOptions{}
	}

	items := make([]CertificateTemplatePredefinedValuesCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesCaOptions(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesCaOptions expands an instance of CertificateTemplatePredefinedValuesCaOptions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesCaOptions(c *Client, f *CertificateTemplatePredefinedValuesCaOptions) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IsCa; !dcl.IsEmptyValueIndirect(v) {
		m["isCa"] = v
	}
	if v := f.MaxIssuerPathLength; !dcl.IsEmptyValueIndirect(v) {
		m["maxIssuerPathLength"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesCaOptions flattens an instance of CertificateTemplatePredefinedValuesCaOptions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesCaOptions(c *Client, i interface{}) *CertificateTemplatePredefinedValuesCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])

	return r
}

// expandCertificateTemplatePredefinedValuesPolicyIdsMap expands the contents of CertificateTemplatePredefinedValuesPolicyIds into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesPolicyIdsMap(c *Client, f map[string]CertificateTemplatePredefinedValuesPolicyIds) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesPolicyIds(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesPolicyIdsSlice expands the contents of CertificateTemplatePredefinedValuesPolicyIds into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesPolicyIdsSlice(c *Client, f []CertificateTemplatePredefinedValuesPolicyIds) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesPolicyIds(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesPolicyIdsMap flattens the contents of CertificateTemplatePredefinedValuesPolicyIds from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesPolicyIdsMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesPolicyIds{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesPolicyIds)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesPolicyIds(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesPolicyIdsSlice flattens the contents of CertificateTemplatePredefinedValuesPolicyIds from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesPolicyIdsSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesPolicyIds{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesPolicyIds{}
	}

	items := make([]CertificateTemplatePredefinedValuesPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesPolicyIds(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesPolicyIds expands an instance of CertificateTemplatePredefinedValuesPolicyIds into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesPolicyIds(c *Client, f *CertificateTemplatePredefinedValuesPolicyIds) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesPolicyIds flattens an instance of CertificateTemplatePredefinedValuesPolicyIds from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesPolicyIds(c *Client, i interface{}) *CertificateTemplatePredefinedValuesPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensionsMap expands the contents of CertificateTemplatePredefinedValuesAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensionsMap(c *Client, f map[string]CertificateTemplatePredefinedValuesAdditionalExtensions) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesAdditionalExtensions(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensionsSlice expands the contents of CertificateTemplatePredefinedValuesAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c *Client, f []CertificateTemplatePredefinedValuesAdditionalExtensions) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesAdditionalExtensions(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensionsMap flattens the contents of CertificateTemplatePredefinedValuesAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensionsMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesAdditionalExtensions{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesAdditionalExtensions(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensionsSlice flattens the contents of CertificateTemplatePredefinedValuesAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensionsSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesAdditionalExtensions{}
	}

	items := make([]CertificateTemplatePredefinedValuesAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesAdditionalExtensions(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensions expands an instance of CertificateTemplatePredefinedValuesAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensions(c *Client, f *CertificateTemplatePredefinedValuesAdditionalExtensions) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, f.ObjectId); err != nil {
		return nil, fmt.Errorf("error expanding ObjectId into objectId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["objectId"] = v
	}
	if v := f.Critical; !dcl.IsEmptyValueIndirect(v) {
		m["critical"] = v
	}
	if v := f.Value; !dcl.IsEmptyValueIndirect(v) {
		m["value"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensions flattens an instance of CertificateTemplatePredefinedValuesAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensions(c *Client, i interface{}) *CertificateTemplatePredefinedValuesAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesAdditionalExtensions
	}
	r.ObjectId = flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, m["objectId"])
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdMap expands the contents of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdMap(c *Client, f map[string]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice expands the contents of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice(c *Client, f []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdMap flattens the contents of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdMap(c *Client, i interface{}) map[string]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice flattens the contents of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSlice(c *Client, i interface{}) []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}
	}

	items := make([]CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId expands an instance of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c *Client, f *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId flattens an instance of CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(c *Client, i interface{}) *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateTemplateIdentityConstraintsMap expands the contents of CertificateTemplateIdentityConstraints into a JSON
// request object.
func expandCertificateTemplateIdentityConstraintsMap(c *Client, f map[string]CertificateTemplateIdentityConstraints) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplateIdentityConstraints(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplateIdentityConstraintsSlice expands the contents of CertificateTemplateIdentityConstraints into a JSON
// request object.
func expandCertificateTemplateIdentityConstraintsSlice(c *Client, f []CertificateTemplateIdentityConstraints) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplateIdentityConstraints(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplateIdentityConstraintsMap flattens the contents of CertificateTemplateIdentityConstraints from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraintsMap(c *Client, i interface{}) map[string]CertificateTemplateIdentityConstraints {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplateIdentityConstraints{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplateIdentityConstraints{}
	}

	items := make(map[string]CertificateTemplateIdentityConstraints)
	for k, item := range a {
		items[k] = *flattenCertificateTemplateIdentityConstraints(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplateIdentityConstraintsSlice flattens the contents of CertificateTemplateIdentityConstraints from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraintsSlice(c *Client, i interface{}) []CertificateTemplateIdentityConstraints {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplateIdentityConstraints{}
	}

	if len(a) == 0 {
		return []CertificateTemplateIdentityConstraints{}
	}

	items := make([]CertificateTemplateIdentityConstraints, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplateIdentityConstraints(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplateIdentityConstraints expands an instance of CertificateTemplateIdentityConstraints into a JSON
// request object.
func expandCertificateTemplateIdentityConstraints(c *Client, f *CertificateTemplateIdentityConstraints) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateTemplateIdentityConstraintsCelExpression(c, f.CelExpression); err != nil {
		return nil, fmt.Errorf("error expanding CelExpression into celExpression: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["celExpression"] = v
	}
	if v := f.AllowSubjectPassthrough; !dcl.IsEmptyValueIndirect(v) {
		m["allowSubjectPassthrough"] = v
	}
	if v := f.AllowSubjectAltNamesPassthrough; !dcl.IsEmptyValueIndirect(v) {
		m["allowSubjectAltNamesPassthrough"] = v
	}

	return m, nil
}

// flattenCertificateTemplateIdentityConstraints flattens an instance of CertificateTemplateIdentityConstraints from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraints(c *Client, i interface{}) *CertificateTemplateIdentityConstraints {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplateIdentityConstraints{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplateIdentityConstraints
	}
	r.CelExpression = flattenCertificateTemplateIdentityConstraintsCelExpression(c, m["celExpression"])
	r.AllowSubjectPassthrough = dcl.FlattenBool(m["allowSubjectPassthrough"])
	r.AllowSubjectAltNamesPassthrough = dcl.FlattenBool(m["allowSubjectAltNamesPassthrough"])

	return r
}

// expandCertificateTemplateIdentityConstraintsCelExpressionMap expands the contents of CertificateTemplateIdentityConstraintsCelExpression into a JSON
// request object.
func expandCertificateTemplateIdentityConstraintsCelExpressionMap(c *Client, f map[string]CertificateTemplateIdentityConstraintsCelExpression) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplateIdentityConstraintsCelExpression(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplateIdentityConstraintsCelExpressionSlice expands the contents of CertificateTemplateIdentityConstraintsCelExpression into a JSON
// request object.
func expandCertificateTemplateIdentityConstraintsCelExpressionSlice(c *Client, f []CertificateTemplateIdentityConstraintsCelExpression) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplateIdentityConstraintsCelExpression(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplateIdentityConstraintsCelExpressionMap flattens the contents of CertificateTemplateIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraintsCelExpressionMap(c *Client, i interface{}) map[string]CertificateTemplateIdentityConstraintsCelExpression {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplateIdentityConstraintsCelExpression{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplateIdentityConstraintsCelExpression{}
	}

	items := make(map[string]CertificateTemplateIdentityConstraintsCelExpression)
	for k, item := range a {
		items[k] = *flattenCertificateTemplateIdentityConstraintsCelExpression(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplateIdentityConstraintsCelExpressionSlice flattens the contents of CertificateTemplateIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraintsCelExpressionSlice(c *Client, i interface{}) []CertificateTemplateIdentityConstraintsCelExpression {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplateIdentityConstraintsCelExpression{}
	}

	if len(a) == 0 {
		return []CertificateTemplateIdentityConstraintsCelExpression{}
	}

	items := make([]CertificateTemplateIdentityConstraintsCelExpression, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplateIdentityConstraintsCelExpression(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplateIdentityConstraintsCelExpression expands an instance of CertificateTemplateIdentityConstraintsCelExpression into a JSON
// request object.
func expandCertificateTemplateIdentityConstraintsCelExpression(c *Client, f *CertificateTemplateIdentityConstraintsCelExpression) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Expression; !dcl.IsEmptyValueIndirect(v) {
		m["expression"] = v
	}
	if v := f.Title; !dcl.IsEmptyValueIndirect(v) {
		m["title"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenCertificateTemplateIdentityConstraintsCelExpression flattens an instance of CertificateTemplateIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCertificateTemplateIdentityConstraintsCelExpression(c *Client, i interface{}) *CertificateTemplateIdentityConstraintsCelExpression {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplateIdentityConstraintsCelExpression{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplateIdentityConstraintsCelExpression
	}
	r.Expression = dcl.FlattenString(m["expression"])
	r.Title = dcl.FlattenString(m["title"])
	r.Description = dcl.FlattenString(m["description"])
	r.Location = dcl.FlattenString(m["location"])

	return r
}

// expandCertificateTemplatePassthroughExtensionsMap expands the contents of CertificateTemplatePassthroughExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensionsMap(c *Client, f map[string]CertificateTemplatePassthroughExtensions) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePassthroughExtensions(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePassthroughExtensionsSlice expands the contents of CertificateTemplatePassthroughExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensionsSlice(c *Client, f []CertificateTemplatePassthroughExtensions) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePassthroughExtensions(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePassthroughExtensionsMap flattens the contents of CertificateTemplatePassthroughExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsMap(c *Client, i interface{}) map[string]CertificateTemplatePassthroughExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePassthroughExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePassthroughExtensions{}
	}

	items := make(map[string]CertificateTemplatePassthroughExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePassthroughExtensions(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePassthroughExtensionsSlice flattens the contents of CertificateTemplatePassthroughExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsSlice(c *Client, i interface{}) []CertificateTemplatePassthroughExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePassthroughExtensions{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePassthroughExtensions{}
	}

	items := make([]CertificateTemplatePassthroughExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePassthroughExtensions(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePassthroughExtensions expands an instance of CertificateTemplatePassthroughExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensions(c *Client, f *CertificateTemplatePassthroughExtensions) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KnownExtensions; v != nil {
		m["knownExtensions"] = v
	}
	if v, err := expandCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c, f.AdditionalExtensions); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePassthroughExtensions flattens an instance of CertificateTemplatePassthroughExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensions(c *Client, i interface{}) *CertificateTemplatePassthroughExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePassthroughExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePassthroughExtensions
	}
	r.KnownExtensions = flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnumSlice(c, m["knownExtensions"])
	r.AdditionalExtensions = flattenCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c, m["additionalExtensions"])

	return r
}

// expandCertificateTemplatePassthroughExtensionsAdditionalExtensionsMap expands the contents of CertificateTemplatePassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensionsAdditionalExtensionsMap(c *Client, f map[string]CertificateTemplatePassthroughExtensionsAdditionalExtensions) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateTemplatePassthroughExtensionsAdditionalExtensions(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice expands the contents of CertificateTemplatePassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c *Client, f []CertificateTemplatePassthroughExtensionsAdditionalExtensions) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateTemplatePassthroughExtensionsAdditionalExtensions(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateTemplatePassthroughExtensionsAdditionalExtensionsMap flattens the contents of CertificateTemplatePassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsAdditionalExtensionsMap(c *Client, i interface{}) map[string]CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePassthroughExtensionsAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePassthroughExtensionsAdditionalExtensions{}
	}

	items := make(map[string]CertificateTemplatePassthroughExtensionsAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePassthroughExtensionsAdditionalExtensions(c, item.(map[string]interface{}))
	}

	return items
}

// flattenCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice flattens the contents of CertificateTemplatePassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsAdditionalExtensionsSlice(c *Client, i interface{}) []CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePassthroughExtensionsAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePassthroughExtensionsAdditionalExtensions{}
	}

	items := make([]CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePassthroughExtensionsAdditionalExtensions(c, item.(map[string]interface{})))
	}

	return items
}

// expandCertificateTemplatePassthroughExtensionsAdditionalExtensions expands an instance of CertificateTemplatePassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCertificateTemplatePassthroughExtensionsAdditionalExtensions(c *Client, f *CertificateTemplatePassthroughExtensionsAdditionalExtensions) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateTemplatePassthroughExtensionsAdditionalExtensions flattens an instance of CertificateTemplatePassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsAdditionalExtensions(c *Client, i interface{}) *CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateTemplatePassthroughExtensionsAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePassthroughExtensionsAdditionalExtensions
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnumMap flattens the contents of CertificateTemplatePassthroughExtensionsKnownExtensionsEnum from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnumMap(c *Client, i interface{}) map[string]CertificateTemplatePassthroughExtensionsKnownExtensionsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateTemplatePassthroughExtensionsKnownExtensionsEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateTemplatePassthroughExtensionsKnownExtensionsEnum{}
	}

	items := make(map[string]CertificateTemplatePassthroughExtensionsKnownExtensionsEnum)
	for k, item := range a {
		items[k] = *flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnumSlice flattens the contents of CertificateTemplatePassthroughExtensionsKnownExtensionsEnum from a JSON
// response object.
func flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnumSlice(c *Client, i interface{}) []CertificateTemplatePassthroughExtensionsKnownExtensionsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateTemplatePassthroughExtensionsKnownExtensionsEnum{}
	}

	if len(a) == 0 {
		return []CertificateTemplatePassthroughExtensionsKnownExtensionsEnum{}
	}

	items := make([]CertificateTemplatePassthroughExtensionsKnownExtensionsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateTemplatePassthroughExtensionsKnownExtensionsEnum with the same value as that string.
func flattenCertificateTemplatePassthroughExtensionsKnownExtensionsEnum(i interface{}) *CertificateTemplatePassthroughExtensionsKnownExtensionsEnum {
	s, ok := i.(string)
	if !ok {
		return CertificateTemplatePassthroughExtensionsKnownExtensionsEnumRef("")
	}

	return CertificateTemplatePassthroughExtensionsKnownExtensionsEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *CertificateTemplate) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCertificateTemplate(b, c)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Project == nil && ncr.Project == nil {
			c.Config.Logger.Info("Both Project fields null - considering equal.")
		} else if nr.Project == nil || ncr.Project == nil {
			c.Config.Logger.Info("Only one Project field is null - considering unequal.")
			return false
		} else if *nr.Project != *ncr.Project {
			return false
		}
		if nr.Location == nil && ncr.Location == nil {
			c.Config.Logger.Info("Both Location fields null - considering equal.")
		} else if nr.Location == nil || ncr.Location == nil {
			c.Config.Logger.Info("Only one Location field is null - considering unequal.")
			return false
		} else if *nr.Location != *ncr.Location {
			return false
		}
		if nr.Name == nil && ncr.Name == nil {
			c.Config.Logger.Info("Both Name fields null - considering equal.")
		} else if nr.Name == nil || ncr.Name == nil {
			c.Config.Logger.Info("Only one Name field is null - considering unequal.")
			return false
		} else if *nr.Name != *ncr.Name {
			return false
		}
		return true
	}
}

type certificateTemplateDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         certificateTemplateApiOperation
}

func convertFieldDiffsToCertificateTemplateDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]certificateTemplateDiff, error) {
	opNamesToFieldDiffs := make(map[string][]*dcl.FieldDiff)
	// Map each operation name to the field diffs associated with it.
	for _, fd := range fds {
		for _, ro := range fd.ResultingOperation {
			if fieldDiffs, ok := opNamesToFieldDiffs[ro]; ok {
				fieldDiffs = append(fieldDiffs, fd)
				opNamesToFieldDiffs[ro] = fieldDiffs
			} else {
				config.Logger.Infof("%s required due to diff in %q", ro, fd.FieldName)
				opNamesToFieldDiffs[ro] = []*dcl.FieldDiff{fd}
			}
		}
	}
	var diffs []certificateTemplateDiff
	// For each operation name, create a certificateTemplateDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := certificateTemplateDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToCertificateTemplateApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToCertificateTemplateApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (certificateTemplateApiOperation, error) {
	switch opName {

	case "updateCertificateTemplateUpdateCertificateTemplateOperation":
		return &updateCertificateTemplateUpdateCertificateTemplateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractCertificateTemplateFields(r *CertificateTemplate) error {
	return nil
}
