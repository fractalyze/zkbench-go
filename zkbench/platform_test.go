// Copyright 2026 The zkbench-go Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ==============================================================================

package zkbench

import (
	"runtime"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	p := GetPlatform()
	if p.OS != runtime.GOOS {
		t.Errorf("OS should be %q, got %q", runtime.GOOS, p.OS)
	}
	if p.Arch != runtime.GOARCH {
		t.Errorf("Arch should be %q, got %q", runtime.GOARCH, p.Arch)
	}
	if p.CPUCount < 1 {
		t.Errorf("CPUCount should be >= 1, got %d", p.CPUCount)
	}
}

func TestGetCPUVendor(t *testing.T) {
	v := getCPUVendor()
	// On CI/Linux this should return something. On unsupported OS it may
	// return empty, which is acceptable.
	if runtime.GOOS == "linux" && v == "" {
		t.Log("warning: getCPUVendor returned empty on linux")
	}
}
