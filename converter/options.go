package converter

import (
	"runtime"
	"time"
)

type ProcessOptions struct {
	Workers    int
	BufferSize int
	Timeout    time.Duration
}

type ProcessOption func(*ProcessOptions)

func WithWorkers(count int) ProcessOption {
	return func(opts *ProcessOptions) {
		opts.Workers = count
	}
}

func WithBuffer(size int) ProcessOption {
	return func(opts *ProcessOptions) {
		opts.BufferSize = size
	}
}

func WithTimeout(timeout time.Duration) ProcessOption {
	return func(opts *ProcessOptions) {
		opts.Timeout = timeout
	}
}

func buildOptions(opts ...ProcessOption) *ProcessOptions {
	options := &ProcessOptions{
		Workers:    runtime.NumCPU(),
		BufferSize: 0,
		Timeout:    0,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Workers <= 0 {
		options.Workers = runtime.NumCPU()
	}
	if options.BufferSize < 0 {
		options.BufferSize = 0
	}

	return options
}
