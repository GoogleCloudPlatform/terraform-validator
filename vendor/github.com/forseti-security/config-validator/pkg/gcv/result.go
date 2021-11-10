// Copyright 2020 Google LLC
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

package gcv

import (
	"encoding/json"
	"fmt"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv/cf"
	"github.com/forseti-security/config-validator/pkg/gcv/configs"
	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	cftypes "github.com/open-policy-agent/frameworks/constraint/pkg/types"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ConstraintKey = "constraint"
)

// Result is the result of reviewing an individual resource
type Result struct {
	// The name of the resource as given by CAI
	Name string
	// CAIResource is the resource as given by CAI
	CAIResource map[string]interface{}
	// ReviewResource is the resource sent to Constraint Framework for review.
	// For GCP types this is the unmodified resource from CAI, for K8S types, this is the unwrapped
	// resource.
	ReviewResource map[string]interface{}
	// ConstraintViolations are the constraints that were not satisfied during review.
	ConstraintViolations []ConstraintViolation
}

// NewResult creates a Result from the provided CF Response.
func NewResult(
	target string,
	caiResource map[string]interface{},
	reviewResource map[string]interface{},
	responses *cftypes.Responses) (*Result, error) {
	cfResponse, found := responses.ByTarget[target]
	if !found {
		return nil, errors.Errorf("No response for target %s", target)
	}

	resNameIface, found := caiResource["name"]
	if !found {
		return nil, errors.Errorf("result missing name field")
	}
	name, ok := resNameIface.(string)
	if !ok {
		return nil, errors.Errorf("failed to convert resource name to string %v", resNameIface)
	}

	result := &Result{
		Name:                 name,
		CAIResource:          caiResource,
		ReviewResource:       reviewResource,
		ConstraintViolations: make([]ConstraintViolation, len(cfResponse.Results)),
	}
	for idx, cfResult := range cfResponse.Results {
		for k, _ := range cfResult.Metadata {
			if k == ConstraintKey {
				return nil, errors.Errorf("constraint template metadata contains reserved key %s", ConstraintKey)
			}
		}
		severity, found, err := unstructured.NestedString(cfResult.Constraint.Object, "spec", "severity")
		if err != nil || !found {
			severity = ""
		}
		result.ConstraintViolations[idx] = ConstraintViolation{
			Message:    cfResult.Msg,
			Metadata:   cfResult.Metadata,
			Constraint: cfResult.Constraint,
			Severity:   severity,
		}
	}
	return result, nil
}

// ConstraintViolations represents an unsatisfied constraint
type ConstraintViolation struct {
	// Message is a human readable message for the violation
	Message string
	// Metadata is the metadata returned by the constraint check
	Metadata map[string]interface{}
	// Constraint is the K8S resource of the constraint that triggered the violation
	Constraint *unstructured.Unstructured
	// Constraint Severity
	Severity string
}

// ToInsights returns the result represented as a slice of insights.
func (r *Result) ToInsights() []*Insight {
	if len(r.ConstraintViolations) == 0 {
		return nil
	}

	insights := make([]*Insight, len(r.ConstraintViolations))
	for idx, cv := range r.ConstraintViolations {
		i := &Insight{
			Description:     cv.Message,
			TargetResources: []string{r.Name},
			InsightSubtype:  cv.name(),
			Content: map[string]interface{}{
				"resource": r.CAIResource,
				"metadata": cv.metadata(nil),
			},
			Category: "SECURITY",
		}
		insights[idx] = i
	}
	return insights
}

func (r *Result) ToViolations() ([]*validator.Violation, error) {
	ancestryPath, found, err := unstructured.NestedString(r.CAIResource, ancestryPathKey)
	if err != nil {

		return nil, errors.Wrapf(err, "error getting ancestry path from %v", r.CAIResource)
	}
	if !found {
		return nil, errors.Errorf("ancestry path not found in %v", r.CAIResource)
	}

	var violations []*validator.Violation
	for _, rv := range r.ConstraintViolations {
		violation, err := rv.toViolation(r.Name, ancestryPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert result")
		}
		violations = append(violations, violation)
	}
	return violations, nil
}

func (cv *ConstraintViolation) metadata(auxMetadata map[string]interface{}) map[string]interface{} {
	labels := cv.Constraint.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	annotations := cv.Constraint.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	params, found, err := unstructured.NestedMap(cv.Constraint.Object, "spec", "parameters")
	if err != nil {
		panic(fmt.Sprintf(
			"constraint has invalid schema (%#v), should have already been validated, "+
				" .spec.parameters got schema error on access: %s", cv.Constraint.Object, err))
	}
	if !found {
		params = map[string]interface{}{}
	}
	metadata := map[string]interface{}{
		ConstraintKey: map[string]interface{}{
			"labels":      labels,
			"annotations": annotations,
			"parameters":  params,
		},
	}
	for k, v := range auxMetadata {
		metadata[k] = v
	}
	for k, v := range cv.Metadata {
		metadata[k] = v
	}
	return metadata
}

// name returns the name for the constraint, this is given as "[Kind].[Name]" to uniquely identify which template and
// constraint the violation came from.
func (cv *ConstraintViolation) name() string {
	name := cv.Constraint.GetName()
	ans := cv.Constraint.GetAnnotations()
	if ans != nil {
		if originalName, ok := ans[configs.OriginalName]; ok {
			name = originalName
		}
	}
	return fmt.Sprintf("%s.%s", cv.Constraint.GetKind(), name)
}

// toViolation converts the constriant to a violation.
func (cv *ConstraintViolation) toViolation(name string, ancestryPath string) (*validator.Violation, error) {
	metadataJson, err := json.Marshal(cv.metadata(map[string]interface{}{ancestryPathKey: ancestryPath}))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal result metadata %v to json", cv.Metadata)
	}
	metadata := &structpb.Value{}
	if err := jsonpb.UnmarshalString(string(metadataJson), metadata); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal json %s into structpb", string(metadataJson))
	}

	// Extract the object fields if they exists.
	var apiVersion string
	if constraintAPIVersion, ok := cv.Constraint.Object["apiVersion"]; ok {
		apiVersion = fmt.Sprintf("%s", constraintAPIVersion)
	}

	var kind string
	if constraintKind, ok := cv.Constraint.Object["kind"]; ok {
		kind = fmt.Sprintf("%s", constraintKind)
	}

	var pbMetadata *structpb.Value
	if constraintMetadata, ok := cv.Constraint.Object["metadata"]; ok {
		if pbMetadata, err = cf.ConvertToProtoVal(constraintMetadata); err != nil {
			return nil, errors.Wrapf(err, "failed to convert constraint metadata into structpb.Value")
		}
	}

	var pbSpec *structpb.Value
	if constraintSpec, ok := cv.Constraint.Object["spec"]; ok {
		if pbSpec, err = cf.ConvertToProtoVal(constraintSpec); err != nil {
			return nil, errors.Wrapf(err, "failed to convert constraint spec into structpb.Value")
		}
	}

	// Build the ConstraintConfig proto.
	constraintConfig := &validator.Constraint{
		ApiVersion: apiVersion,
		Kind:       kind,
		Metadata:   pbMetadata,
		Spec:       pbSpec,
	}

	return &validator.Violation{
		Constraint:       cv.name(),
		ConstraintConfig: constraintConfig,
		Resource:         name,
		Message:          cv.Message,
		Metadata:         metadata,
		Severity:         cv.Severity,
	}, nil
}
