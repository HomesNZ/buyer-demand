package auth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

const (
	authorizationSchemeService = "HomesService"
)

// AuthenticateUserRequest forwards the Authorization header from an existing
// request object (src) to a new request object (dest). It returns a boolean
// which indicates whether the header was successfully copied.
//
// The header will only be copied if:
// - The JWT token is valid and was already parsed and added to the context of
//   the request (any UserMiddleware in this package will do this for you).
// - The token is a valid JWT token. Service tokens will not be forwarded.
func AuthenticateUserRequest(src, dest *http.Request) bool {
	token, err := TokenFromContext(src.Context())
	if err != nil {
		return false
	}

	// If the resulting token is not a 'User' token, then we don't forward it on.
	// Service tokens should explicitly not be forwarded.
	if token.User == nil || token.Service != nil {
		return false
	}

	// Copy the Authorization header from src to dest, return true to indicate
	// that src was indeed a valid user request.
	dest.Header.Set(
		"Authorization",
		src.Header.Get("Authorization"),
	)
	return true
}

// AuthenticateServiceRequest sets the Authorization of the given HTTP request
// (dest) to a base64-encoded string containing the service name and
// authorisation key, separated by a colon (:). This can be parsed at the other
// end using the middlewares contained in this package.
//
// The service name and key should not contain any colon characters.
//
// TODO: Ideally we would use JWT for service-to-service authentication as well,
// but Auth0 doesn't seem to have any solid pattern for doing this. This
// 'service key' implementation should suffice for our requirements for the time
// being.
func AuthenticateServiceRequest(dest *http.Request, service, key string) {
	encoded := base64.StdEncoding.EncodeToString([]byte(service + ":" + key))
	dest.Header.Set("Authorization", authorizationSchemeService+" "+encoded)
}

type ServiceAuthenticator struct {
	Key    string
	Legacy bool
}

func (a *ServiceAuthenticator) Authenticate(r *http.Request) (Token, error) {
	if a.Legacy {
		token, err := a.getServiceFromKeyParam(r)
		if err == nil && token != nil {
			return *token, nil
		}
	}
	return a.getServiceFromHeaders(r)
}

func (a *ServiceAuthenticator) getServiceFromHeaders(r *http.Request) (Token, error) {
	authorization := r.Header.Get("Authorization")
	prefix := authorizationSchemeService + " "
	if !strings.HasPrefix(authorization, prefix) {
		return Token{}, nil
	}

	decoded, _ := base64.StdEncoding.DecodeString(authorization[len(prefix):])
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return Token{}, errAuthInvalid
	}

	if parts[1] != a.Key {
		return Token{}, errAuthInvalid
	}

	s := Service{
		Name: parts[0],
	}
	if s.Name == "" {
		return Token{}, errAuthInvalid
	}

	return Token{
		Valid:   true,
		Service: &s,
	}, nil
}

// getServiceFromKeyParam is a deprecated method of authenticating a service
// request. It retrieves the token from the 'key' URL parameter, and returns
// 'deprecated' as the service name.
func (a *ServiceAuthenticator) getServiceFromKeyParam(r *http.Request) (*Token, error) {
	key := r.URL.Query().Get("key")
	if key == "" {
		return nil, nil
	}
	if key != a.Key {
		return &Token{}, errAuthInvalid
	}
	return &Token{
		Valid: true,
		Service: &Service{
			Name: "deprecated",
		},
	}, nil
}
