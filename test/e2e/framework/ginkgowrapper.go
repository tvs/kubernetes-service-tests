/*
Copyright 2022 The Kubernetes Authors.

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

package framework

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Environment is the name for the environment in which a test can run, like
// "Linux" or "Windows".
type Environment string

type Valid[T comparable] struct {
	items  sets.Set[T]
	frozen bool
}

// Add registers a new valid item name. The expected usage is
//
//	var SomeEnvironment = framework.ValidEnvironments.Add("Some")
//
// during the init phase of an E2E suite. Individual tests should not register
// their own, to avoid uncontrolled proliferation of new items. E2E suites can,
// but don't have to, enforce that by freezing the set of valid names.
func (v *Valid[T]) Add(item T) T {
	if v.frozen {
		RecordBug(NewBug(fmt.Sprintf(`registry %T is already frozen, "%v" must not be added anymore`, *v, item), 1))
	}
	if v.items == nil {
		v.items = sets.New[T]()
	}
	if v.items.Has(item) {
		RecordBug(NewBug(fmt.Sprintf(`registry %T already contains "%v", it must not be added again`, *v, item), 1))
	}
	v.items.Insert(item)
	return item
}

func (v *Valid[T]) Freeze() {
	v.frozen = true
}

// These variables contain the parameters that [WithEnvironment] accepts. The
// framework itself has no pre-defined constants. Test suites and tests may
// define their own and then add them here before calling these With functions.
var (
	ValidEnvironments Valid[Environment]
)

var errInterface = reflect.TypeOf((*error)(nil)).Elem()

// IgnoreNotFound can be used to wrap an arbitrary function in a call to
// [ginkgo.DeferCleanup]. When the wrapped function returns an error that
// `apierrors.IsNotFound` considers as "not found", the error is ignored
// instead of failing the test during cleanup. This is useful for cleanup code
// that just needs to ensure that some object does not exist anymore.
func IgnoreNotFound(in any) any {
	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)
	return reflect.MakeFunc(inType, func(args []reflect.Value) []reflect.Value {
		out := inValue.Call(args)
		if len(out) > 0 {
			lastValue := out[len(out)-1]
			last := lastValue.Interface()
			if last != nil && lastValue.Type().Implements(errInterface) && apierrors.IsNotFound(last.(error)) {
				out[len(out)-1] = reflect.Zero(errInterface)
			}
		}
		return out
	}).Interface()
}

// AnnotatedLocation can be used to provide more informative source code
// locations by passing the result as additional parameter to a
// BeforeEach/AfterEach/DeferCleanup/It/etc.
func AnnotatedLocation(annotation string) types.CodeLocation {
	return AnnotatedLocationWithOffset(annotation, 1)
}

// AnnotatedLocationWithOffset skips additional call stack levels. With 0 as offset
// it is identical to [AnnotatedLocation].
func AnnotatedLocationWithOffset(annotation string, offset int) types.CodeLocation {
	codeLocation := types.NewCodeLocation(offset + 1)
	codeLocation.FileName = path.Base(codeLocation.FileName)
	codeLocation = types.NewCustomCodeLocation(annotation + " | " + codeLocation.String())
	return codeLocation
}

// KubernetesServiceDescribe returns a wrapper function for [ginkgo.Describe].
// Adds "[KubernetesService]" tag and makes static analysis easier.
func KubernetesServiceDescribe(args ...interface{}) bool {
	args = append(args, ginkgo.Offset(1), WithKubernetesService())
	return Describe(args...)
}

// WorkloadClusterDescribe returns a wrapper function for [ginkgo.Describe].
// Adds "[WorkloadCluster]" tag and makes static analysis easier.
func WorkloadClusterDescribe(args ...interface{}) bool {
	args = append(args, ginkgo.Offset(1), WithWorkloadCluster())
	return Describe(args...)
}

// ConformanceIt is wrapper function for [ginkgo.It]. Adds "[Conformance]" tag
// and makes static analysis easier.
func ConformanceIt(args ...interface{}) bool {
	args = append(args, ginkgo.Offset(1), WithConformance())
	return It(args...)
}

// It is a wrapper around [ginkgo.It] which supports framework With* labels as
// optional arguments in addition to those already supported by ginkgo itself,
// like [ginkgo.Label] and [ginkgo.Offset].
//
// Text and arguments may be mixed. The final text is a concatenation
// of the text arguments and special tags from the With functions.
func It(args ...interface{}) bool {
	return registerInSuite(ginkgo.It, args)
}

// Describe is a wrapper around [ginkgo.Describe] which supports framework
// With* labels as optional arguments in addition to those already supported by
// ginkgo itself, like [ginkgo.Label] and [ginkgo.Offset].
//
// Text and arguments may be mixed. The final text is a concatenation
// of the text arguments and special tags from the With functions.
func Describe(args ...interface{}) bool {
	return registerInSuite(ginkgo.Describe, args)
}

// Context is a wrapper around [ginkgo.Context] which supports framework With*
// labels as optional arguments in addition to those already supported by
// ginkgo itself, like [ginkgo.Label] and [ginkgo.Offset].
//
// Text and arguments may be mixed. The final text is a concatenation
// of the text arguments and special tags from the With functions.
func Context(args ...interface{}) bool {
	return registerInSuite(ginkgo.Context, args)
}

// registerInSuite is the common implementation of all wrapper functions. It
// expects to be called through one intermediate wrapper.
func registerInSuite(ginkgoCall func(string, ...interface{}) bool, args []interface{}) bool {
	var ginkgoArgs []interface{}
	var offset ginkgo.Offset
	var texts []string

	addLabel := func(label string) {
		texts = append(texts, fmt.Sprintf("[%s]", label))
		ginkgoArgs = append(ginkgoArgs, ginkgo.Label(label))
	}

	haveEmptyStrings := false
	for _, arg := range args {
		switch arg := arg.(type) {
		case label:
			fullLabel := strings.Join(arg.parts, ":")
			addLabel(fullLabel)
			if arg.extra != "" {
				addLabel(arg.extra)
			}
			if fullLabel == "Serial" {
				ginkgoArgs = append(ginkgoArgs, ginkgo.Serial)
			}
		case ginkgo.Offset:
			offset = arg
		case string:
			if arg == "" {
				haveEmptyStrings = true
			}
			texts = append(texts, arg)
		default:
			ginkgoArgs = append(ginkgoArgs, arg)
		}
	}
	offset += 2 // This function and its direct caller.

	// Now that we have the final offset, we can record bugs.
	if haveEmptyStrings {
		RecordBug(NewBug("empty strings as separators are unnecessary and need to be removed", int(offset)))
	}

	// Enforce that text snippets to not start or end with spaces because
	// those lead to double spaces when concatenating below.
	for _, text := range texts {
		if strings.HasPrefix(text, " ") || strings.HasSuffix(text, " ") {
			RecordBug(NewBug(fmt.Sprintf("trailing or leading spaces are unnecessary and need to be removed: %q", text), int(offset)))
		}
	}

	ginkgoArgs = append(ginkgoArgs, offset)
	text := strings.Join(texts, " ")
	return ginkgoCall(text, ginkgoArgs...)
}

var (
	tagRe                 = regexp.MustCompile(`\[.*?\]`)
	deprecatedTags        = sets.New("Conformance", "Flaky", "NodeConformance", "Disruptive", "Serial", "Slow")
	deprecatedTagPrefixes = sets.New("Environment", "Feature", "NodeFeature", "FeatureGate")
	deprecatedStability   = sets.New("Alpha", "Beta")
)

// validateSpecs checks that the test specs were registered as intended.
func validateSpecs(specs types.SpecReports) {
	checked := sets.New[call]()

	for _, spec := range specs {
		for i, text := range spec.ContainerHierarchyTexts {
			c := call{
				text:     text,
				location: spec.ContainerHierarchyLocations[i],
			}
			if checked.Has(c) {
				// No need to check the same container more than once.
				continue
			}
			checked.Insert(c)
			validateText(c.location, text, spec.ContainerHierarchyLabels[i])
		}
		c := call{
			text:     spec.LeafNodeText,
			location: spec.LeafNodeLocation,
		}
		if !checked.Has(c) {
			validateText(spec.LeafNodeLocation, spec.LeafNodeText, spec.LeafNodeLabels)
			checked.Insert(c)
		}
	}
}

// call acts as (mostly) unique identifier for a container node call like
// Describe or Context. It's not perfect because theoretically a line might
// have multiple calls with the same text, but that isn't a problem in
// practice.
type call struct {
	text     string
	location types.CodeLocation
}

// validateText checks for some known tags that should not be added through the
// plain text strings anymore. Eventually, all such tags should get replaced
// with the new APIs.
func validateText(location types.CodeLocation, text string, labels []string) {
	for _, tag := range tagRe.FindAllString(text, -1) {
		if tag == "[]" {
			recordTextBug(location, "[] in plain text is invalid")
			continue
		}
		// Strip square brackets.
		tag = tag[1 : len(tag)-1]
		if slices.Contains(labels, tag) {
			// Okay, was also set as label.
			continue
		}
		if deprecatedTags.Has(tag) {
			recordTextBug(location, fmt.Sprintf("[%s] in plain text is deprecated and must be added through With%s instead", tag, tag))
		}
		if deprecatedStability.Has(tag) {
			recordTextBug(location, fmt.Sprintf("[%s] in plain text is deprecated and must be added by defining the feature gate through WithFeatureGate instead", tag))
		}
		if index := strings.Index(tag, ":"); index > 0 {
			prefix := tag[:index]
			if deprecatedTagPrefixes.Has(prefix) {
				recordTextBug(location, fmt.Sprintf("[%s] in plain text is deprecated and must be added through With%s(%s) instead", tag, prefix, tag[index+1:]))
			}
		}
	}
}

func recordTextBug(location types.CodeLocation, message string) {
	RecordBug(Bug{FileName: location.FileName, LineNumber: location.LineNumber, Message: message})
}

// WithEnvironment specifies that a certain test or group of tests only works
// in a certain environment. The return value must be passed as additional
// argument to [framework.It], [framework.Describe], [framework.Context].
//
// The environment must be listed in ValidEnvironments.
func WithEnvironment(name Environment) interface{} {
	return withEnvironment(name)
}

func withEnvironment(name Environment) interface{} {
	if !ValidEnvironments.items.Has(name) {
		RecordBug(NewBug(fmt.Sprintf("WithEnvironment: unknown environment %q", name), 2))
	}
	return newLabel("Environment", string(name))
}

// WithKubernetesService specifies that a certain test or group of tests are
// related to the Kubernetes Service. The return value must be passed as an
// additional argument to [framework.It], [framework.Describe], or
// [framework.Context].
func WithKubernetesService() interface{} {
	return withKubernetesService()
}

func withKubernetesService() interface{} {
	return newLabel("KubernetesService")
}

// WithWorkloadCluster specifies that a certain test or group of tests are
// related to the Kubernetes Service's workload cluster. The return value must
// be passed as an additional argument to [framework.It], [framework.Describe],
// or [framework.Context].
func WithWorkloadCluster() interface{} {
	return withWorkloadCluster()
}

func withWorkloadCluster() interface{} {
	return newLabel("WorkloadCluster")
}

// WithConformance specifies that a certain test or group of tests must pass in
// all conformant Kubernetes clusters. The return value must be passed as
// additional argument to [framework.It], [framework.Describe],
// [framework.Context].
func WithConformance() interface{} {
	return withConformance()
}

func withConformance() interface{} {
	return newLabel("Conformance")
}

// WithDisruptive specifies that a certain test or group of tests temporarily
// affects the functionality of the Kubernetes cluster. The return value must
// be passed as additional argument to [framework.It], [framework.Describe],
// [framework.Context].
func WithDisruptive() interface{} {
	return withDisruptive()
}

func withDisruptive() interface{} {
	return newLabel("Disruptive")
}

// WithSlow specifies that a certain test or group of tests must not run in
// parallel with other tests. The return value must be passed as additional
// argument to [framework.It], [framework.Describe], [framework.Context].
func WithSlow() interface{} {
	return withSlow()
}

func withSlow() interface{} {
	return newLabel("Slow")
}

// WithLabel is a wrapper around [ginkgo.Label]. Besides adding an arbitrary
// label to a test, it also injects the label in square brackets into the test
// name.
func WithLabel(label string) interface{} {
	return withLabel(label)
}

func withLabel(label string) interface{} {
	return newLabel(label)
}

// WithFlaky specifies that a certain test or group of tests are failing randomly.
// These tests are usually filtered out and ran separately from other tests.
func WithFlaky() interface{} {
	return withFlaky()
}

func withFlaky() interface{} {
	return newLabel("Flaky")
}

type label struct {
	// parts get concatenated with ":" to build the full label.
	parts []string
	// extra is an optional fully-formed extra label.
	extra string
	// explanation gets set for each label to help developers
	// who pass a label to a ginkgo function. They need to use
	// the corresponding framework function instead.
	explanation string
}

func newLabel(parts ...string) label {
	return label{
		parts:       parts,
		explanation: "If you see this as part of an 'Unknown Decorator' error from Ginkgo, then you need to replace the ginkgo.It/Context/Describe call with the corresponding framework.It/Context/Describe or (if available) f.It/Context/Describe.",
	}
}

// TagsEqual can be used to check whether two tags are the same.
// It's safe to compare e.g. the result of WithSlow() against the result
// of WithSerial(), the result will be false. False is also returned
// when a parameter is some completely different value.
func TagsEqual(a, b interface{}) bool {
	al, ok := a.(label)
	if !ok {
		return false
	}
	bl, ok := b.(label)
	if !ok {
		return false
	}
	if al.extra != bl.extra {
		return false
	}
	return slices.Equal(al.parts, bl.parts)
}
