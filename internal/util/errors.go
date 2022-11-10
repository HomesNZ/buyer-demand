package util

import (
	"net/http"

	"github.com/pkg/errors"
)

var New = errors.New
var Wrap = errors.Wrap
var Wrapf = errors.Wrapf

type statusError struct {
	Message string
	Status  int
}

func (e statusError) Error() string {
	return e.Message
}

func (e statusError) StatusCode() int {
	return e.Status
}
func (e statusError) NotFound() bool {
	return e.Status == http.StatusNotFound
}

func NotFound() error {
	return statusError{
		Message: "Not found",
		Status:  http.StatusNotFound,
	}
}

func Unauthorized(message string) error {
	return statusError{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func Forbidden() error {
	return statusError{
		Message: "Forbidden",
		Status:  http.StatusForbidden,
	}
}

func BadRequest(message string) error {
	return statusError{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}
