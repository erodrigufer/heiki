package handlers

import (
	"log/slog"
)

type Handlers struct {
	infoLog  *slog.Logger
	errorLog *slog.Logger
}

// NewHandlers creates a struct that contains all handlers of the application.
func NewHandlers(infoLog, errorLog *slog.Logger) *Handlers {
	handlers := new(Handlers)
	handlers.infoLog = infoLog
	handlers.errorLog = errorLog

	return handlers
}
