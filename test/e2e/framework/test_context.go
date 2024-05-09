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
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

// TestContextType contains test settings and global state.
type TestContextType struct {
	// TestEnv is a global test environment for use with all test runners
	TestEnv env.Environment

	// timeouts contains user-configurable timeouts for various operations.
	// Individual Framework instance also have such timeouts which may be
	// different from these here. To avoid confusion, this field is not
	// exported. Its values can be accessed through
	// NewTimeoutContext.
	timeouts TimeoutContext
}

// TestContext should be used by all tests to access common context data.
var TestContext = TestContextType{
	timeouts: defaultTimeouts,
}

// AfterReadingAllFlags makes changes to the context after all flags
// have been read and prepares the process for a test run.
func AfterReadingAllFlags(t *TestContextType, cfg *envconf.Config) {
	cfg.WithKubeconfigFile(conf.ResolveKubeConfigFile())

	t.TestEnv = env.NewWithConfig(cfg)
}
