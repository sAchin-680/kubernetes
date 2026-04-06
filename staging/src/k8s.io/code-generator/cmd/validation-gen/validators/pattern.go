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

package validators

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/code-generator/cmd/validation-gen/util"
	"k8s.io/gengo/v2/codetags"
	"k8s.io/gengo/v2/types"
)

const (
	patternTagName = "k8s:pattern"
)

func init() {
	RegisterTagValidator(patternTagValidator{})
}

type patternTagValidator struct{}

func (patternTagValidator) Init(_ Config) {}

func (patternTagValidator) TagName() string {
	return patternTagName
}

var patternTagValidScopes = sets.New(ScopeType, ScopeField, ScopeListVal, ScopeMapKey, ScopeMapVal)

func (patternTagValidator) ValidScopes() sets.Set[Scope] {
	return patternTagValidScopes
}

var (
	patternValidator  = types.Name{Package: libValidationPkg, Name: "Pattern"}
	regexpMustCompile = types.Name{Package: "regexp", Name: "MustCompile"}
)

func (patternTagValidator) GetValidations(context Context, tag codetags.Tag) (Validations, error) {
	// This tag can apply to value and pointer fields, as well as typedefs
	// (which should never be pointers). We need to check the concrete type.
	if t := util.NonPointer(util.NativeType(context.Type)); t != types.String {
		return Validations{}, fmt.Errorf("can only be used on string types (%s)", rootTypeString(context.Type, t))
	}

	pattern := tag.Value
	// Validate the regex at code generation time to catch invalid patterns early.
	if _, err := regexp.Compile(pattern); err != nil {
		return Validations{}, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}

	var result Validations

	// Generate a package-level variable for the compiled regexp.
	// The regex is compiled once at startup, not on every validation call.
	// TODO: Avoid the "local" here. This was added to avoid errors caused when the package is an empty string.
	//       The correct package would be the output package but is not known here. This does not show up in generated code.
	// TODO: Append a consistent hash suffix to avoid generated name conflicts?
	patternVarName := PrivateVar{Name: "patternFor_" + sanitizeName(context.Path.String()), Package: "local"}
	result.AddVariable(Variable(patternVarName, Function(patternTagName, DefaultFlags, regexpMustCompile, Literal(fmt.Sprintf("%q", pattern)))))
	result.AddFunction(Function(patternTagName, DefaultFlags, patternValidator, patternVarName))

	return result, nil
}

func (ptv patternTagValidator) Docs() TagDoc {
	return TagDoc{
		Tag:            ptv.TagName(),
		StabilityLevel: TagStabilityLevelAlpha,
		Scopes:         sets.List(ptv.ValidScopes()),
		Description: `Indicates that a string field must match the specified regular expression pattern.
The pattern uses Go's RE2 engine (no lookahead/lookbehind or backreferences).
The pattern is validated at code generation time; invalid patterns are caught as lint errors.
The compiled regexp is cached as a package-level variable for performance.`,
		Payloads: []TagPayloadDoc{{
			Description: "<regex>",
			Docs:        "The regular expression pattern the string value must match.",
		}},
		PayloadsType:     codetags.ValueTypeString,
		PayloadsRequired: true,
	}
}
