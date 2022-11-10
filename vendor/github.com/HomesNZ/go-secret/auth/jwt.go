package auth

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

const jwkExpiry = 24 * time.Hour

type parser interface {
	Parse(ctx context.Context, token string) (Token, error)
}

type JWTAuthenticator struct {
	configs []parser
}

func (a *JWTAuthenticator) Authenticate(r *http.Request) (Token, error) {
	token, err := a.extractJWTFromAuthorizationHeader(r)
	if err != nil || token == "" {
		return Token{}, nil
	}
	for _, config := range a.configs {
		t, err := config.Parse(r.Context(), token)
		if err == nil {
			return t, nil
		}
	}
	return Token{}, errAuthInvalid
}

// extractJWTFromAuthorizationHeader extracts the JWT token from the
// Authorization header.
// Copied from https://github.com/auth0/go-jwt-middleware/blob/master/jwtmiddleware.go
func (a *JWTAuthenticator) extractJWTFromAuthorizationHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}
	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if strings.EqualFold(authHeaderParts[0], "oauth") {
		// If request contains an OAuth token, quietly ignore it. This is probably
		// Quinovic, and probably a request to the Listing Service (TM API)...
		return "", nil
	}
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "bearer") {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}
	return authHeaderParts[1], nil
}

type hs256Parser struct {
	Secret []byte
	Parser jwt.Parser
}

func (p hs256Parser) key(token *jwt.Token) (interface{}, error) {
	return p.Secret, nil
}

func (p hs256Parser) Parse(ctx context.Context, token string) (Token, error) {
	claims := &rawClaims{}
	_, err := jwt.ParseWithClaims(token, claims, p.key)
	if err != nil {
		return Token{}, err
	}
	if time.Unix(int64(claims.Exp), 0).Before(time.Now()) {
		return Token{}, errors.New("Token expired")
	}
	return claims.Token(), nil
}

type jwksParser struct {
	URL    string
	Parser jwt.Parser

	mu      sync.RWMutex
	expires time.Time
	cache   *jwk.Set
}

func (p *jwksParser) key(token *jwt.Token) (interface{}, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing header kid")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cache == nil || time.Now().After(p.expires) {
		// Refresh cache
		var err error
		p.cache, err = jwk.Fetch(p.URL)
		if err != nil {
			return nil, err
		}
		p.expires = time.Now().Add(jwkExpiry)
	}
	keys := p.cache.LookupKeyID(kid)
	if len(keys) < 1 {
		return nil, errors.Errorf("JWKS has no kid: %s", kid)
	}
	return keys[0].Materialize()
}

func (p *jwksParser) Parse(ctx context.Context, token string) (Token, error) {
	claims := &rawClaims{}
	_, err := jwt.ParseWithClaims(token, claims, p.key)
	if err != nil {
		return Token{}, err
	}
	if time.Unix(int64(claims.Exp), 0).Before(time.Now()) {
		return Token{}, errors.New("Token expired")
	}
	return claims.Token(), nil
}
