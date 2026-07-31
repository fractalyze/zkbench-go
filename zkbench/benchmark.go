// Copyright 2026 zkbench Authors
// SPDX-License-Identifier: Apache-2.0

// Template-based benchmark harness for Go benchmarks.
//
// Mirrors the JaxBenchmark pattern from zkbench-py: implement the Benchmark
// interface (GetConfig + GetOps), then call Run() to handle CLI, warmup,
// timing, statistics, and JSON output.
//
// Example:
//
//	type MyBench struct{}
//	func (b MyBench) GetConfig() BenchmarkConfig { ... }
//	func (b MyBench) GetOps(sizes []int) []BenchmarkOp { ... }
//	func main() { os.Exit(zkbench.Run(MyBench{})) }
package zkbench

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BenchmarkConfig defines a benchmark suite's metadata and defaults.
type BenchmarkConfig struct {
	Implementation string
	Version        string
	Iterations     int // default measured iterations
	Warmup         int // default warmup iterations
}

// BenchmarkOp encapsulates a single benchmarkable operation.
type BenchmarkOp struct {
	Name     string         // benchmark key (e.g., "fft/16")
	Fn       func()         // the timed operation
	Setup    func()         // per-iteration setup (optional)
	Sync     func()         // post-op sync (optional, e.g., cuda_device_sync)
	Metadata map[string]any // extra metadata (field, degree, etc.)

	InputHash    string        // pre-computed input hash
	OutputHashFn func() string // returns output hash after timing (optional)
	VerifyFn     func() bool   // returns true if output is correct (optional)
}

// Benchmark is the interface that benchmarks implement.
type Benchmark interface {
	GetConfig() BenchmarkConfig
	GetOps(sizes []int) []BenchmarkOp
}

// Run orchestrates the benchmark lifecycle: CLI → warmup → timing → stats → JSON.
// Returns 0 on success, 1 on verification failure.
func Run(b Benchmark) int {
	config := b.GetConfig()

	iterationsFlag := flag.Int("iterations", config.Iterations, "Measured iterations")
	warmupFlag := flag.Int("warmup", config.Warmup, "Warmup iterations")
	sizesFlag := flag.String("sizes", "16,18,20,22", "Comma-separated log2 sizes")
	outputFlag := flag.String("output", "", "Output file (default: stdout)")
	flag.Parse()

	sizes := parseSizes(*sizesFlag)

	fmt.Fprintf(os.Stderr, "%s v%s — %d warmup + %d measured iterations\n",
		config.Implementation, config.Version, *warmupFlag, *iterationsFlag)

	ops := b.GetOps(sizes)
	report := NewReport(config.Implementation, config.Version)

	allVerified := true
	for i := range ops {
		op := &ops[i]
		fmt.Fprintf(os.Stderr, "  %s...\n", op.Name)
		result := runSingleOp(op, *iterationsFlag, *warmupFlag)

		if result.TestVectors != nil && !result.TestVectors.Verified {
			fmt.Fprintf(os.Stderr, "  VERIFICATION FAILED: %s\n", op.Name)
			allVerified = false
		}
		if result.Latency != nil {
			fmt.Fprintf(os.Stderr, "    median: %.1f us\n", result.Latency.Value)
		}

		report.Benchmarks[op.Name] = result
	}

	data, err := report.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		return 1
	}

	if *outputFlag != "" {
		if err := os.WriteFile(*outputFlag, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *outputFlag, err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Wrote %s\n", *outputFlag)
	} else {
		fmt.Println(string(data))
	}

	if !allVerified {
		return 1
	}
	return 0
}

func runSingleOp(op *BenchmarkOp, iterations, warmup int) BenchmarkResult {
	callSetup := func() {
		if op.Setup != nil {
			op.Setup()
		}
	}
	callSync := func() {
		if op.Sync != nil {
			op.Sync()
		}
	}

	// Warmup
	for range warmup {
		callSetup()
		callSync()
		op.Fn()
		callSync()
	}

	// Measured iterations (collect in μs directly)
	timesUs := make([]float64, iterations)
	for i := range iterations {
		callSetup()
		callSync()
		start := time.Now()
		op.Fn()
		callSync()
		timesUs[i] = float64(time.Since(start).Nanoseconds()) / 1000.0
	}

	result := BenchmarkResult{
		Iterations: iterations,
	}

	if iterations > 0 {
		mean, stdev, _ := CalculateStatistics(timesUs)
		med, _ := Median(timesUs)
		lower, upper := CalculateConfidenceInterval(mean, stdev, 0.95)

		lower = roundTo2(lower)
		upper = roundTo2(upper)

		result.Latency = &MetricValue{
			Value:      roundTo2(med),
			Unit:       "us",
			LowerValue: &lower,
			UpperValue: &upper,
		}
	}

	if len(op.Metadata) > 0 {
		result.Metadata = op.Metadata
	}

	if op.OutputHashFn != nil || op.VerifyFn != nil {
		outputHash := ""
		if op.OutputHashFn != nil {
			outputHash = op.OutputHashFn()
		}
		verified := true
		if op.VerifyFn != nil {
			verified = op.VerifyFn()
		}
		result.TestVectors = &TestVectors{
			InputHash:  op.InputHash,
			OutputHash: outputHash,
			Verified:   verified,
		}
	}

	return result
}

func roundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}

func parseSizes(s string) []int {
	parts := strings.Split(s, ",")
	sizes := make([]int, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid size: %s\n", p)
			os.Exit(1)
		}
		sizes[i] = v
	}
	sort.Ints(sizes)
	return sizes
}
