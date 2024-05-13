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

package tanzukubernetescluster

import (
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cloudprovider"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cni"
)

func TanzuKubernetesClusterTests(t *testing.T, tc *framework.TestContextType) {
	builder := framework.NewTestRunner()

	// TODO(tvs): TanzuKubernetesCluster creation test
	//builder.WithSerialSequence(CreateClusterTests(t, tc))

	//builder.WithParallelSequence(CAPIResourceTests(t, tc)...)

	// Run feature tests in parallel
	feat := []features.Feature{}
	feat = append(feat, cni.Features(t, tc)...)
	feat = append(feat, cloudprovider.Features(t, tc)...)
	builder.WithParallelSequence(feat...)

	// TODO(tvs): Cluster scale test
	//builder.WithSerialSequence(ScaleClusterTests(t, tc))

	// TODO(tvs): Cluster remediation test
	//builder.WithSerialSequence(RemediateClusterTests(t, tc))

	// TODO(tvs): Cluster upgrade test

	// TODO(tvs): Cluster deletion test

	builder.Runner().Test(t, tc)
}
