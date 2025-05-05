package middlewares

import (
	"net/http"
)

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.disableAuthentication {
			userID := m.sessionManager.GetString(r.Context(), "userID")
			if userID == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// AuthenticateLogin checks in the login route if the user already has a
// valid auth token. If so, it redirects the user to the login page.
func (m *Middlewares) AuthenticateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := m.sessionManager.GetString(r.Context(), "userID")
		if userID != "" || m.disableAuthentication {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
