package qp

import (
	"github.com/pkg/errors"
)

// ValidationFunc represents a validator for Filters.
type ValidationFunc func(value interface{}) error

// Validations type replacement for map.
// Used in NewParse(), NewQV(), SetValidations().
type Validations map[string]ValidationFunc

// Multi combines multiple validation functions.
// Usage: Multi(Min(10), Max(100)).
func Multi(validators ...ValidationFunc) ValidationFunc {
	return func(value interface{}) error {
		for _, validator := range validators {
			if err := validator(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// In checks if the value is within the provided values.
func In(validValues ...interface{}) ValidationFunc {
	return func(value interface{}) error {
		for _, validValue := range validValues {
			if validValue == value {
				return nil
			}
		}
		return errors.Wrapf(ErrNotInScope, "%v", value)
	}
}

// Min checks if the value is greater than or equal to the minimum.
func Min(min int) ValidationFunc {
	return func(value interface{}) error {
		if intValue, ok := value.(int); ok {
			if intValue >= min {
				return nil
			}
		}
		return errors.Wrapf(ErrNotInScope, "%v", value)
	}
}

// Max checks if the value is less than or equal to the maximum.
func Max(max int) ValidationFunc {
	return func(value interface{}) error {
		if intValue, ok := value.(int); ok {
			if intValue <= max {
				return nil
			}
		}
		return errors.Wrapf(ErrNotInScope, "%v", value)
	}
}

// MinMax checks if the value is between or equal to the minimum and maximum.
func MinMax(min, max int) ValidationFunc {
	return func(value interface{}) error {
		if intValue, ok := value.(int); ok {
			if min <= intValue && intValue <= max {
				return nil
			}
		}
		return errors.Wrapf(ErrNotInScope, "%v", value)
	}
}

// NotEmpty checks if the string value is not empty.
func NotEmpty() ValidationFunc {
	return func(value interface{}) error {
		if str, ok := value.(string); ok {
			if len(str) > 0 {
				return nil
			}
		}
		return errors.Wrapf(ErrNotInScope, "%v", value)
	}
}
