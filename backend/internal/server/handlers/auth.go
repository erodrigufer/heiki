package handlers

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/views"
	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) GetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := web.RenderComponent(r.Context(), w, views.Login())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
		}
	}
}

func (h *Handlers) PostLogin(authorizedUsername, authorizedPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			web.HandleBadRequest(w, fmt.Sprintf("unable to parse form: %s", err.Error()))
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if username != authorizedUsername || password != authorizedPassword {
			web.SendHTMXErrorMessage(w, r, http.StatusUnauthorized, h.errorLog,
				`<p class="error-response"><b>Username</b> and/or <b>password</b> are invalid.</p>`)
			return
		}

		// Renew the session token before making the privilege-level change.
		err = h.sessionManager.RenewToken(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Make the privilege-level change.
		h.sessionManager.Put(r.Context(), "userID", username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handlers) PostLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionManager.Destroy(r.Context())
		if err != nil {
			web.HandleServerError(w, r, fmt.Errorf("unable to destroy session: %w", err), h.errorLog)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
