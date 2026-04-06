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

package pattern

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func Test(t *testing.T) {
	st := localSchemeBuilder.Test(t)

	// Valid: all fields set to a matching value.
	st.Value(&Struct{
		PatternField:                      "abc",
		PatternPtrField:                   ptr.To("abc"),
		PatternUnvalidatedTypedefField:    UnvalidatedStringType("abc"),
		PatternUnvalidatedTypedefPtrField: ptr.To(UnvalidatedStringType("abc")),
		PatternValidatedTypedefField:      PatternType("abc"),
		PatternValidatedTypedefPtrField:   ptr.To(PatternType("abc")),
	}).ExpectValid()

	// Valid: pointer fields are nil (Pattern() accepts nil).
	st.Value(&Struct{
		PatternField:                   "abc",
		PatternUnvalidatedTypedefField: UnvalidatedStringType("abc"),
		PatternValidatedTypedefField:   PatternType("abc"),
	}).ExpectValid()

	// Invalid: all fields set to a non-matching value.
	testVal := &Struct{
		PatternField:                      "123",
		PatternPtrField:                   ptr.To("123"),
		PatternUnvalidatedTypedefField:    UnvalidatedStringType("123"),
		PatternUnvalidatedTypedefPtrField: ptr.To(UnvalidatedStringType("123")),
		PatternValidatedTypedefField:      PatternType("123"),
		PatternValidatedTypedefPtrField:   ptr.To(PatternType("123")),
	}
	st.Value(testVal).ExpectMatches(field.ErrorMatcher{}.ByType().ByField(), field.ErrorList{
		field.Invalid(field.NewPath("patternField"), "", ""),
		field.Invalid(field.NewPath("patternPtrField"), "", ""),
		field.Invalid(field.NewPath("patternUnvalidatedTypedefField"), "", ""),
		field.Invalid(field.NewPath("patternUnvalidatedTypedefPtrField"), "", ""),
		field.Invalid(field.NewPath("patternValidatedTypedefField"), "", ""),
		field.Invalid(field.NewPath("patternValidatedTypedefPtrField"), "", ""),
	})

	// Test validation ratcheting: updating a non-matching value to the same
	// non-matching value should pass (unchanged data is not re-validated).
	st.Value(testVal).OldValue(testVal).ExpectValid()
}
