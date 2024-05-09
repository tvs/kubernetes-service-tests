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
	"fmt"
	"math/rand"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/types"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework/testlabels"
)

// internal interface for sequences of feature tests
type testSequence interface {
	Features() []types.Feature
	Validate()
	Test(*testing.T, *TestContextType) context.Context
}

// SerialSequence are feature tests that must be run in serial.
// These tests are run in the order of the array unless configured to be
// shuffled.
type SerialSequence struct {
	features []types.Feature
}

// NewSerialSequence creates a SerialSequence of feature tests to be executed
// by the test runner.
func NewSerialSequence(features ...types.Feature) *SerialSequence {
	return &SerialSequence{
		features: features,
	}
}

func (s *SerialSequence) SetFeatures(features ...types.Feature) {
	s.features = features
}

func (s *SerialSequence) Features() []types.Feature {
	return s.features
}

// Validate ensures that features are being run in accordance to label and
// configuration semantics.
// Invalid configuration is recorded as a source code bug and can be retrieved
// with FormatBugs.
func (s *SerialSequence) Validate() {
	// TODO(tvs): Validate testFeatures are being run in accordance to label
	// semantics
}

// Test is a wrapper function around [TestEnv.Test] that offers additional
// execution configuration.
// If TestContext.Shuffle is set, the tests will be shuffled before execution.
func (s *SerialSequence) Test(t *testing.T, tc *TestContextType) context.Context {
	if tc.Shuffle {
		rng := rand.New(rand.NewSource(tc.ShuffleSeed))
		rng.Shuffle(len(s.features), func(i, j int) { s.features[i], s.features[j] = s.features[j], s.features[i] })
	}

	return tc.TestEnv.Test(t, s.features...)
}

// ParallelSequence are feature tests that can be run in parallel when the
// framework is configured to run tests in parallel. If the framework is not
// configured to run in parallel, these are still run serially.
type ParallelSequence struct {
	features []types.Feature
}

// NewParallelSequence creates a ParallelSequence of feature tests to be
// executed by the test runner.
func NewParallelSequence(features ...types.Feature) *ParallelSequence {
	return &ParallelSequence{
		features: features,
	}
}

func (p *ParallelSequence) SetFeatures(features ...types.Feature) {
	p.features = features
}

func (p *ParallelSequence) Features() []types.Feature {
	return p.features
}

// Validate ensures that features are being run in accordance to label and
// configuration semantics.
// Tests labeled as disruptive or slow may not be run in parallel with other
// tests.
// Invalid configuration is recorded as a source code bug and can be retrieved
// with FormatBugs.
func (p *ParallelSequence) Validate() {
	for _, f := range p.features {
		if f.Labels().Contains(testlabels.Disruptive()) {
			RecordBug(NewBug(fmt.Sprintf("disruptive tests must not be run in parallel: %q", f.Name()), 1))
		}

		if f.Labels().Contains(testlabels.Slow()) {
			RecordBug(NewBug(fmt.Sprintf("slow tests must not be run in parallel: %q", f.Name()), 1))
		}
	}
}

// Test is a wrapper function for [TestEnv.TestInParallel] that can fall back
// to serial testing when parallelism is desired.
func (p *ParallelSequence) Test(t *testing.T, tc *TestContextType) context.Context {
	// TODO(tvs): Allow parallelism to be optional and the scope of parallelism
	// to be restricted.
	return tc.TestEnv.TestInParallel(t, p.features...)
}
