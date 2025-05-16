package handlers

import (
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/views"
	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) GetProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := h.sm.GetAllProjects(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
		err = web.RenderComponent(r.Context(), w, views.ProjectsPageView(projects))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}

func (h *Handlers) PostProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			web.HandleBadRequest(w, "unable to parse form")
			return
		}

		name := r.PostForm.Get("name")

		err = h.sm.InsertProject(r.Context(), name)
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		http.Redirect(w, r, "/projects", http.StatusSeeOther)
	}
}
