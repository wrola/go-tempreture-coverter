package converter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseUnit normalizes the candidate string and returns the corresponding Unit.
func ParseUnit(unitValue string) (Unit, error) {
	normalized := strings.TrimSpace(strings.ToUpper(unitValue))
	if normalized == "" {
		return 0, fmt.Errorf("unknown unit %q", unitValue)
	}

	if len(normalized) != 1 {
		return 0, fmt.Errorf("unknown unit %q", unitValue)
	}

	unit := Unit(normalized[0])
	if err := unit.Validate(); err != nil {
		return 0, err
	}

	return unit, nil
}

// MustParseUnit panics if the unit string cannot be parsed.
func MustParseUnit(unitValue string) Unit {
	unit, err := ParseUnit(unitValue)
	if err != nil {
		panic(err)
	}
	return unit
}

// UnmarshalJSON converts a JSON string into a Unit.
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

// MarshalJSON serializes the Unit as a single-character string.
func (unit Unit) MarshalJSON() ([]byte, error) {
	if err := unit.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(unit.String())
}
