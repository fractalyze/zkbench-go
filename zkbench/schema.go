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

// Package zkbench provides a common JSON schema for cross-implementation
// benchmark results.
package zkbench

import (
	"encoding/json"
	"time"
)

// MetricValue represents a benchmark metric with optional confidence bounds.
type MetricValue struct {
	Value      float64  `json:"value"`
	Unit       string   `json:"unit"`
	LowerValue *float64 `json:"lower_value,omitempty"`
	UpperValue *float64 `json:"upper_value,omitempty"`
}

// TestVectors holds test vector verification information.
type TestVectors struct {
	InputHash  string `json:"input_hash"`
	OutputHash string `json:"output_hash"`
	Verified   bool   `json:"verified"`
}

// BenchmarkResult represents results from a single benchmark.
type BenchmarkResult struct {
	Latency     *MetricValue           `json:"latency,omitempty"`
	Throughput  *MetricValue           `json:"throughput,omitempty"`
	Memory      *MetricValue           `json:"memory,omitempty"`
	Iterations  int                    `json:"iterations,omitempty"`
	TestVectors *TestVectors           `json:"test_vectors,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Platform holds platform information.
type Platform struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	CPUCount  int    `json:"cpu_count"`
	CPUVendor string `json:"cpu_vendor,omitempty"`
	GPUVendor string `json:"gpu_vendor,omitempty"`
}

// Metadata holds benchmark metadata.
type Metadata struct {
	Implementation string   `json:"implementation"`
	Version        string   `json:"version"`
	CommitSHA      string   `json:"commit_sha"`
	Timestamp      string   `json:"timestamp"`
	Platform       Platform `json:"platform"`
}

// NewMetadata creates Metadata with auto-filled commit SHA, timestamp, and
// platform information.
func NewMetadata(implementation, version string) Metadata {
	return Metadata{
		Implementation: implementation,
		Version:        version,
		CommitSHA:      GetGitCommitSHA(),
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Platform:       GetPlatform(),
	}
}

// BenchmarkReport is the complete benchmark report.
type BenchmarkReport struct {
	Metadata   Metadata                   `json:"metadata"`
	Benchmarks map[string]BenchmarkResult `json:"benchmarks"`
}

// NewReport creates a BenchmarkReport with auto-filled metadata.
func NewReport(implementation, version string) *BenchmarkReport {
	return &BenchmarkReport{
		Metadata:   NewMetadata(implementation, version),
		Benchmarks: make(map[string]BenchmarkResult),
	}
}

// ToJSON serializes the report to indented JSON bytes.
func (r *BenchmarkReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
