package respond

import (
	"github.com/revel/revel"
)

// JSONResponse is gameshelf's standard json response
type JSONResponse struct {
	Status int `json:"status"`
}

// JSONErrors hold errors for a json response
type JSONErrors struct {
	Errors []string `json:"errors"`
}

// Add adds errors to a response
func (j *JSONErrors) Add(errors ...string) {
	j.Errors = append(j.Errors, errors...)
}

// AddFromValidation adds errors to a response from a revel validation
func (j *JSONErrors) AddFromValidation(errors []*revel.ValidationError) {
	for _, vError := range errors {
		j.Add(vError.Message)
	}
}

// ErrorResponse is gameshelf's standard json response when something goes wrong
type ErrorResponse struct {
	JSONResponse
	JSONErrors
}

// MessageResponse is gameshelf's standard json response when messages need to be conveyed
type MessageResponse struct {
	JSONResponse
	Messages []string `json:"messages"`
}

// SingleEntityResponse is gameshelf's standard single entity json response
type SingleEntityResponse struct {
	JSONResponse
	Entity interface{} `json:"entity"`
}

// MultipleEntityResponse is gameshelf's standard single entity json response
type MultipleEntityResponse struct {
	JSONResponse
	Entities []interface{} `json:"entities"`
}

// NewErrors creates a new JSONErrors
func NewErrors(initialErrors ...string) *JSONErrors {
	return &JSONErrors{
		Errors: initialErrors,
	}
}
