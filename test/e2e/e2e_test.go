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
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"

	"github.com/tvs/kubernetes-service-tests/test/e2e/framework"
)

func handleFlags() {
	framework.RegisterTimeoutFlags(flag.CommandLine)
}

func TestMain(m *testing.M) {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Displays version information.")

	handleFlags()

	cfg, err := envconf.NewFromFlags()
	if err != nil {
		log.Fatalf("failed to build envconf from flags: %s", err)
	}

	if versionFlag {
		fmt.Printf("%s\n", versionString())
		os.Exit(0)
	}

	framework.AfterReadingAllFlags(&framework.TestContext, cfg)

	// Generate a namespace for the e2e test and ensure it's cleaned up
	namespace := envconf.RandomName("k8s-svc-e2e", 16)
	framework.TestContext.TestEnv.Setup(
		envfuncs.CreateNamespace(namespace),
	)
	framework.TestContext.TestEnv.Finish(
		envfuncs.DeleteNamespace(namespace),
	)

	os.Exit(framework.TestContext.TestEnv.Run(m))
}
