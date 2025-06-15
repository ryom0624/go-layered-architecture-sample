package entity

import "fmt"

// ValidationError represents a validation error for entity fields
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// BusinessLogicError represents a business logic error
type BusinessLogicError struct {
	Code    string
	Message string
}

func (e *BusinessLogicError) Error() string {
	return fmt.Sprintf("business logic error [%s]: %s", e.Code, e.Message)
}