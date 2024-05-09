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

package workload

import (
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/workload/features/addons/cloudprovider"
	"github.com/tvs/kubernetes-service-tests/test/e2e/workload/features/addons/cni"
)

func WorkloadClusterTests(t *testing.T, tc *framework.TestContextType) {
	feat := []features.Feature{}
	feat = append(feat, cni.Features(t, tc)...)
	feat = append(feat, cloudprovider.Features(t, tc)...)

	framework.TestContext.TestEnv.Test(t, feat...)
}
