package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// Job represents a temperature conversion task.
type Job struct {
	Value float64 `json:"value"`
	From  Unit    `json:"from"`
	To    Unit    `json:"to"`
}

// Validate ensures the job contains supported units.
func (j Job) Validate() error {
	if err := j.From.Validate(); err != nil {
		return fmt.Errorf("invalid source unit: %w", err)
	}
	if err := j.To.Validate(); err != nil {
		return fmt.Errorf("invalid target unit: %w", err)
	}
	return nil
}

// Result captures the outcome of processing a Job.
type Result struct {
	Err       error
	Job       Job
	Elapsed   time.Duration
	Converted float64
	Index     int
}

// LoadJobsFromReader decodes jobs from JSON.
func LoadJobsFromReader(ctx context.Context, r io.Reader) ([]Job, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var jobs []Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if err := job.Validate(); err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

// LoadJobsFromFile reads a JSON file and returns the jobs.
func LoadJobsFromFile(ctx context.Context, filename string) (jobs []Job, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return LoadJobsFromReader(ctx, f)
}

// ProcessSequential converts the jobs one by one.
func ProcessSequential(ctx context.Context, jobs []Job) ([]Result, error) {
	results := make([]Result, 0, len(jobs))
	for idx, job := range jobs {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		start := time.Now()
		converted, err := Convert(job.Value, job.From, job.To)
		results = append(results, Result{
			Job:       job,
			Converted: converted,
			Err:       err,
			Index:     idx,
			Elapsed:   time.Since(start),
		})
	}

	return results, nil
}

type jobRequest struct {
	index int
	job   Job
}

func processJob(ctx context.Context, req jobRequest) (Result, bool) {
	select {
	case <-ctx.Done():
		return Result{}, false
	default:
	}

	start := time.Now()
	converted, err := Convert(req.job.Value, req.job.From, req.job.To)

	return Result{
		Job:       req.job,
		Converted: converted,
		Err:       err,
		Index:     req.index,
		Elapsed:   time.Since(start),
	}, true
}

func startWorkers(ctx context.Context, workerCount int, jobCh <-chan jobRequest, resultCh chan<- Result) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range jobCh {
				result, ok := processJob(ctx, req)
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case resultCh <- result:
				}
			}
		}()
	}
	return &wg
}

// ProcessConcurrent converts the jobs using worker goroutines and channels.
func ProcessConcurrent(ctx context.Context, jobs []Job, workerCount, bufferSize int) ([]Result, error) {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	if bufferSize < 0 {
		bufferSize = 0
	}

	jobCh := make(chan jobRequest)
	resultCh := make(chan Result, bufferSize)

	wg := startWorkers(ctx, workerCount, jobCh, resultCh)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	go func() {
		defer close(jobCh)
		for idx, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- jobRequest{index: idx, job: job}:
			}
		}
	}()

	return collectResults(ctx, resultCh, len(jobs))
}

func collectResults(ctx context.Context, resultCh <-chan Result, total int) ([]Result, error) {
	results := make([]Result, total)
	collected := 0
	for {
		select {
		case <-ctx.Done():
			return results[:collected], ctx.Err()
		case result, ok := <-resultCh:
			if !ok {
				if ctx.Err() != nil {
					return results[:collected], ctx.Err()
				}
				return results[:collected], nil
			}
			results[result.Index] = result
			collected++
			if collected == total {
				return results, nil
			}
		}
	}
}
