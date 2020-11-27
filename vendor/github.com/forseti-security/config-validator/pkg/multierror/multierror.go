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

package multierror

import (
	"fmt"
	"io"
	"strings"
)

// errorImpl handles implementing a list of errors.
type errorImpl []error

// Error implements error.
func (errs errorImpl) Error() string {
	var s []string
	for _, err := range errs {
		s = append(s, err.Error())
	}
	return strings.Join(s, ", ")
}

// Format implements fmt.Formatter to make this play nice with handling stack traces produced from
// github.com/pkg/errors
func (errs errorImpl) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = fmt.Fprintf(s, "errors (%d):\n", len(errs))
		for _, err := range errs {
			if formatter, ok := err.(fmt.Formatter); ok {
				_, _ = io.WriteString(s, "  ")
				formatter.Format(s, verb)
				_, _ = io.WriteString(s, "\n")
			} else {
				_, _ = fmt.Fprintf(s, "  %v\n", err)
			}
		}

	case 's':
		_, _ = io.WriteString(s, errs.Error())

	case 'q':
		_, _ = fmt.Fprintf(s, "%q", errs.Error())
	}
}

// Errors allows for returning multiple errors in one error
type Errors struct {
	errs []error
}

// ToError returns the error if populated, or nil if none exists.
func (e *Errors) ToError() error {
	if len(e.errs) == 0 {
		return nil
	}
	return errorImpl(e.errs)
}

func (e *Errors) Empty() bool {
	return len(e.errs) == 0
}

func (e *Errors) Add(err error) {
	if err == nil {
		return
	}
	if ei, ok := err.(errorImpl); ok {
		e.errs = append(e.errs, ei...)
		return
	}
	e.errs = append(e.errs, err)
}

func (e *Errors) AddF(err error, mod func(error) error) {
	if err == nil {
		return
	}
	e.errs = append(e.errs, mod(err))
}
