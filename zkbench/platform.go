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
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// GetPlatform returns the current platform information including OS, arch,
// CPU count, and optional CPU/GPU vendor strings.
func GetPlatform() Platform {
	return Platform{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CPUCount:  runtime.NumCPU(),
		CPUVendor: getCPUVendor(),
		GPUVendor: getGPUVendor(),
	}
}

func getCPUVendor() string {
	switch runtime.GOOS {
	case "linux":
		return getCPUVendorLinux()
	case "darwin":
		return getCPUVendorMacOS()
	case "windows":
		return getCPUVendorWindows()
	}
	return ""
}

func getCPUVendorLinux() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`model name\s*:\s*(.+)`)
	match := re.FindSubmatch(data)
	if match != nil {
		return strings.TrimSpace(string(match[1]))
	}
	return ""
}

func getCPUVendorMacOS() string {
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getCPUVendorWindows() string {
	return os.Getenv("PROCESSOR_IDENTIFIER")
}

func getGPUVendor() string {
	if runtime.GOOS == "darwin" {
		return getGPUVendorMacOS()
	}
	// Try NVIDIA first, then ROCm.
	if v := getGPUVendorNvidia(); v != "" {
		return v
	}
	return getGPUVendorROCm()
}

func getGPUVendorNvidia() string {
	out, err := exec.Command(
		"nvidia-smi", "--query-gpu=name", "--format=csv,noheader",
	).Output()
	if err != nil {
		return ""
	}
	line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
	return strings.TrimSpace(line)
}

func getGPUVendorROCm() string {
	out, err := exec.Command("rocm-smi", "--showproductname").Output()
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`Card Series:\s*(.+)`)
	match := re.FindSubmatch(out)
	if match != nil {
		return strings.TrimSpace(string(match[1]))
	}
	return ""
}

func getGPUVendorMacOS() string {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`Chipset Model:\s*(.+)`)
	match := re.FindSubmatch(out)
	if match != nil {
		return strings.TrimSpace(string(match[1]))
	}
	return ""
}
