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

package cluster

import (
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cloudprovider"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cni"
)

func ClusterTests(t *testing.T, tc *framework.TestContextType) {
	// TODO(tvs): Cluster creation test

	// Run feature tests in parallel
	feat := []features.Feature{}
	feat = append(feat, cni.Features(t, tc)...)
	feat = append(feat, cloudprovider.Features(t, tc)...)

	// TODO(tvs): Ensure features do not have the slow or serialized label when
	// run in parallel
	framework.TestContext.TestEnv.TestInParallel(t, feat...)

	// TODO(tvs): Cluster upgrade test

	// TODO(tvs):

}
