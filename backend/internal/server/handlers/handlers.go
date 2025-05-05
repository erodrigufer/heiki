package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/erodrigufer/serenitynow/internal/state"
)

type Handlers struct {
	infoLog        *slog.Logger
	errorLog       *slog.Logger
	sm             state.StateManager
	sessionManager *scs.SessionManager
}

// NewHandlers creates a struct that contains all handlers of the application.
func NewHandlers(infoLog, errorLog *slog.Logger, db *sql.DB, sessions *scs.SessionManager) *Handlers {
	handlers := new(Handlers)
	handlers.infoLog = infoLog
	handlers.errorLog = errorLog
	handlers.sm = *state.NewStateManager(db)
	handlers.sessionManager = sessions

	return handlers
}
