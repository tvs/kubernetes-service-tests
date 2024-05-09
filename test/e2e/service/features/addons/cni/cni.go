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

package cni

import (
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cni/antrea"
	"github.com/tvs/kubernetes-service-tests/test/e2e/service/features/addons/cni/calico"
)

// Features returns a list of CNI test features to be run in a given context
func Features(t *testing.T, tc *framework.TestContextType) []features.Feature {
	return []features.Feature{
		antrea.Feature(t, tc),
		calico.Feature(t, tc),
	}
}
