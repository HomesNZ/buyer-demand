package allow

import (
	"context"
	"strings"

	"github.com/HomesNZ/go-secret/auth"
)

type Never struct{}

func (Never) Allow(_ context.Context) bool {
	return false
}

type Always struct{}

func (Always) Allow(_ context.Context) bool {
	return true
}

// Every allows if every Allower allows
type Every []auth.Allower

func (e Every) Allow(ctx context.Context) bool {
	for _, allower := range e {
		if !allower.Allow(ctx) {
			return false
		}
	}
	return true
}

// Any allows if any Allower allows
type Any []auth.Allower

func (a Any) Allow(ctx context.Context) bool {
	for _, allower := range a {
		if allower.Allow(ctx) {
			return true
		}
	}
	return false
}

// Role allows if there is a user and they have the specified role
type Role auth.Role

func (r Role) Allow(ctx context.Context) bool {
	t, ok := auth.ContextToken(ctx)
	if !ok || t.User == nil {
		return false
	}
	return t.User.HasRole(auth.Role(r))
}

// Service allows if there is a service. The zero value matches any service name
type Service string

func (s Service) Allow(ctx context.Context) bool {
	t, ok := auth.ContextToken(ctx)
	if !ok || t.Service == nil {
		return false
	}
	return s == "" || strings.EqualFold(t.Service.Name, string(s))
}

// Audience allows if there is a audience. The zero value matches any audience
type Audience string

func (a Audience) Allow(ctx context.Context) bool {
	t, ok := auth.ContextToken(ctx)
	if !ok || len(t.Audience) == 0 {
		return false
	}
	if a == "" {
		return true
	}
	for _, audience := range t.Audience {
		if strings.EqualFold(audience, string(a)) {
			return true
		}
	}
	return false
}

// Authenticated allows if there is a token or not.
type Authenticated bool

func (a Authenticated) Allow(ctx context.Context) bool {
	t, _ := auth.ContextToken(ctx)
	return t.Valid == bool(a)
}
