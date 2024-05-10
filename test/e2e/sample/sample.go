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

package sample

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
	"github.com/tvs/kubernetes-service-tests/test/e2e/framework/testlabels"
)

// Features returns a list of cloud provider test features to be run in a given context
func SerialFeatures(t *testing.T, tc *framework.TestContextType) []features.Feature {
	return []features.Feature{
		SerialFeature(t, tc),
	}
}

// Feature returns a test feature for samples.
// Sample tests are meant to be run against a plain kind cluster and validate
// the existence of simple components.
// TODO(tvs): Remove this sample testing functionality once we have real tests.
// TODO(tvs): Look into using testify for assertions, requirements, and
// expectations
func SerialFeature(t *testing.T, tc *framework.TestContextType) features.Feature {
	builder := features.New("sample")
	builder.WithLabel(testlabels.Sample())

	// Of course it's running, we can access its client!
	builder.Assess("Kube APIServer Pod Running",
		func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			var pods corev1.PodList

			err := c.Client().Resources("kube-system").List(ctx, &pods,
				resources.WithLabelSelector(
					labels.FormatLabels(
						map[string]string{"component": "kube-apiserver"},
					)))
			if err != nil {
				t.Errorf("unexpected error retrieving pods: %s", err)
			}

			if len(pods.Items) < 1 {
				t.Error("unable to find kube-apiserver")
			}

			for _, pod := range pods.Items {
				if pod.Status.Phase != "Running" {
					t.Error("pod is not running")
				}
			}

			return ctx
		})

	return builder.Feature()
}

func SampleTests(t *testing.T, tc *framework.TestContextType) {
	builder := framework.NewTestRunner()
	builder.WithSerialSequence(SerialFeatures(t, tc)...)
	// Sample parallel sequence
	builder.Runner().Test(t, tc)
}
