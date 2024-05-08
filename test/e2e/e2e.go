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

package e2e

import (
	"testing"

	"github.com/onsi/ginkgo/v2"

	"sigs.k8s.io/e2e-framework/pkg/env"
)

// RunE2ETests runs the E2E tests using the Ginkgo runner
// This function is called on each Ginkgo node in parallel mode.
func RunE2ETests(t *testing.T, testenv env.Environment) {
	ginkgo.RunSpecs(t, "Kubernetes Service e2e suite")
}
