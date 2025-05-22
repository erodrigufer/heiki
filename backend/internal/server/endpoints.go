package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/server/handlers"
	"github.com/erodrigufer/serenitynow/internal/server/middlewares"
	"github.com/erodrigufer/serenitynow/internal/static"
)

func (app *Application) defineEndpoints() (http.Handler, error) {
	mux := http.NewServeMux()

	disableAuth, err := app.GetConfigValueBool("DISABLE_AUTH")
	if err != nil {
		disableAuth = false
	}
	app.InfoLog.Info("Configuring authentication", slog.Bool("DISABLE_AUTH", disableAuth))
	m := middlewares.NewMiddlewares(app.InfoLog, app.ErrorLog, app.sessionManager, disableAuth)

	h := handlers.NewHandlers(app.InfoLog, app.ErrorLog, app.db, app.sessionManager)

	fileServer := http.FileServer(http.FS(static.STATIC_CONTENT))

	middlewareChain := func(next http.Handler) http.Handler {
		return m.RecoverPanic(m.SecureHeaders(m.LogRequest(app.sessionManager.LoadAndSave(next))))
	}

	authorizedUsername, err := app.GetConfigValueString("AUTH_USERNAME")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve authorized username: %w", err)
	}

	authorizedPassword, err := app.GetConfigValueString("AUTH_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve authorized password: %w", err)
	}

	mux.Handle("GET /login", m.AuthenticateLogin(h.GetLogin()))
	mux.Handle("POST /login", h.PostLogin(authorizedUsername, authorizedPassword))
	mux.Handle("POST /logout", h.PostLogout())
	mux.Handle("GET /content/", fileServer)
	mux.Handle("GET /api/v1/health", h.GetHealth())

	protectedMux := http.NewServeMux()
	mux.Handle("/", m.Authenticate(protectedMux))

	protectedMux.Handle("GET /", h.GetHome())
	protectedMux.Handle("GET /tasks", h.GetTasks())
	protectedMux.Handle("POST /tasks", h.PostTasks())
	protectedMux.Handle("PUT /tasks/{id}", h.PutCompletedTask())
	protectedMux.Handle("GET /projects", h.GetProjects())
	protectedMux.Handle("GET /projects/{id}", h.GetAllTasksByProjectID())
	protectedMux.Handle("POST /projects", h.PostProjects())
	protectedMux.Handle("GET /contexts", h.GetContexts())
	protectedMux.Handle("POST /contexts", h.PostContexts())

	return middlewareChain(mux), nil
}
