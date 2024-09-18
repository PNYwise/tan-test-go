package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
)

type Validator interface {
	ValidateStruct(s interface{}) error
}

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &validatorImpl{
		validate: validator.New(),
	}
}

type ValidationError struct {
	Messages []string
}

func (v *ValidationError) Error() string {
	return strings.Join(v.Messages, ", ")
}

func (v *validatorImpl) ValidateStruct(s interface{}) error {
	// Check if the input is a slice of structs
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Slice {
		// Iterate over the slice and validate each struct
		var allErrors []string
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			if err := v.validate.Struct(item); err != nil {
				validationErrors := err.(validator.ValidationErrors)
				allErrors = append(allErrors, formatValidationErrors(validationErrors, i)...)
			}
		}
		if len(allErrors) > 0 {
			return &ValidationError{Messages: allErrors}
		}
		return nil
	}

	// If it's a single struct
	err := v.validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return newValidationError(validationErrors)
	}

	return nil
}

func newValidationError(validationErrors validator.ValidationErrors) *ValidationError {
	return &ValidationError{Messages: formatValidationErrors(validationErrors, -1)}
}

// formatValidationErrors converts validation errors into a slice of strings.
// If `index` is provided (i.e., for slices), it includes the index in the error message.
func formatValidationErrors(validationErrors validator.ValidationErrors, index int) []string {
	var errorMessages []string
	for _, err := range validationErrors {
		if index >= 0 {
			errorMessages = append(errorMessages, fmt.Sprintf("Element[%d] Field '%s' failed validation on rule '%s'", index, err.Field(), err.Tag()))
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation on rule '%s'", err.Field(), err.Tag()))
		}
	}
	return errorMessages
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateStruct(s interface{}) error {
	args := m.Called(s)
	return args.Error(0)
}
