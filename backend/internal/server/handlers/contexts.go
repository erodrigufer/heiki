package handlers

import (
	"net/http"
	"strconv"

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

func (h *Handlers) GetAllTasksByContextID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextIDStr := r.PathValue("id")
		contextID, err := strconv.Atoi(contextIDStr)
		if err != nil {
			web.HandleBadRequest(w, "unable to parse context ID")
			return
		}
		tasks, err := h.sm.GetAllTasksByContextID(r.Context(), contextID)
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
		projects, err := h.sm.GetAllProjects(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		contexts, err := h.sm.GetAllContexts(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		err = web.RenderComponent(r.Context(), w, views.TasksPageView(tasks, projects, contexts))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}
