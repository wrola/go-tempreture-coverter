package converter

import "fmt"

// Unit represents the supported temperature scales.
type Unit rune

const (
	// UnitCelsius is the Celsius temperature scale.
	UnitCelsius Unit = 'C'
	// UnitFahrenheit is the Fahrenheit temperature scale.
	UnitFahrenheit Unit = 'F'
	// UnitKelvin is the Kelvin temperature scale.
	UnitKelvin Unit = 'K'
)

// String converts a Unit to its string representation.
func (unit Unit) String() string {
	return string(unit)
}

// Validate ensures the unit is one of the supported scales.
func (unit Unit) Validate() error {
	switch unit {
	case UnitCelsius, UnitFahrenheit, UnitKelvin:
		return nil
	default:
		return fmt.Errorf("unknown unit %q", string(unit))
	}
}

// CelsiusToFahrenheit converts from Celsius to Fahrenheit.
func CelsiusToFahrenheit(celsiusValue float64) float64 {
	return (celsiusValue * 9 / 5) + 32
}

// FahrenheitToCelsius converts from Fahrenheit to Celsius.
func FahrenheitToCelsius(fahrenheitValue float64) float64 {
	return (fahrenheitValue - 32) * 5 / 9
}

// CelsiusToKelvin converts from Celsius to Kelvin.
func CelsiusToKelvin(celsiusValue float64) float64 {
	return celsiusValue + 273.15
}

// KelvinToCelsius converts from Kelvin to Celsius.
func KelvinToCelsius(kelvinValue float64) float64 {
	return kelvinValue - 273.15
}

// ToCelsius converts the value expressed in the given unit to Celsius.
func ToCelsius(value float64, unit Unit) (float64, error) {
	switch unit {
	case UnitCelsius:
		return value, nil
	case UnitFahrenheit:
		return FahrenheitToCelsius(value), nil
	case UnitKelvin:
		return KelvinToCelsius(value), nil
	default:
		return 0, fmt.Errorf("unknown unit %q", string(unit))
	}
}

// FromCelsius converts a Celsius value into the requested unit.
func FromCelsius(value float64, unit Unit) (float64, error) {
	switch unit {
	case UnitCelsius:
		return value, nil
	case UnitFahrenheit:
		return CelsiusToFahrenheit(value), nil
	case UnitKelvin:
		return CelsiusToKelvin(value), nil
	default:
		return 0, fmt.Errorf("unknown unit %q", string(unit))
	}
}

// Convert transforms a temperature value from one unit to another.
func Convert(value float64, from Unit, to Unit) (float64, error) {
	if err := from.Validate(); err != nil {
		return 0, err
	}
	if err := to.Validate(); err != nil {
		return 0, err
	}

	celsius, err := ToCelsius(value, from)
	if err != nil {
		return 0, err
	}

	return FromCelsius(celsius, to)
}

// DefaultTargetUnit suggests a sensible target unit based on the source.
func DefaultTargetUnit(source Unit) Unit {
	switch source {
	case UnitCelsius:
		return UnitFahrenheit
	case UnitFahrenheit:
		return UnitCelsius
	case UnitKelvin:
		return UnitCelsius
	default:
		return 0
	}
}

// UnitSuffix returns the human-friendly suffix for display purposes.
func UnitSuffix(unit Unit) string {
	switch unit {
	case UnitCelsius:
		return "°C"
	case UnitFahrenheit:
		return "°F"
	case UnitKelvin:
		return "K"
	default:
		return string(unit)
	}
}
