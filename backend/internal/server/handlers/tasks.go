package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/erodrigufer/serenitynow/internal/views"
	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) GetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		showCompletedTasks := queryValues.Has("completed")
		tasks, err := h.sm.GetAllTasks(r.Context(), showCompletedTasks)
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

func (h *Handlers) PostTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			web.HandleBadRequest(w, "unable to parse form")
			return
		}

		description := r.PostForm.Get("description")
		priority := r.PostForm.Get("priority")
		dueDateStr := r.PostForm.Get("due-date")
		projectIDStr := r.PostForm.Get("project-id")
		contextIDStr := r.PostForm.Get("context-id")

		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			web.HandleBadRequest(w, "project-id cannot be parsed to an int")
			return
		}

		contextID, err := strconv.Atoi(contextIDStr)
		if err != nil {
			web.HandleBadRequest(w, "context-id cannot be parsed to an int")
			return
		}

		err = h.sm.InsertTask(r.Context(), priority, description, dueDateStr, projectID, contextID)
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	}
}

func (h *Handlers) PutCompletedTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 0)
		if err != nil {
			web.HandleBadRequest(w, "unable to parse task id")
			return
		}

		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			web.HandleServerError(w, r, fmt.Errorf("unable to parse URL query: %w", err), h.errorLog)
			return
		}
		completedStr := values.Get("completed")
		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			web.HandleBadRequest(w, "unable to parse task completed status")
			return
		}

		err = h.sm.UpdateCompletedTask(r.Context(), completed, int(id))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}

		tasks, err := h.sm.GetAllTasks(r.Context(), true)
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

		err = web.RenderComponent(r.Context(), w, views.TasksPageContent(tasks, projects, contexts))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}
