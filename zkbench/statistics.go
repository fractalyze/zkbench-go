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
	"fmt"
	"math"
	"sort"
)

// CalculateStatistics returns the mean and sample standard deviation (n-1) of
// the given values.
func CalculateStatistics(values []float64) (mean, stdev float64, err error) {
	if len(values) == 0 {
		return 0, 0, fmt.Errorf("cannot calculate statistics on empty slice")
	}

	n := float64(len(values))
	for _, v := range values {
		mean += v
	}
	mean /= n

	if len(values) < 2 {
		return mean, 0, nil
	}

	var variance float64
	for _, v := range values {
		d := v - mean
		variance += d * d
	}
	variance /= n - 1
	stdev = math.Sqrt(variance)

	return mean, stdev, nil
}

// CalculateConfidenceInterval returns (lower, upper) bounds for the given
// confidence level. Uses z ≈ 2.0 for 95% and z ≈ 2.576 for 99%.
func CalculateConfidenceInterval(mean, stdev, confidence float64) (lower, upper float64) {
	z := 2.0
	if confidence == 0.99 {
		z = 2.576
	}
	margin := z * stdev
	return mean - margin, mean + margin
}

// Median returns the median of a slice of float64 values. The input slice is
// not modified.
func Median(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("cannot calculate median of empty slice")
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	return sorted[len(sorted)/2], nil
}
