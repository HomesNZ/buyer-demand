package auth

import (
	"context"
	"errors"
	"net/http"
)

type Token struct {
	Valid bool

	User     *User
	Service  *Service
	Audience MultiString
}

func SetContextToken(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, contextKeyToken, token)
}

func ContextToken(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(contextKeyToken).(Token)
	return token, ok
}

// TokenFromContext returns the first valid token found in the the given
// context's data. This could be either a user or service token.
//
// DEPRECATED in favour of auth.ContextToken()
func TokenFromContext(ctx context.Context) (Token, error) {
	token, ok := ContextToken(ctx)
	if !ok {
		return token, errors.New("token not found")
	}
	return token, nil
}

// TokenFromHTTPRequest returns the first valid token found in the the given
// HTTP request's context data. This could be either a user or service token.
//
// DEPRECATED in favour of auth.ContextToken()
func TokenFromHTTPRequest(r *http.Request) (Token, error) {
	return TokenFromContext(r.Context())
}
