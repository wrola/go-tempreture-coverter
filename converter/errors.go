package converter

import "fmt"

type ConversionError struct {
	Err       error
	Operation string
	Value     float64
	Unit      Unit
}

func (e *ConversionError) Error() string {
	if e.Unit != 0 {
		return fmt.Sprintf("conversion error in %s (unit=%s): %v", e.Operation, e.Unit, e.Err)
	}
	return fmt.Sprintf("conversion error in %s: %v", e.Operation, e.Err)
}

func (e *ConversionError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Operation string
	Field     string
	Value     string
	Reason    string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error in %s: field %s with value %q: %s", e.Operation, e.Field, e.Value, e.Reason)
	}
	return fmt.Sprintf("validation error in %s: %s", e.Operation, e.Reason)
}

type ProcessingError struct {
	Err       error
	Operation string
	JobIndex  int
}

func (e *ProcessingError) Error() string {
	if e.JobIndex >= 0 {
		return fmt.Sprintf("processing error in %s (job %d): %v", e.Operation, e.JobIndex, e.Err)
	}
	return fmt.Sprintf("processing error in %s: %v", e.Operation, e.Err)
}

func (e *ProcessingError) Unwrap() error {
	return e.Err
}
