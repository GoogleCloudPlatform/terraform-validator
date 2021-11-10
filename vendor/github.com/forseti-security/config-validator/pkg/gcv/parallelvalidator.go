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

package gcv

import (
	"context"
	"flag"
	"runtime"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/multierror"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var flags struct {
	workerCount int
}

func init() {
	flag.IntVar(
		&flags.workerCount,
		"workerCount",
		runtime.NumCPU(),
		"Number of workers that Validator will spawn to handle validate calls, this defaults to core count on the host")
}

// ParallelValidator handles making parallel calls to Validator during a Review call.
type ParallelValidator struct {
	cv   ConfigValidator
	work chan func()
}

type assetResult struct {
	violations []*validator.Violation
	err        error
}

// NewParallelValidator creates a new instance with the given stop channel and validator
func NewParallelValidator(stopChannel <-chan struct{}, cv ConfigValidator) *ParallelValidator {
	pv := &ParallelValidator{
		// channel size of number of workers seems sufficient to prevent blocking,
		// this is really just an assumption with no actual perf benchmarking.
		work: make(chan func(), flags.workerCount),
		cv:   cv,
	}

	go func() {
		<-stopChannel
		glog.Infof("validator shutdown requested via stopChannel close")
		close(pv.work)
	}()

	workerCount := flags.workerCount
	glog.Infof("validator starting %d workers", workerCount)
	for i := 0; i < workerCount; i++ {
		go pv.reviewWorker(i)
	}

	return pv
}

// reviewWorker is the function that each worker goroutine will use
func (v *ParallelValidator) reviewWorker(idx int) {
	glog.V(1).Infof("worker %d starting", idx)
	for f := range v.work {
		f()
	}
	glog.V(1).Infof("worker %d terminated", idx)
}

// handleReview is the wrapper function for individual asset reviews.
func (v *ParallelValidator) handleReview(ctx context.Context, idx int, asset *validator.Asset, resultChan chan<- *assetResult) func() {
	return func() {
		resultChan <- func() *assetResult {
			violations, err := v.cv.ReviewAsset(ctx, asset)
			if err != nil {
				return &assetResult{err: errors.Wrapf(err, "index %d", idx)}
			}
			return &assetResult{violations: violations}
		}()
	}
}

// Review evaluates each asset in the review request in parallel and returns any
// violations found.
func (v *ParallelValidator) Review(ctx context.Context, request *validator.ReviewRequest) (*validator.ReviewResponse, error) {
	assetCount := len(request.Assets)
	// channel size of number of workers seems sufficient to prevent blocking,
	// this is really just an assumption with no actual perf benchmarking.
	resultChan := make(chan *assetResult, flags.workerCount)
	defer close(resultChan)

	go func() {
		for idx, asset := range request.Assets {
			v.work <- v.handleReview(ctx, idx, asset, resultChan)
		}
	}()

	response := &validator.ReviewResponse{}
	var errs multierror.Errors
	for i := 0; i < assetCount; i++ {
		result := <-resultChan
		if result.err != nil {
			errs.Add(result.err)
			continue
		}
		response.Violations = append(response.Violations, result.violations...)
	}

	if !errs.Empty() {
		return response, errs.ToError()
	}
	return response, nil
}
