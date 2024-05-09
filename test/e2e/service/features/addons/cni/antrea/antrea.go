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

package antrea

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/framework/testlabels"
)

// Feature returns a test feature for the Antrea CNI
func Feature(t *testing.T, tc *framework.TestContextType) features.Feature {
	builder := features.New("antrea")
	builder.WithLabel(testlabels.KubernetesService())

	builder.Assess("ClusterBootstrap reports success", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
		// TODO(tvs): Add test content
		return ctx
	})

	return builder.Feature()
}
