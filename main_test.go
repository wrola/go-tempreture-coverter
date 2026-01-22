package main

import (
	"strings"
	"testing"
)

func TestConvertTemperature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantContains []string
		args         []string
		want         string
		name         string
		wantExitCode int
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
			name:         "KelvinToCelsius",
			args:         []string{"300k"},
			wantExitCode: 0,
			want:         "300.00K = 26.85°C",
		},
		{
			name:         "CelsiusToKelvinTarget",
			args:         []string{"25c", "K"},
			wantExitCode: 0,
			want:         "25.00°C = 298.15K",
		},
		{
			name:         "FahrenheitToKelvinTarget",
			args:         []string{"32F", "k"},
			wantExitCode: 0,
			want:         "32.00°F = 273.15K",
		},
		{
			name:         "NoArguments",
			args:         nil,
			wantExitCode: 1,
			wantContains: []string{
				"Usage: go run main.go <temperature><C|F|K> [<C|F|K>]",
				"Example: go run main.go 36.6C F, go run main.go 100F C, or go run main.go 300K C",
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
			args:         []string{"10X"},
			wantExitCode: 1,
			want:         "Unknown unit. Please use C, F, or K (e.g., 36.6C, 100F, or 300K).",
		},
		{
			name:         "UnknownTargetUnit",
			args:         []string{"10C", "X"},
			wantExitCode: 1,
			want:         "Unknown target unit. Please use C, F, or K.",
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
