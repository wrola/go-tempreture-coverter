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
	Job       Job
	Converted float64
	Err       error
	Index     int
	Elapsed   time.Duration
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
func LoadJobsFromFile(ctx context.Context, filename string) ([]Job, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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

	var workerWG sync.WaitGroup
	for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			for workRequest := range jobCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				start := time.Now()
				converted, err := Convert(workRequest.job.Value, workRequest.job.From, workRequest.job.To)

				result := Result{
					Job:       workRequest.job,
					Converted: converted,
					Err:       err,
					Index:     workRequest.index,
					Elapsed:   time.Since(start),
				}

				select {
				case <-ctx.Done():
					return
				case resultCh <- result:
				}
			}
		}()
	}

	go func() {
		workerWG.Wait()
		close(resultCh)
	}()

	go func() {
		defer close(jobCh)
		for jobIndex, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- jobRequest{index: jobIndex, job: job}:
			}
		}
	}()

	results := make([]Result, len(jobs))
	collected := 0
	for {
		select {
		case <-ctx.Done():
			return results[:collected], ctx.Err()
		case result, ok := <-resultCh:
			if !ok {
				return results[:collected], nil
			}
			results[result.Index] = result
			collected++
			if collected == len(jobs) {
				return results, nil
			}
		}
	}
}
