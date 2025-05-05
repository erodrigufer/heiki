package handlers

import (
	"net/http"
)

func (h *Handlers) GetHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	}
}
