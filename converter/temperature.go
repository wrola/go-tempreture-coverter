package converter

type Unit rune

const (
	UnitCelsius Unit = 'C'
	UnitFahrenheit Unit = 'F'
	UnitKelvin Unit = 'K'
)

func (unit Unit) String() string {
	return string(unit)
}

func (unit Unit) Validate() error {
	switch unit {
	case UnitCelsius, UnitFahrenheit, UnitKelvin:
		return nil
	default:
		return &ValidationError{
			Operation: "unit validation",
			Field:  "unit",
			Value:  string(unit),
			Reason: "must be one of C, F, or K",
		}
	}
}

func CelsiusToFahrenheit(celsiusValue float64) float64 {
	return (celsiusValue * 9 / 5) + 32
}

func FahrenheitToCelsius(fahrenheitValue float64) float64 {
	return (fahrenheitValue - 32) * 5 / 9
}

func CelsiusToKelvin(celsiusValue float64) float64 {
	return celsiusValue + 273.15
}

func KelvinToCelsius(kelvinValue float64) float64 {
	return kelvinValue - 273.15
}

func ToCelsius(value float64, unit Unit) (float64, error) {
	switch unit {
	case UnitCelsius:
		return value, nil
	case UnitFahrenheit:
		return FahrenheitToCelsius(value), nil
	case UnitKelvin:
		return KelvinToCelsius(value), nil
	default:
		return 0, &ConversionError{
			Operation: "convert to celsius",
			Unit: unit,
			Err:  &ValidationError{
				Operation: "unit validation",
				Field:  "source unit",
				Value:  string(unit),
				Reason: "unknown unit",
			},
		}
	}
}

func FromCelsius(value float64, unit Unit) (float64, error) {
	switch unit {
	case UnitCelsius:
		return value, nil
	case UnitFahrenheit:
		return CelsiusToFahrenheit(value), nil
	case UnitKelvin:
		return CelsiusToKelvin(value), nil
	default:
		return 0, &ConversionError{
			Operation: "convert from celsius",
			Unit: unit,
			Err:  &ValidationError{
				Operation: "unit validation",
				Field:  "target unit",
				Value:  string(unit),
				Reason: "unknown unit",
			},
		}
	}
}

func Convert(value float64, from Unit, to Unit) (float64, error) {
	if err := from.Validate(); err != nil {
		return 0, &ConversionError{
			Operation: "validate source unit",
			Unit:  from,
			Value: value,
			Err:   err,
		}
	}
	if err := to.Validate(); err != nil {
		return 0, &ConversionError{
			Operation: "validate target unit",
			Unit: to,
			Err:  err,
		}
	}

	celsius, err := ToCelsius(value, from)
	if err != nil {
		return 0, err
	}

	return FromCelsius(celsius, to)
}

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
