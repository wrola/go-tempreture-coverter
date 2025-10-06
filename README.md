# Go Temperature Converter

A tiny command-line utility for converting temperatures between Celsius, Fahrenheit, and Kelvin. It accepts a value with its unit (for example `36.6C`, `100F`, or `300K`) and optionally a target unit, performs the conversion, and prints the result with two decimal places. The project is intentionally small and self-contained so it can serve as a reference for basic Go CLI patterns and table-driven tests.

## Requirements

- Go 1.24 or newer (see `go.mod`)

## Getting Started

Clone the repository and navigate into the project directory:

```bash
git clone <repository-url>
cd go-tempreture-coverter
```

You can run the tool without building a binary:

```bash
go run . 36.6C
```

Or build a reusable executable:

```bash
go build -o temp-converter
./temp-converter 100F
```

## Usage

Provide a value (with `C`, `F`, or `K`) and optionally a target scale (`C`, `F`, or `K`). If you skip the target unit, the tool converts to a sensible default (for example Celsius to Fahrenheit).

```bash
go run . 36.6C
# 36.60°C = 97.88°F

go run . 100F
# 100.00°F = 37.78°C

go run . 300K
# 300.00K = 26.85°C

go run . 36.6C K
# 36.60°C = 309.75K
```

Invalid input is reported with an error message and a non-zero exit code. For example:

```bash
go run . 10C X
# Unknown target unit. Please use C, F, or K.
```

## Testing

The project includes table-driven unit tests that cover successful conversions and common error paths. Run them with:

```bash
go test ./...
```

## Project Structure

- `main.go` – command-line entry point and conversion helpers.
- `main_test.go` – table-driven tests for the conversion logic.

## License

This project is licensed under the terms of the [MIT License](LICENSE).
