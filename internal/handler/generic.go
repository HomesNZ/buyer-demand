package handler

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type genericResponse struct {
	Error string `json:"error"`
}

type statusCoder interface {
	error
	StatusCode() int
}

func EncodeErrorResponse(logger *logrus.Entry, w http.ResponseWriter, err error) {
	var res genericResponse
	code := http.StatusInternalServerError
	switch er := errors.Cause(err).(type) {
	case *json.SyntaxError, validation.Errors, schema.MultiError:
		code = http.StatusBadRequest
		res.Error = er.Error()
	case statusCoder:
		code = er.StatusCode()
		res.Error = er.Error()
	default:
		logger.WithError(err).Error()
		res.Error = "Something went wrong"
	}
	if code == http.StatusNoContent {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		logger.WithError(err).Error()
	}
}

func EncodeJSONResponse(logger *logrus.Entry, w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		logger.WithError(err).Error()
	}
}
