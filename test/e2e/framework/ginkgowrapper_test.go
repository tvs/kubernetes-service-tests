/*
Copyright 2023 The Kubernetes Authors.

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
	"fmt"
	"testing"
)

func TestTagsEqual(t *testing.T) {
	testcases := []struct {
		a, b        interface{}
		expectEqual bool
	}{
		{1, 2, false},
		{2, 2, false},
		{WithSlow(), 2, false},
		{WithSlow(), WithSlow(), true},
		{WithLabel("hello"), WithLabel("world"), false},
		{WithLabel("hello"), WithLabel("hello"), true},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v=%v", tc.a, tc.b), func(t *testing.T) {
			actualEqual := TagsEqual(tc.a, tc.b)
			if actualEqual != tc.expectEqual {
				t.Errorf("expected %v, got %v", tc.expectEqual, actualEqual)
			}
		})
	}
}
