package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type appMetadata struct {
	UID   int    `json:"uid"`
	Roles []Role `json:"roles"`
}

type userMetadata struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type identity struct {
	UserID     string `json:"user_id"`
	Provider   string `json:"provider"`
	Connection string `json:"connection"`
	IsSocial   bool   `json:"isSocial"`
}

type MultiString []string

func (s *MultiString) UnmarshalJSON(b []byte) error {
	e := multiError{}

	// try to parse as single string
	var str string
	err := json.Unmarshal(b, &str)
	if err == nil {
		*s = append(*s, str)
		return nil
	}
	e = append(e, err)

	// try to parse as []string
	var strs []string
	err = json.Unmarshal(b, &strs)
	if err != nil {
		return append(e, err)
	}
	*s = MultiString(strs)
	return nil
}

type rawClaims struct {
	ISS string      `json:"iss"`
	Sub string      `json:"sub"`
	Aud MultiString `json:"aud"`
	Exp float64     `json:"exp"`
	IAT float64     `json:"iat"`

	Email             string        `json:"email"`
	Name              string        `json:"name"`
	Nickname          string        `json:"nickname"`
	Picture           string        `json:"picture"`
	UpdatedAt         time.Time     `json:"updated_at"`
	HomesAppMetadata  appMetadata   `json:"https://homes.co.nz/app_metadata"`
	HomesUserMetadata *userMetadata `json:"https://homes.co.nz/user_metadata"`
	HomesIdentities   []identity    `json:"https://homes.co.nz/identities"`
	HomesRoles        []Role        `json:"https://homes.co.nz/roles"`
}

func claimsFromContext(ctx context.Context) (Token, error) {
	token, ok := ContextToken(ctx)
	if !ok {
		return token, errors.New("token not found")
	}
	return token, nil
}

func (c rawClaims) Valid() error {
	if c.userID() > 0 && len(c.roles()) == 0 {
		return errors.New("expected at least 1 role, got: 0")
	}
	return nil
}

func (c rawClaims) userID() int {
	return c.HomesAppMetadata.UID
}

func (c rawClaims) roles() []Role {
	return c.HomesAppMetadata.Roles
}

func (c rawClaims) userMetadata() userMetadata {
	if c.HomesUserMetadata != nil {
		return *c.HomesUserMetadata
	}

	return userMetadata{}
}

func (c rawClaims) User() *User {
	if c.userID() == 0 {
		return nil
	}
	metadata := c.userMetadata()
	name := c.Name
	if name == "" {
		name = metadata.FirstName + " " + metadata.LastName
	}
	return &User{
		UserID:    c.userID(),
		FirstName: metadata.FirstName,
		LastName:  metadata.LastName,
		Nickname:  c.Nickname,
		Name:      name,
		Email:     c.Email,
		Picture:   c.Picture,
		Roles:     c.roles(),

		UpdatedAt: c.UpdatedAt,
	}
}

func (t rawClaims) Token() Token {
	return Token{
		Valid:    true,
		User:     t.User(),
		Audience: t.Aud,
	}
}
