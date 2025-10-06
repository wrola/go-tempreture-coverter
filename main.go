package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CelsiusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func CelsiusToKelvin(c float64) float64 {
	return c + 273.15
}

func KelvinToCelsius(k float64) float64 {
	return k - 273.15
}

func toCelsius(value float64, unit rune) (float64, error) {
	switch unit {
	case 'C':
		return value, nil
	case 'F':
		return FahrenheitToCelsius(value), nil
	case 'K':
		return KelvinToCelsius(value), nil
	default:
		return 0, fmt.Errorf("unknown unit %q", string(unit))
	}
}

func fromCelsius(value float64, unit rune) (float64, error) {
	switch unit {
	case 'C':
		return value, nil
	case 'F':
		return CelsiusToFahrenheit(value), nil
	case 'K':
		return CelsiusToKelvin(value), nil
	default:
		return 0, fmt.Errorf("unknown unit %q", string(unit))
	}
}

func unitSuffix(unit rune) string {
	switch unit {
	case 'C':
		return "°C"
	case 'F':
		return "°F"
	case 'K':
		return "K"
	default:
		return string(unit)
	}
}

func defaultTargetUnit(source rune) rune {
	switch source {
	case 'C':
		return 'F'
	case 'F':
		return 'C'
	case 'K':
		return 'C'
	default:
		return 0
	}
}

func convertTemperature(args []string) (string, int) {
	if len(args) < 1 {
		message := strings.Join([]string{
			"Usage: go run main.go <temperature><C|F|K> [<C|F|K>]",
			"Example: go run main.go 36.6C F, go run main.go 100F C, or go run main.go 300K C",
		}, "\n")
		return message, 1
	}

	input := args[0]
	sanitizeInput := strings.TrimSpace(strings.ToUpper(input))

	if len(sanitizeInput) < 2 {
		return "Unknown unit. Please use C, F, or K (e.g., 36.6C, 100F, or 300K).", 1
	}

	temperatureUnit := rune(sanitizeInput[len(sanitizeInput)-1])
	temperatureValueInput := sanitizeInput[:len(sanitizeInput)-1]

	temperatureValue, err := strconv.ParseFloat(temperatureValueInput, 64)
	if err != nil {
		return fmt.Sprintf("Invalid number: %v", err), 1
	}

	if temperatureUnit != 'C' && temperatureUnit != 'F' && temperatureUnit != 'K' {
		return "Unknown unit. Please use C, F, or K (e.g., 36.6C, 100F, or 300K).", 1
	}

	targetUnit := defaultTargetUnit(temperatureUnit)
	if len(args) >= 2 {
		targetInput := strings.TrimSpace(strings.ToUpper(args[1]))
		if len(targetInput) != 1 {
			return "Unknown target unit. Please use C, F, or K.", 1
		}
		targetUnit = rune(targetInput[0])
	}

	if targetUnit != 'C' && targetUnit != 'F' && targetUnit != 'K' {
		return "Unknown target unit. Please use C, F, or K.", 1
	}

	valueInCelsius, err := toCelsius(temperatureValue, temperatureUnit)
	if err != nil {
		return "Unknown unit. Please use C, F, or K (e.g., 36.6C, 100F, or 300K).", 1
	}

	convertedValue, err := fromCelsius(valueInCelsius, targetUnit)
	if err != nil {
		return "Unknown target unit. Please use C, F, or K.", 1
	}

	return fmt.Sprintf("%.2f%v = %.2f%v", temperatureValue, unitSuffix(temperatureUnit), convertedValue, unitSuffix(targetUnit)), 0
}

func main() {
	result, exitCode := convertTemperature(os.Args[1:])
	if result != "" {
		fmt.Println(result)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
