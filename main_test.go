package main

import (
	"strings"
	"testing"
)

func TestConvertTemperature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		want         string
		wantContains []string
	}{
		{
			name:         "CelsiusToFahrenheit",
			args:         []string{"36.6C"},
			wantExitCode: 0,
			want:         "36.60°C = 97.88°F",
		},
		{
			name:         "FahrenheitToCelsius",
			args:         []string{"100f"},
			wantExitCode: 0,
			want:         "100.00°F = 37.78°C",
		},
		{
			name:         "NoArguments",
			args:         nil,
			wantExitCode: 1,
			wantContains: []string{
				"Usage: go run main.go <temperature><C|F>",
				"Example: go run main.go 36.6C or go run main.go 100F",
			},
		},
		{
			name:         "InvalidNumber",
			args:         []string{"abcC"},
			wantExitCode: 1,
			want:         "Invalid number: strconv.ParseFloat: parsing \"ABC\": invalid syntax",
		},
		{
			name:         "UnknownUnit",
			args:         []string{"10K"},
			wantExitCode: 1,
			want:         "Unknown unit. Please use C or F (e.g., 36.6C or 100F).",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, exitCode := convertTemperature(tc.args)

			if exitCode != tc.wantExitCode {
				t.Fatalf("exit code = %d, want %d; output: %s", exitCode, tc.wantExitCode, got)
			}

			if tc.want != "" && got != tc.want {
				t.Fatalf("output = %q, want %q", got, tc.want)
			}

			for _, want := range tc.wantContains {
				if !strings.Contains(got, want) {
					t.Fatalf("output %q does not contain %q", got, want)
				}
			}
		})
	}
}
