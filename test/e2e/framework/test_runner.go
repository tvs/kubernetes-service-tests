/*
Copyright 2024.

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
	"context"
	"log"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/types"
)

type TestRunnerBuilder struct {
	steps []testSequence
}

// NewTestRunner provides a builder for test runs
func NewTestRunner() *TestRunnerBuilder {
	return &TestRunnerBuilder{}
}

// WithSequence adds a generic sequence to the test run
func (b *TestRunnerBuilder) WithSequence(s testSequence) *TestRunnerBuilder {
	b.steps = append(b.steps, s)
	return b
}

// WithSerialSequence adds a serialized sequence to the test run
func (b *TestRunnerBuilder) WithSerialSequence(features ...types.Feature) *TestRunnerBuilder {
	b.steps = append(b.steps, NewSerialSequence(features...))
	return b
}

// WithParallelSequence adds a parallel sequence to the test run
func (b *TestRunnerBuilder) WithParallelSequence(features ...types.Feature) *TestRunnerBuilder {
	b.steps = append(b.steps, NewParallelSequence(features...))
	return b
}

// Runner validates and returns a TestRunner configured by the builder.
func (b *TestRunnerBuilder) Runner() *testRunner {
	// Start by validating tests before we run them. This way we can bail out
	// angrily before we start anything.
	for _, s := range b.steps {
		s.Validate()
	}
	if err := FormatBugs(); err != nil {
		log.Fatalf("ERROR: E2E suite initialization was faulty, these errors must be fixed:\n%s", err)
	}

	return &testRunner{
		sequence: b.steps,
	}
}

// testRunner is a mechanism for validating and executing tests in a series of
// sequences.
// Each item in a sequence is a collection of tests themselves and each
// collection may have its own test execution semantics, such as running in
// serial or parallel.
type testRunner struct {
	sequence []testSequence
}

// Test runs the test sequences.
func (tr *testRunner) Test(t *testing.T, tc *TestContextType) context.Context {
	// Context is shared off the TestEnv. We probably don't need to throw it back
	// but we choose to send it back anyway just to preserve our wrapping around
	// the [TestEnv.Test] calls.
	var ctx context.Context
	for _, s := range tr.sequence {
		ctx = s.Test(t, tc)
	}

	return ctx
}
