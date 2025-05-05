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
	// task1 := tasks.Task{Priority: "B", Description: "Buy a new notebook"}
	//
	// task6 := tasks.Task{Priority: "C", Description: "Find a new スマホ"}
	// task2 := tasks.Task{Priority: "A", Description: "Insurance claim", Contexts: []string{"5分", "メール"}}
	//
	// task3 := tasks.Task{Completed: true, Description: "Send gift to X", Contexts: []string{"仕事"}}
	// task4 := tasks.Task{Description: "Implement MVP backend template renderer", Contexts: []string{"programming"}, Projects: []string{"serenitynow"}}
	//
	// task5 := tasks.Task{Priority: "B", Description: "Buy train ticket", Contexts: []string{"Conference25", "仕事"}}
	//
	// tasks := []tasks.Task{task1, task2, task3, task4, task5, task6}
	//
	// tasks = sortTasksByPriority(tasks)

	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.sm.GetAllTasks(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
		err = web.RenderComponent(r.Context(), w, views.TasksPageView(tasks))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}

// func sortTasksByPriority(ts []tasks.Task) []tasks.Task {
// 	cmpFn := func(a, b tasks.Task) int {
// 		// Tasks with no Priority field should be sorted to the bottom.
// 		if a.Priority == "" {
// 			return 1
// 		}
// 		if b.Priority == "" {
// 			return -1
// 		}
// 		return cmp.Compare(a.Priority, b.Priority)
// 	}
// 	slices.SortStableFunc(ts, cmpFn)
// 	return ts
// }

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

		err = h.sm.InsertTask(r.Context(), priority, description, dueDateStr)
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

		tasks, err := h.sm.GetAllTasks(r.Context())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
		err = web.RenderComponent(r.Context(), w, views.TasksPageContent(tasks))
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
			return
		}
	}
}
