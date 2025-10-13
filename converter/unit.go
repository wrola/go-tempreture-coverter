package converter

import (
	"encoding/json"
	"strings"
)

func ParseUnit(unitValue string) (Unit, error) {
	normalized := strings.TrimSpace(strings.ToUpper(unitValue))
	if normalized == "" {
		return 0, &ValidationError{
			Operation: "parse unit",
			Field:  "unit",
			Value:  unitValue,
			Reason: "cannot be empty",
		}
	}

	if len(normalized) != 1 {
		return 0, &ValidationError{
			Operation: "parse unit",
			Field:  "unit",
			Value:  unitValue,
			Reason: "must be exactly one character",
		}
	}

	unit := Unit(normalized[0])
	if err := unit.Validate(); err != nil {
		return 0, err
	}

	return unit, nil
}

func MustParseUnit(unitValue string) Unit {
	unit, err := ParseUnit(unitValue)
	if err != nil {
		panic(err)
	}
	return unit
}

func (unit *Unit) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	parsedUnit, err := ParseUnit(raw)
	if err != nil {
		return err
	}

	*unit = parsedUnit
	return nil
}

func (unit Unit) MarshalJSON() ([]byte, error) {
	if err := unit.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(unit.String())
}
