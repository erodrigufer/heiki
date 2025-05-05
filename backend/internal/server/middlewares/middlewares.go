package middlewares

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
)

type Middlewares struct {
	infoLog               *slog.Logger
	errorLog              *slog.Logger
	sessionManager        *scs.SessionManager
	disableAuthentication bool
}

// NewMiddlewares creates a struct that contains all middlewares of the application.
func NewMiddlewares(infoLog, errorLog *slog.Logger, sm *scs.SessionManager, disableAuthentication bool) *Middlewares {
	middlewares := new(Middlewares)
	middlewares.infoLog = infoLog
	middlewares.errorLog = errorLog
	middlewares.sessionManager = sm
	middlewares.disableAuthentication = disableAuthentication

	return middlewares
}
