package server

import (
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/server/handlers"
	"github.com/erodrigufer/serenitynow/internal/server/middlewares"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	m := middlewares.NewMiddlewares(app.InfoLog, app.ErrorLog)

	h := handlers.NewHandlers(app.InfoLog, app.ErrorLog, app.db)

	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

	middlewareChain := func(next http.Handler) http.Handler {
		return m.RecoverPanic(m.SecureHeaders(m.LogRequest(next)))
	}

	mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /api/v1/health", h.GetHealth())
	mux.Handle("GET /", h.GetHome())
	mux.Handle("GET /tasks", h.GetTasks())
	mux.Handle("POST /tasks", h.PostTasks())

	return middlewareChain(mux)
}
