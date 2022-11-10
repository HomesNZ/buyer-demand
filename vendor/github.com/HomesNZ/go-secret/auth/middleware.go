package auth

import (
	"net/http"

	"github.com/justinas/alice"
)

// UserMiddleware requires the request to contain a valid JWT token within the
// Authorization header.
func (a Auth) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, success := a.CheckRequest(r)
		_, err := UserFromHTTPRequest(r)
		if success && err == nil {
			if next != nil {
				next.ServeHTTP(w, r)
			}
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

// ServiceMiddleware requires that the request is being sent on behalf of a
// microservice, verified by the Authorization header.
func (a Auth) ServiceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, success := a.CheckRequest(r)
		_, err := ServiceFromHTTPRequest(r)
		if success && err == nil {
			if next != nil {
				next.ServeHTTP(w, r)
			}
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

// UserOrServiceMiddleware requires the request to contain a valid JWT token
// within the Authorization header, or that the request has been sent on behalf
// of a microservice, verified by the Authorization header.
func (a Auth) UserOrServiceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, success := a.CheckRequest(r)
		t, err := TokenFromHTTPRequest(r)
		if success && err == nil && (t.User != nil || t.Service != nil) {
			if next != nil {
				next.ServeHTTP(w, r)
			}
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// OptionalUserMiddleware allows user authentication, or no authentication.
func (a Auth) OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, _ = a.CheckRequest(r)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// OptionalServiceMiddleware allows service authentication, or no authentication.
func (a Auth) OptionalServiceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, _ = a.CheckRequest(r)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// OptionalUserOrServiceMiddleware allows service authentication, user
// authentication, or no authentication.
func (a Auth) OptionalUserOrServiceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, _ = a.CheckRequest(r)
		// Always execute next even if authentication fails, as this is an optional
		// auth handler.
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// RoleMiddleware returns a middleware which checks to see if the requesting
// user matches at least one of the role lists.
//
// For a match to be successful, the user must have *all* of the listed roles.
//
// Example to restrict access users with either the 'admin' role, or 'agent' AND
// 'premium_agent':
// a.RoleMiddleware(
//   []auth.Role{auth.RoleAdmin},
//   []auth.Role{auth.RoleAgent, auth.RolePremiumAgent},
// )
func (a Auth) RoleMiddleware(roleLists ...[]Role) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := UserFromHTTPRequest(r)
			// TODO: should this actually validate the user?
			if err != nil {
				a.logger.Debug(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			matched := false
			for _, roles := range roleLists {
				matched = true
				for _, role := range roles {
					if !u.HasRole(role) {
						matched = false
						break
					}
				}
				if matched {
					break
				}
			}

			if !matched {
				a.logger.Debug("user did not match any role lists")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}
