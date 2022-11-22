package auth

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrNoServiceFound = errors.New("Unable to retrieve service from token")
)

type Service struct {
	Name string
}

// DEPRECATED in favour of auth.ContextToken()
func ServiceFromContext(c context.Context) (*Service, error) {
	token, err := TokenFromContext(c)
	if err != nil {
		return nil, err
	}

	if token.Service == nil {
		return nil, ErrNoServiceFound
	}

	return token.Service, nil
}

// DEPRECATED in favour of auth.ContextToken()
func ServiceFromHTTPRequest(r *http.Request) (*Service, error) {
	return ServiceFromContext(r.Context())
}
