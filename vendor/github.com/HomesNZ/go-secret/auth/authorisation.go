package auth

import (
	"context"
	"net/http"
)

type Allower interface {
	Allow(ctx context.Context) bool
}

type Rules map[string]Allower

func (r Rules) Allow(ctx context.Context, name string) bool {
	allower, ok := r[name]
	if !ok {
		// TODO: should we panic or add error to result?
		return false
	}
	return allower.Allow(ctx)
}

type Authorisation struct {
	Authenticator RequestAuthenticator
	Rules         Rules
}

func (a Authorisation) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := a.Authenticator.Authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := WithRules(r.Context(), a.Rules)
		if t.Valid {
			ctx = SetContextToken(ctx, t)
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a Authorisation) MiddlewareAllow(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := a.Authenticator.Authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := WithRules(r.Context(), a.Rules)
		if t.Valid {
			ctx = SetContextToken(ctx, t)
		}
		r = r.WithContext(ctx)
		if Allow(ctx, name) {
			next.ServeHTTP(w, r)
			return
		}
		if t.Valid {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func WithRules(ctx context.Context, r Rules) context.Context {
	return context.WithValue(ctx, contextKeyRules, r)
}

func Allow(ctx context.Context, name string) bool {
	r, ok := ctx.Value(contextKeyRules).(Rules)
	if !ok {
		return false
	}
	return r.Allow(ctx, name)
}
