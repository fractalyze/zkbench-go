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
	"encoding/json"
	"testing"
)

func TestNewMetadata(t *testing.T) {
	m := NewMetadata("test-impl", "1.0.0")
	if m.Implementation != "test-impl" {
		t.Errorf("got implementation %q, want %q", m.Implementation, "test-impl")
	}
	if m.Version != "1.0.0" {
		t.Errorf("got version %q, want %q", m.Version, "1.0.0")
	}
	if m.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
	if m.Platform.OS == "" {
		t.Error("platform OS should not be empty")
	}
	if m.Platform.CPUCount < 1 {
		t.Errorf("cpu_count should be >= 1, got %d", m.Platform.CPUCount)
	}
}

func TestNewReport(t *testing.T) {
	r := NewReport("test", "0.1.0")
	if r.Benchmarks == nil {
		t.Fatal("benchmarks map should be initialized")
	}
	if len(r.Benchmarks) != 0 {
		t.Error("benchmarks should be empty initially")
	}
}

func TestBenchmarkReportToJSON(t *testing.T) {
	r := NewReport("test", "0.1.0")
	lower, upper := 90.0, 110.0
	r.Benchmarks["test_bench"] = BenchmarkResult{
		Latency: &MetricValue{
			Value:      100.0,
			Unit:       "ms",
			LowerValue: &lower,
			UpperValue: &upper,
		},
		Iterations: 5,
		TestVectors: &TestVectors{
			InputHash:  "abc123",
			OutputHash: "def456",
			Verified:   true,
		},
	}

	data, err := r.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON and round-trips.
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	benchmarks, ok := parsed["benchmarks"].(map[string]interface{})
	if !ok {
		t.Fatal("benchmarks should be a JSON object")
	}
	bench, ok := benchmarks["test_bench"].(map[string]interface{})
	if !ok {
		t.Fatal("test_bench should be a JSON object")
	}
	latency, ok := bench["latency"].(map[string]interface{})
	if !ok {
		t.Fatal("latency should be a JSON object")
	}
	if latency["value"].(float64) != 100.0 {
		t.Errorf("latency value should be 100.0, got %v", latency["value"])
	}
	if latency["lower_value"].(float64) != 90.0 {
		t.Errorf("lower_value should be 90.0, got %v", latency["lower_value"])
	}
}

func TestBenchmarkResultOmitsNilFields(t *testing.T) {
	r := NewReport("test", "0.1.0")
	r.Benchmarks["empty"] = BenchmarkResult{}

	data, err := r.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)
	benchmarks := parsed["benchmarks"].(map[string]interface{})
	bench := benchmarks["empty"].(map[string]interface{})

	if _, ok := bench["latency"]; ok {
		t.Error("nil latency should be omitted")
	}
	if _, ok := bench["throughput"]; ok {
		t.Error("nil throughput should be omitted")
	}
	if _, ok := bench["test_vectors"]; ok {
		t.Error("nil test_vectors should be omitted")
	}
}
