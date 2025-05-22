package handlers

import (
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/views"
	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) GetContexts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contexts, err := h.sm.GetAllContexts(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
		err = web.RenderComponent(r.Context(), w, views.ContextsPageView(contexts))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}

func (h *Handlers) PostContexts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			web.HandleBadRequest(w, "unable to parse form")
			return
		}

		name := r.PostForm.Get("name")

		err = h.sm.InsertContext(r.Context(), name)
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		http.Redirect(w, r, "/contexts", http.StatusSeeOther)
	}
}
