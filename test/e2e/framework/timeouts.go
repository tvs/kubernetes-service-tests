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
	"flag"
	"time"
)

var defaultTimeouts = TimeoutContext{
	Poll:            2 * time.Second,
	PodStart:        5 * time.Minute,
	PodDelete:       5 * time.Minute,
	NodeSchedulable: 30 * time.Minute,
	ClusterReady:    30 * time.Minute,
}

// TimeoutContext contains timeout settings for several actions.
type TimeoutContext struct {
	// Poll is how long to wait between API calls when waiting for a condition.
	Poll time.Duration

	// PodStart is how long to wait for the pod to be started.
	// This value is the default for gomega.Eventually
	PodStart time.Duration

	// PodDelete is how long to wait for the pod to be deleted.
	PodDelete time.Duration

	// NodeSchedulable is how long to wait for a/all nodes to be schedulable.
	NodeSchedulable time.Duration

	// ClusterReady is how long to wait for a Cluster be ready.
	ClusterReady time.Duration
}

// RegisterTimeoutFlags registers flags related to timeouts
func RegisterTimeoutFlags(flags *flag.FlagSet) {
	flags.DurationVar(&TestContext.timeouts.Poll, "polling-interval", TestContext.timeouts.Poll, "Interval between API calls when waiting for a condition.")
	flags.DurationVar(&TestContext.timeouts.PodStart, "pod-start-timeout", TestContext.timeouts.PodStart, "Timeout for waiting for a pod to be started.")
	flags.DurationVar(&TestContext.timeouts.PodDelete, "pod-delete-timeout", TestContext.timeouts.PodDelete, "Timeout for waiting for a pod to be deleted.")
	flags.DurationVar(&TestContext.timeouts.NodeSchedulable, "node-schedulable-timeout", TestContext.timeouts.NodeSchedulable, "Timeout for waiting for a/all nodes to be schedulable.")
	flags.DurationVar(&TestContext.timeouts.ClusterReady, "cluster-ready-timeout", TestContext.timeouts.ClusterReady, "Timeout for waiting for a cluster to be ready.")
}

// NewTimeoutContext returns a TimeoutContext with all values set either to
// hard-coded defaults or a value that was configured when running the E2E
// suite. Should be called after command line parsing.
func NewTimeoutContext() *TimeoutContext {
	// Make a copy, otherwise the caller would have the ability to modify
	// the original values.
	copy := TestContext.timeouts
	return &copy
}

// PollInterval defines how long to wait between API server queries while
// waiting for some condition.
//
// This value is the default for gomega.Eventually and gomega.Consistently.
func PollInterval() time.Duration {
	return TestContext.timeouts.Poll
}
