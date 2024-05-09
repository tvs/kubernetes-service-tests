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
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/types"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework/testlabels"
)

// TestContextType contains test settings and global state.
type TestContextType struct {
	// TestEnv is a global test environment for use with all test runners
	TestEnv env.Environment

	// shuffleFlag contains the contents of the command line flag that is
	// used to set the Shuffle boolean and the ShuffleSeed integer
	shuffleFlag string

	// Shuffle indicates that tests within a sequence should be shuffled.
	// This does not shuffle the order of sequences.
	Shuffle bool

	// ShuffleSeed is the seed used when setting up the RNG used for shuffling.
	ShuffleSeed int64

	// timeouts contains user-configurable timeouts for various operations.
	// Individual Framework instance also have such timeouts which may be
	// different from these here. To avoid confusion, this field is not
	// exported. Its values can be accessed through
	// NewTimeoutContext.
	timeouts TimeoutContext

	// versionFlag displays version information then exits.
	versionFlag bool
}

// Test is a wrapper around [TestEnv.Test] that ensures features are being
// run in accordance to label semantics.
// For example, it will ensure that disruptive tests are not run alongside
// other tests.
func (tc *TestContextType) Test(t *testing.T, testFeatures ...types.Feature) context.Context {
	// TODO(tvs): Validate testFeatures are being run in accordance to label
	// semantics
	return tc.TestEnv.Test(t, testFeatures...)
}

// TestInParallel is a wrapper around [TestEnv.TestInParallel] that ensures
// features are being run in accordance to label semantics.
// For example, it will ensure that disruptive or slow tests are not run
// in parallel.
func (tc *TestContextType) TestInParallel(t *testing.T, testFeatures ...types.Feature) context.Context {
	for _, f := range testFeatures {
		if f.Labels().Contains(testlabels.Disruptive()) {
			RecordBug(NewBug(fmt.Sprintf("disruptive tests must not be run in parallel: %q", f.Name()), 1))
		}

		if f.Labels().Contains(testlabels.Slow()) {
			RecordBug(NewBug(fmt.Sprintf("slow tests must not be run in parallel: %q", f.Name()), 1))
		}
	}

	// TODO(tvs): Validate testFeatures are being run in accordance to label
	// semantics
	return tc.TestEnv.TestInParallel(t, testFeatures...)
}

// TestContext should be used by all tests to access common context data.
var TestContext = TestContextType{
	timeouts: defaultTimeouts,
}

// RegisterCommonFlags registers flags common to all e2e test suites.
// The flag set can be flag.CommandLine (if desired) or a custom
// flag set that then gets passed to viperconfig.ViperizeFlags.
//
// The other Register*Flags methods below can be used to add more
// test-specific flags. However, those settings then get added
// regardless whether the test is actually in the test suite.
//
// For tests that have been converted to registering their
// options themselves, copy flags from test/e2e/framework/config
// as shown in HandleFlags.
func RegisterCommonFlags(flags *flag.FlagSet, tc *TestContextType) {
	flags.BoolVar(&tc.versionFlag, "version", false, "Displays version information")
	flags.StringVar(&tc.shuffleFlag, "shuffle", "off", "Shuffle tests within testing sequences. Valid values are 'off', 'on', or a valid integer that will be used as the RNG seed.")
}

// DefaultTestFlags establishes the common default flags that configure a
// TestContextType
func DefaultTestFlags(t *TestContextType) {
	RegisterCommonFlags(flag.CommandLine, t)
	RegisterTimeoutFlags(flag.CommandLine, t)
}

// processAndValidateFlags semantically validates flag values and processes
// flag behavior before tests start. This is used to configure the TestContext,
// as well as to report information about the tests themselves.
func processAndValidateFlags(t *TestContextType) {
	if t.versionFlag {
		fmt.Printf("%s\n", versionString())
		os.Exit(0)
	}

	if t.shuffleFlag == "off" {
		t.Shuffle = false
	} else {
		t.Shuffle = true
		var err error
		if t.shuffleFlag == "on" {
			t.ShuffleSeed = time.Now().UnixNano()
		} else {
			t.ShuffleSeed, err = strconv.ParseInt(t.shuffleFlag, 10, 64)
			if err != nil {
				log.Fatalf(`-shuffle should be "off", "on", or a valid integer: %s`, err)
			}
		}
	}
	// TODO(tvs): Log shuffle seed
}

// AfterReadingAllFlags makes changes to the context after all flags
// have been read and prepares the process for a test run.
func AfterReadingAllFlags(t *TestContextType) {
	processAndValidateFlags(t)

	cfg, err := envconf.NewFromFlags()
	if err != nil {
		log.Fatalf("failed to build envconf from flags: %s", err)
	}

	cfg.WithKubeconfigFile(conf.ResolveKubeConfigFile())

	t.TestEnv = env.NewWithConfig(cfg)
}
