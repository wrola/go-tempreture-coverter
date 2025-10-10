package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go-temperature/converter"
)

const (
	usageMessage = "Usage: go run main.go <temperature><C|F|K> [<C|F|K>]\n" +
		"Example: go run main.go 36.6C F, go run main.go 100F C, or go run main.go 300K C\n\n" +
		"Batch mode:\n  go run main.go --file temperatures.json --mode concurrent --workers 4 --buffer 4"

	errUnknownUnit       = "Unknown unit. Please use C, F, or K (e.g., 36.6C, 100F, or 300K)."
	errUnknownTargetUnit = "Unknown target unit. Please use C, F, or K."
)

func convertTemperature(args []string) (string, int) {
	if len(args) == 0 {
		return usageMessage, 1
	}

	input := strings.TrimSpace(args[0])

	if len(input) < 2 {
		return errUnknownUnit, 1
	}

	temperatureUnit, err := converter.ParseUnit(input[len(input)-1:])
	if err != nil {
		return errUnknownUnit, 1
	}

	temperatureValueInput := strings.TrimSpace(strings.ToUpper(input[:len(input)-1]))

	temperatureValue, err := strconv.ParseFloat(temperatureValueInput, 64)
	if err != nil {
		return fmt.Sprintf("Invalid number: %v", err), 1
	}

	targetUnit := converter.DefaultTargetUnit(temperatureUnit)
	if len(args) >= 2 {
		targetUnit, err = converter.ParseUnit(args[1])
		if err != nil {
			return errUnknownTargetUnit, 1
		}
	}

	convertedValue, err := converter.Convert(temperatureValue, temperatureUnit, targetUnit)
	if err != nil {
		return errUnknownTargetUnit, 1
	}

	return fmt.Sprintf("%.2f%v = %.2f%v", temperatureValue, converter.UnitSuffix(temperatureUnit), convertedValue, converter.UnitSuffix(targetUnit)), 0
}

func processBatch(filePath, mode string, workers, buffer int, timeout time.Duration) error {
	batchCtx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		batchCtx, cancel = context.WithTimeout(batchCtx, timeout)
		defer cancel()
	}

	jobs, err := converter.LoadJobsFromFile(batchCtx, filePath)
	if err != nil {
		return fmt.Errorf("load jobs: %w", err)
	}

	var (
		results []converter.Result
	)

	start := time.Now()
	switch strings.ToLower(mode) {
	case "sequential":
		results, err = converter.ProcessSequential(batchCtx, jobs)
	case "concurrent":
		results, err = converter.ProcessConcurrent(batchCtx, jobs, workers, buffer)
	default:
		return fmt.Errorf("unknown mode %q (expected sequential or concurrent)", mode)
	}

	if err != nil {
		return err
	}

	fmt.Printf("Processed %d job(s) in %s using %s mode\n", len(results), time.Since(start).Round(time.Millisecond), strings.ToLower(mode))

	for _, result := range results {
		source := converter.UnitSuffix(result.Job.From)
		target := converter.UnitSuffix(result.Job.To)
		if result.Err != nil {
			fmt.Printf("  [%d] %.2f%v -> %s ERROR: %v\n", result.Index, result.Job.Value, source, target, result.Err)
			continue
		}

		fmt.Printf("  [%d] %.2f%v = %.2f%v (worker time %s)\n",
			result.Index,
			result.Job.Value,
			source,
			result.Converted,
			target,
			result.Elapsed.Round(time.Microsecond),
		)
	}

	return nil
}

func main() {
	filePath := flag.String("file", "", "Path to a JSON file containing temperature conversion jobs.")
	mode := flag.String("mode", "sequential", "Processing mode: sequential or concurrent.")
	workers := flag.Int("workers", 0, "Number of worker goroutines for concurrent mode.")
	buffer := flag.Int("buffer", 0, "Buffer size for the results channel in concurrent mode.")
	timeout := flag.Duration("timeout", 0, "Optional timeout (e.g. 2s, 500ms) for batch processing.")
	flag.Parse()

	if *filePath != "" {
		if err := processBatch(*filePath, *mode, *workers, *buffer, *timeout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	result, exitCode := convertTemperature(flag.Args())
	if result != "" {
		fmt.Println(result)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
