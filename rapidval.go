// Package rapidval provides a high-performance, zero-dependency validation library for Go.
// It focuses on simplicity, extensibility, and performance while providing a clean and intuitive API.
//
// Basic usage:
//
//	type User struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	func (u *User) Validate(v *rapidval.Validator) error {
//	    return v.Validate(rapidval.P{
//	        rapidval.Required("Name", u.Name),
//	        rapidval.MinLength("Name", u.Name, 2),
//	        rapidval.Email("Email", u.Email),
//	        rapidval.Between("Age", u.Age, 18, 100),
//	    })
//	}
package rapidval

import (
	"strings"
	"time"
)

// Validateable is an interface that can be implemented by any struct to add custom validation logic.
// It allows for custom validation logic to be added to a struct without having to implement the Validate method.
type Validateable interface {

	// Validations is a method that can be implemented by any struct to add custom validation logic.
	// It allows for custom validation logic to be added to a struct without having to implement the Validate method.
	Validations() P
}

// ValidationError represents a single validation error.
// It contains the field name, message key, and any parameters needed for translation.
type ValidationError struct {
	Field         string
	MessageKey    string
	MessageParams map[string]interface{}
	CurrentValue  interface{}
}

// Error implements the error interface.
// It returns the message key by default, which can be translated using a Translator.
func (ve *ValidationError) Error() string {
	return ve.MessageKey
}

// ValidationErrors represents a collection of validation errors.
type ValidationErrors []*ValidationError

// Error implements the error interface.
// It joins all validation error messages with semicolons.
func (ve ValidationErrors) Error() string {
	var errors []string
	for _, err := range ve {
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, "; ")
}

// Validator handles the validation process and collects validation errors.
type Validator struct {
	errors ValidationErrors
}

// P (Params) is a collection of validation errors used for grouping validations.
type P []*ValidationError

// Validate processes all validation rules and returns any validation errors.
// If there are no errors, it returns nil.
func (v *Validator) Validate(val Validateable) error {
	params := val.Validations()
	if len(params) == 0 {
		return nil
	}

	for _, err := range params {
		if err != nil && err.MessageKey != "" {
			v.errors = append(v.errors, err)
		}
	}

	if len(v.errors) > 0 {
		return v.errors
	}

	return nil
}

// Message Keys
const (
	MsgRequired        = "validation.required"
	MsgInvalidEmail    = "validation.email"
	MsgMinLength       = "validation.min_length"
	MsgMaxLength       = "validation.max_length"
	MsgBetween         = "validation.between"
	MsgDateGreaterThan = "validation.date_greater_than"
	MsgDateLessThan    = "validation.date_less_than"
)

// MessageParam keys
const (
	Field = "Field"
	Min   = "Min"
	Max   = "Max"
	Value = "Value"
)

// Required checks if a value is not zero according to its type.
// For strings, it checks if the string is not empty.
// For numbers, it checks if the number is not zero.
// For time.Time, it checks if the time is not zero.
// For pointers and interfaces, it checks if the value is not nil.
func Required(field string, value interface{}) *ValidationError {
	if isZero(value) {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgRequired,
			MessageParams: map[string]interface{}{
				Field: field,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// Email validates if a string is a valid email address.
// Currently checks for @ and . characters.
func Email(field string, value string) *ValidationError {
	if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgInvalidEmail,
			MessageParams: map[string]interface{}{
				Field: field,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// MinLength validates if a string's length is at least the specified minimum.
func MinLength(field string, value string, min int) *ValidationError {
	if len(value) < min {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgMinLength,
			MessageParams: map[string]interface{}{
				Field: field,
				Min:   min,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// MaxLength validates if a string's length is at most the specified maximum.
func MaxLength(field string, value string, max int) *ValidationError {
	if len(value) > max {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgMaxLength,
			MessageParams: map[string]interface{}{
				Field: field,
				Max:   max,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// Between validates if a number is between the specified minimum and maximum values (inclusive).
func Between(field string, value int, min, max int) *ValidationError {
	if value < min || value > max {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgBetween,
			MessageParams: map[string]interface{}{
				Field: field,
				Min:   min,
				Max:   max,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// DateGreaterThan validates if a time.Time is after the specified minimum time.
func DateGreaterThan(field string, value, min time.Time) *ValidationError {
	if value.Before(min) {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgDateGreaterThan,
			MessageParams: map[string]interface{}{
				Field: field,
				Min:   min,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// DateLessThan validates if a time.Time is before the specified maximum time.
func DateLessThan(field string, value, max time.Time) *ValidationError {
	if value.After(max) {
		return &ValidationError{
			Field:      field,
			MessageKey: MsgDateLessThan,
			MessageParams: map[string]interface{}{
				Field: field,
				Max:   max,
				Value: value,
			},
			CurrentValue: value,
		}
	}
	return nil
}

// isZero checks if a value is the zero value for its type.
// This is used internally by the Required validation.
func isZero(v interface{}) bool {
	switch v := v.(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	case bool:
		return !v
	case time.Time:
		return v.IsZero()
	case nil:
		return true
	}
	return false
}

// New returns a new Validator.
func New() *Validator {
	return &Validator{}
}
