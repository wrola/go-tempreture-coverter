package converter

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestLoadJobsFromReader(t *testing.T) {
	t.Parallel()

	payload := []byte(`[
		{ "value": 36.6, "from": "c", "to": "f" },
		{ "value": 100, "from": "F", "to": "c" }
	]`)

	jobs, err := LoadJobsFromReader(context.Background(), bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("LoadJobsFromReader returned error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].From != UnitCelsius || jobs[0].To != UnitFahrenheit {
		t.Fatalf("unexpected unit normalization: %+v", jobs[0])
	}
}

func TestSequentialAndConcurrentProduceSameResults(t *testing.T) {
	t.Parallel()

	jobs := []Job{
		{Value: 0, From: UnitCelsius, To: UnitKelvin},
		{Value: 100, From: UnitFahrenheit, To: UnitCelsius},
		{Value: 273.15, From: UnitKelvin, To: UnitFahrenheit},
	}

	ctx := context.Background()

	sequential, err := ProcessSequential(ctx, jobs)
	if err != nil {
		t.Fatalf("ProcessSequential error: %v", err)
	}

	concurrent, err := ProcessConcurrent(ctx, jobs, WithWorkers(2), WithBuffer(2))
	if err != nil {
		t.Fatalf("ProcessConcurrent error: %v", err)
	}

	for jobIndex := range jobs {
		compareResult(t, sequential[jobIndex], concurrent[jobIndex])
	}
}

func TestProcessSequentialHonorsContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results, err := ProcessSequential(ctx, []Job{{Value: 0, From: UnitCelsius, To: UnitFahrenheit}})
	if err == nil {
		t.Fatalf("expected context cancellation error, got nil")
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestProcessConcurrentHonorsContext(t *testing.T) {
	t.Parallel()

	jobs := makeTestJobs(100)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results, err := ProcessConcurrent(ctx, jobs, WithWorkers(4), WithBuffer(4))
	if err == nil {
		t.Fatalf("expected context cancellation error, got nil")
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func compareResult(t *testing.T, firstResult, secondResult Result) {
	t.Helper()

	if firstResult.Err != nil || secondResult.Err != nil {
		if fmt.Sprint(firstResult.Err) != fmt.Sprint(secondResult.Err) {
			t.Fatalf("errors differ: %v vs %v", firstResult.Err, secondResult.Err)
		}
		return
	}

	if firstResult.Converted != secondResult.Converted {
		t.Fatalf("converted mismatch: %v vs %v", firstResult.Converted, secondResult.Converted)
	}
	if firstResult.Job.From != secondResult.Job.From || firstResult.Job.To != secondResult.Job.To || firstResult.Job.Value != secondResult.Job.Value {
		t.Fatalf("job mismatch: %+v vs %+v", firstResult.Job, secondResult.Job)
	}
}

func makeTestJobs(jobCount int) []Job {
	jobs := make([]Job, jobCount)
	for jobIndex := range jobs {
		value := float64(jobIndex%212) - 100 // mix negatives and positives
		jobs[jobIndex] = Job{
			Value: value,
			From:  UnitCelsius,
			To:    UnitFahrenheit,
		}
	}
	return jobs
}

func BenchmarkProcessSequential(b *testing.B) {
	ctx := context.Background()
	jobs := makeTestJobs(2048)

	// Warm-up to avoid measuring initial allocations.
	if _, err := ProcessSequential(ctx, jobs); err != nil {
		b.Fatalf("warm-up failed: %v", err)
	}

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		results, err := ProcessSequential(ctx, jobs)
		if err != nil {
			b.Fatalf("ProcessSequential error: %v", err)
		}
		if len(results) != len(jobs) {
			b.Fatalf("expected %d results, got %d", len(jobs), len(results))
		}
	}
}

func BenchmarkProcessConcurrent(b *testing.B) {
	ctx := context.Background()
	jobs := makeTestJobs(2048)

	benchmarks := []struct {
		name    string
		workers int
		buffer  int
	}{
		{name: "workers=1_buffer=0", workers: 1, buffer: 0},
		{name: "workers=cpu_buffer=0", workers: runtime.NumCPU(), buffer: 0},
		{name: "workers=cpu_buffer=cpu", workers: runtime.NumCPU(), buffer: runtime.NumCPU()},
	}

	for _, benchmarkCase := range benchmarks {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			// Warm-up
			if _, err := ProcessConcurrent(ctx, jobs, WithWorkers(benchmarkCase.workers), WithBuffer(benchmarkCase.buffer)); err != nil {
				b.Fatalf("warm-up failed: %v", err)
			}

			b.ResetTimer()
			for iteration := 0; iteration < b.N; iteration++ {
				results, err := ProcessConcurrent(ctx, jobs, WithWorkers(benchmarkCase.workers), WithBuffer(benchmarkCase.buffer))
				if err != nil {
					b.Fatalf("ProcessConcurrent error: %v", err)
				}
				if len(results) != len(jobs) {
					b.Fatalf("expected %d results, got %d", len(jobs), len(results))
				}
			}
		})
	}
}

func BenchmarkProcessConcurrentWithTimeout(b *testing.B) {
	jobs := makeTestJobs(2048)
	timeout := 1 * time.Nanosecond

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		_, _ = ProcessConcurrent(timeoutCtx, jobs, WithWorkers(runtime.NumCPU()), WithBuffer(runtime.NumCPU()))
		cancel()
	}
}
