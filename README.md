# Go Temperature Converter

A tiny command-line utility for converting temperatures between Celsius and Fahrenheit. It accepts a value with its unit (for example `36.6C` or `100F`), performs the conversion, and prints the result with two decimal places. The project is intentionally small and self-contained so it can serve as a reference for basic Go CLI patterns and table-driven tests.

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

Provide a single argument combining the numeric value with either `C` or `F`.

```bash
go run . 36.6C
# 36.60°C = 97.88°F

go run . 100F
# 100.00°F = 37.78°C
```

Invalid input is reported with an error message and a non-zero exit code. For example:

```bash
go run . 10K
# Unknown unit. Please use C or F (e.g., 36.6C or 100F).
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
