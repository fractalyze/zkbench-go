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
	"math"
	"testing"
)

func TestCalculateStatisticsEmpty(t *testing.T) {
	_, _, err := CalculateStatistics(nil)
	if err == nil {
		t.Error("expected error for empty slice")
	}
}

func TestCalculateStatisticsSingle(t *testing.T) {
	mean, stdev, err := CalculateStatistics([]float64{42.0})
	if err != nil {
		t.Fatal(err)
	}
	if mean != 42.0 {
		t.Errorf("mean should be 42.0, got %f", mean)
	}
	if stdev != 0.0 {
		t.Errorf("stdev should be 0.0 for single value, got %f", stdev)
	}
}

func TestCalculateStatisticsMultiple(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, stdev, err := CalculateStatistics(values)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(mean-5.0) > 0.001 {
		t.Errorf("mean should be 5.0, got %f", mean)
	}
	// Sample stdev of {2,4,4,4,5,5,7,9} = sqrt(32/7) ≈ 2.138
	if math.Abs(stdev-2.138) > 0.01 {
		t.Errorf("stdev should be ~2.138, got %f", stdev)
	}
}

func TestCalculateConfidenceInterval95(t *testing.T) {
	lower, upper := CalculateConfidenceInterval(100.0, 10.0, 0.95)
	if math.Abs(lower-80.0) > 0.001 {
		t.Errorf("lower should be 80.0, got %f", lower)
	}
	if math.Abs(upper-120.0) > 0.001 {
		t.Errorf("upper should be 120.0, got %f", upper)
	}
}

func TestCalculateConfidenceInterval99(t *testing.T) {
	lower, upper := CalculateConfidenceInterval(100.0, 10.0, 0.99)
	if math.Abs(lower-74.24) > 0.01 {
		t.Errorf("lower should be ~74.24, got %f", lower)
	}
	if math.Abs(upper-125.76) > 0.01 {
		t.Errorf("upper should be ~125.76, got %f", upper)
	}
}

func TestMedian(t *testing.T) {
	m, err := Median([]float64{3, 1, 4, 1, 5})
	if err != nil {
		t.Fatal(err)
	}
	// sorted: [1, 1, 3, 4, 5], index 2 → 3
	if m != 3.0 {
		t.Errorf("median should be 3.0, got %f", m)
	}
}

func TestMedianEmpty(t *testing.T) {
	_, err := Median(nil)
	if err == nil {
		t.Error("expected error for empty slice")
	}
}
