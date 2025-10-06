package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CelviusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func convertTemperature(args []string) (string, int) {
	if len(args) < 1 {
		message := strings.Join([]string{
			"Usage: go run main.go <temperature><C|F>",
			"Example: go run main.go 36.6C or go run main.go 100F",
		}, "\n")
		return message, 1
	}

	input := args[0]
	sanitizeInput := strings.TrimSpace(strings.ToUpper(input))

	if len(sanitizeInput) < 2 {
		return "Unknown unit. Please use C or F (e.g., 36.6C or 100F).", 1
	}

	temperatureUnit := sanitizeInput[len(sanitizeInput)-1]
	temperatureValueInput := sanitizeInput[:len(sanitizeInput)-1]

	temperatureValue, err := strconv.ParseFloat(temperatureValueInput, 64)
	if err != nil {
		return fmt.Sprintf("Invalid number: %v", err), 1
	}

	switch temperatureUnit {
	case 'C':
		return fmt.Sprintf("%.2f°C = %.2f°F", temperatureValue, CelviusToFahrenheit(temperatureValue)), 0
	case 'F':
		return fmt.Sprintf("%.2f°F = %.2f°C", temperatureValue, FahrenheitToCelsius(temperatureValue)), 0
	default:
		return "Unknown unit. Please use C or F (e.g., 36.6C or 100F).", 1
	}
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
