package errors

import "strings"

type ValidationError struct {
	Errors []ValidationErrorItem `json:"errors"`
}

func NewValidationError() *ValidationError {
	return &ValidationError{}
}

type ValidationErrorItem struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func (v *ValidationError) Add(field, error string) {
	v.Errors = append(v.Errors, ValidationErrorItem{Field: field, Error: error})
}

func (v *ValidationError) HasErrors() bool {
    return len(v.Errors) > 0
}

func (v *ValidationError) Error() string {
	var result strings.Builder

	for _, item := range v.Errors {
		result.WriteString(";")
		result.WriteString(item.Error)
	}

    return result.String()
}
