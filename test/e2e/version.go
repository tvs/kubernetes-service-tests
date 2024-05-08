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
	"fmt"
)

// var needs to be used instead of const as ldflags are used to fill this
// information during the build process
var (
	serviceVersion = "unknown"
	serviceVendor  = "unknown"
	goos           = "unknown"
	goarch         = "unknown"
	gitCommit      = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y)-%m-%dT%H:%M:%SZ'
)

type version struct {
	ServiceVersion string `json:"serviceVersion"`
	ServiceVendor  string `json:"serviceVendor"`
	GitCommit      string `json:"gitCommit"`
	BuildDate      string `json:"buildDate"`
	GoOs           string `json:"goOs"`
	GoArch         string `json:"goArch"`
}

func versionString() string {
	return fmt.Sprintf("Version: %#v", version{
		serviceVersion,
		serviceVendor,
		gitCommit,
		buildDate,
		goos,
		goarch,
	})
}
