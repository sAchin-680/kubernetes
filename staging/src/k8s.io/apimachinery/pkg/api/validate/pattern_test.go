/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validate

import (
	"context"
	"regexp"
	"testing"

	"k8s.io/apimachinery/pkg/api/operation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestPattern(t *testing.T) {
	alphaPattern := regexp.MustCompile(`^[a-z]+$`)
	digitPattern := regexp.MustCompile(`^[0-9]+$`)

	cases := []struct {
		name     string
		value    *string
		pattern  *regexp.Regexp
		wantErrs field.ErrorList
	}{
		{
			name:     "nil value",
			value:    nil,
			pattern:  alphaPattern,
			wantErrs: nil,
		},
		{
			name:     "value matches pattern",
			value:    strPtr("abc"),
			pattern:  alphaPattern,
			wantErrs: nil,
		},
		{
			name:    "value does not match pattern",
			value:   strPtr("abc123"),
			pattern: alphaPattern,
			wantErrs: field.ErrorList{
				field.Invalid(field.NewPath("fldPath"), "", "").WithOrigin("pattern"),
			},
		},
		{
			name:     "empty string matches empty-accepting pattern",
			value:    strPtr(""),
			pattern:  regexp.MustCompile(`^.*$`),
			wantErrs: nil,
		},
		{
			name:    "empty string does not match non-empty pattern",
			value:   strPtr(""),
			pattern: alphaPattern,
			wantErrs: field.ErrorList{
				field.Invalid(field.NewPath("fldPath"), "", "").WithOrigin("pattern"),
			},
		},
		{
			name:     "digit value matches digit pattern",
			value:    strPtr("12345"),
			pattern:  digitPattern,
			wantErrs: nil,
		},
		{
			name:    "alpha value does not match digit pattern",
			value:   strPtr("abc"),
			pattern: digitPattern,
			wantErrs: field.ErrorList{
				field.Invalid(field.NewPath("fldPath"), "", "").WithOrigin("pattern"),
			},
		},
	}

	matcher := field.ErrorMatcher{}.ByType().ByField().ByOrigin()
	fldPath := field.NewPath("fldPath")
	op := operation.Operation{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := Pattern(context.Background(), op, fldPath, tc.value, nil, tc.pattern)
			matcher.Test(t, tc.wantErrs, errs)
		})
	}
}
