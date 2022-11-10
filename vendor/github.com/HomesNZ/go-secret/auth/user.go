package auth

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	RoleDev             Role = "dev"
	RoleTest            Role = "test"
	RoleAdmin           Role = "admin"
	RoleBranchAdmin     Role = "branch_admin"
	RoleUser            Role = "user"
	RoleAgent           Role = "agent"
	RolePremiumAgent    Role = "premium_agent"
	RoleBranchPrincipal Role = "branch_principal"
	RoleRayWhiteData    Role = "ray_white_data"
)

type Role string

type User struct {
	UserID    int
	FirstName string
	LastName  string
	Nickname  string
	Name      string
	Email     string
	Picture   string

	Roles []Role

	UpdatedAt time.Time
}

func (u User) HasRole(role Role) bool {
	for _, v := range u.Roles {
		if v == role {
			return true
		}
	}
	return false
}

// DEPRECATED in favour of auth.ContextToken()
func UserFromContext(c context.Context) (*User, error) {
	token, err := TokenFromContext(c)
	if err != nil {
		return nil, err
	}

	if token.User == nil {
		return nil, errors.New("Unable to retrieve user from token")
	}

	return token.User, nil
}

// DEPRECATED deprecated in favour of auth.ContextToken()
func UserFromHTTPRequest(r *http.Request) (*User, error) {
	return UserFromContext(r.Context())
}
