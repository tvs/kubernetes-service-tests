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

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/sample"
)

func RunE2ETests(t *testing.T, tc *framework.TestContextType) {
	// TODO(tvs): Ensure that tests are invoked for the correct context
	// e.g., service vs workload
	sample.SampleTests(t, tc)

	// TODO(tvs): Ensure Cluster and TKC tests can be run separately (or
	// conjoined) so we can pointedly limit which of the slow lifecycle tests we
	// exercise

	// TODO(tvs): Ensure upgrade tests can be exclusively targeted for
	// KubernetesDistribution testing

	// TODO(tvs): Chaos testing should be mutually exclusive from other forms
	// of testing
}
