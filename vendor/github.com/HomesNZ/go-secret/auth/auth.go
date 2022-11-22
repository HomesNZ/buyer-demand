package auth

import (
	"encoding/base64"
	"net/http"

	"github.com/sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
)

type RequestAuthenticator interface {
	Authenticate(r *http.Request) (Token, error)
}

type Auth struct {
	logger                   logrus.FieldLogger
	authenticators           []RequestAuthenticator
	deprecatedServiceKeyAuth bool
}

type Option func(*Auth) error

func New(opts ...Option) (*Auth, error) {
	a := &Auth{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	if a.logger == nil {
		a.logger = logrus.StandardLogger()
	}
	return a, nil
}

func ServiceKey(key string) Option {
	return func(a *Auth) error {
		a.authenticators = append(a.authenticators, &ServiceAuthenticator{
			Key:    key,
			Legacy: a.deprecatedServiceKeyAuth,
		})
		return nil
	}
}

func hs256(base64Secret string) Option {
	return func(a *Auth) error {
		decoded, err := base64.RawURLEncoding.DecodeString(base64Secret)
		if err != nil {
			return err
		}
		authenticator := a.findOrCreateJWTAuth()
		authenticator.configs = append(authenticator.configs, hs256Parser{
			Parser: jwt.Parser{
				ValidMethods:         []string{"HS256"},
				SkipClaimsValidation: false,
			},
			Secret: decoded,
		})
		return nil
	}
}

func ClientSecret(base64Secret string) Option {
	return hs256(base64Secret)
}

func APISecret(base64Secret string) Option {
	return hs256(base64Secret)
}

func JWKS(url string) Option {
	return func(a *Auth) error {
		authenticator := a.findOrCreateJWTAuth()
		authenticator.configs = append(authenticator.configs, &jwksParser{
			Parser: jwt.Parser{
				ValidMethods:         []string{"RS256"},
				SkipClaimsValidation: false,
			},
			URL: url,
		})
		return nil
	}
}

func DeprecatedServiceKeyAuth(enabled bool) Option {
	return func(a *Auth) error {
		a.deprecatedServiceKeyAuth = enabled // so this can be called before the authenticator is created
		for _, a := range a.authenticators {
			if s, ok := a.(*ServiceAuthenticator); ok {
				s.Legacy = enabled
			}
		}
		return nil
	}
}

func Logger(logger logrus.FieldLogger) Option {
	return func(a *Auth) error {
		a.logger = logger
		return nil
	}
}

func (a *Auth) Authenticate(r *http.Request) (Token, error) {
	for _, authenticator := range a.authenticators {
		t, err := authenticator.Authenticate(r)
		if err == nil && !t.Valid {
			continue
		}
		return t, err
	}
	return Token{}, nil
}

func (a *Auth) CheckRequest(r *http.Request) (*http.Request, bool) {
	t, err := a.Authenticate(r)
	if err != nil {
		return r, false
	}
	return r.WithContext(SetContextToken(r.Context(), t)), t.Valid
}

// findOrCreateJWTAuth is a helper to manage a single instance JWTAuthenticator
func (a *Auth) findOrCreateJWTAuth() *JWTAuthenticator {
	for _, authenticator := range a.authenticators {
		if authenticator, ok := authenticator.(*JWTAuthenticator); ok {
			return authenticator
		}
	}
	authenticator := &JWTAuthenticator{}
	a.authenticators = append(a.authenticators, authenticator)
	return authenticator
}
