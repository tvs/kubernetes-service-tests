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

package testlabels

const (
	kindLabelKey               = "kind"
	kindKubernetesService      = "KubernetesService"
	kindKubernetesDistribution = "KubernetesDistribution"
	kindWorkloadCluster        = "WorkloadCluster"
	kindSample                 = "Sample"
)

func kindLabel(v string) (string, string) {
	return kindLabelKey, v
}

// KubernetesService specifies that a certain test or group of tests are
// targeted to the Kubernetes Service. The return value must be passed into
// [features.WithLabel].
func KubernetesService() (string, string) {
	return kindLabel(kindKubernetesService)
}

// KubernetesDistribution specifies that a certain test or group of tests are
// targeted to the Kubernetes Distribution. The return value must be passed
// into [features.WithLabel].
func KubernetesDistribution() (string, string) {
	return kindLabel(kindKubernetesService)
}

// WorkloadCluster specifies that a certain test or group of tests are targeted
// to a workload cluster provisioned by the Kubernetes Service. The return
// value must be passed into [features.WithLabel].
func WorkloadCluster() (string, string) {
	return kindLabel(kindWorkloadCluster)
}

// Sample specifies that a certain test or group of tests are targeted to a
// cluster being used for example purposes. The return value must be passed
// into [features.WithLabel].
//
// TODO(tvs): Remove this once we have legitimate testing
func Sample() (string, string) {
	return kindLabel(kindSample)
}

const (
	environmentLabelKey = "environment"
	environmentLinux    = "Linux"
	environmentWindows  = "Windows"
)

func environmentLabel(v string) (string, string) {
	return environmentLabelKey, v
}

// Linux specifies that a certain test or group of tests only work in a Linux
// environment. The return value must be passed into [features.WithLabel].
func Linux() (string, string) {
	return environmentLabel(environmentLinux)
}

// Windows specifies that a certain test or group of tests only work in a
// Windows environment. The return value must be passed into
// [features.WithLabel].
func Windows() (string, string) {
	return environmentLabel(environmentWindows)
}

const (
	typeLabelKey    = "type"
	typeConformance = "Conformance"
	typeFlaky       = "Flaky"
	typeDisruptive  = "Disruptive"
	typeSlow        = "Slow"
)

func typeLabel(v string) (string, string) {
	return typeLabelKey, v
}

// Conformance specifies that a certain test or group of tests must pass
// regardless of feature specification or system configuration. The return
// value must be passed into [features.WithLabel].
func Conformance() (string, string) {
	return typeLabel(typeConformance)
}

// Flaky specifies that a certain test or group of tests are failing randomly.
// These tests are usually filtered out and ran separately from other tests.
// The return value must be passed into [features.WithLabel].
func Flaky() (string, string) {
	return typeLabel(typeFlaky)
}

// Disruptive specifies that a certain test or group of tests temporarily
// affect the functionality of the Kubernetes Service or Workload Cluster.
// The return value must be passed into [features.WithLabel].
func Disruptive() (string, string) {
	return typeLabel(typeDisruptive)
}

// Slow specifies that a certain test or group of tests must not run in
// parallel with other tests. The return value must be passed into
// [features.WithLabel].
func Slow() (string, string) {
	return typeLabel(typeSlow)
}
