package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/erodrigufer/serenitynow/internal/state"
)

type Handlers struct {
	infoLog  *slog.Logger
	errorLog *slog.Logger
	sm       state.StateManager
}

// NewHandlers creates a struct that contains all handlers of the application.
func NewHandlers(infoLog, errorLog *slog.Logger, db *sql.DB) *Handlers {
	handlers := new(Handlers)
	handlers.infoLog = infoLog
	handlers.errorLog = errorLog
	handlers.sm = *state.NewStateManager(db)

	return handlers
}
