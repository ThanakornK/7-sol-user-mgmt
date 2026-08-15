// package utils provides utility functions.
package utils

import (
	"errors"
	"net/http"
)

// ResponseStruct is the response struct used in the API.
type ResponseStruct struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

// NewResponseStruct creates a new response struct. can customize the status, msg, data.
func NewResponseStruct(status bool, msg string, data interface{}) *ResponseStruct {
	return &ResponseStruct{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
}

// NewSuccessResponseStruct creates a new success response struct. status is true and assign only msg and data.
func NewSuccessResponseStruct(msg string, data interface{}) *ResponseStruct {
	return NewResponseStruct(true, msg, data)
}

// NewErrorResponseStruct creates a new error response struct. status is false and assign only msg and error.
func NewErrorResponseStruct(msg string, error interface{}) *ResponseStruct {
	return &ResponseStruct{
		Status: false,
		Msg:    msg,
		Error:  error,
	}
}

// MapErrorResponse maps known errors to a public HTTP response and hides unknown errors.
func MapErrorResponse(err error, status int, msg string, publicError interface{}, knownErrors ...error) (int, *ResponseStruct) {
	for _, knownError := range knownErrors {
		if errors.Is(err, knownError) {
			return status, NewErrorResponseStruct(msg, publicError)
		}
	}
	return http.StatusInternalServerError, NewErrorResponseStruct(msg, "internal server error")
}

// ErrorMessages returns the error messages from the slice of errors.
func ErrorMessages(errs []error) []string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}
	return messages
}
