package errors

import (
	"errors"
	"strings"
)

var ErrAdsList = errors.New("ads list error")
var ErrAd = errors.New("ad error")
var ErrCreateAd = errors.New("ad create error")
var ErrUpdateAd = errors.New("ad update error")
var ErrDeleteAd = errors.New("ad delete error")
var ErrSoldAd = errors.New("ad sold error")
var ErrCreateAdValidation = errors.New("ad create error validation")
var ErrUpdateAdValidation = errors.New("ad update error validation")
var ErrForbidden = errors.New("forbidden")
var ErrNotValidUuid = errors.New("not valid uuid")

var ErrAdNotFound = errors.New("ad not found")
var ErrAdUserNotFound = errors.New("ad user not found")

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
