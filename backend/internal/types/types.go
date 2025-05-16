package types

import "time"

type Task struct {
	ID          int
	Completed   bool
	Priority    string
	Description string
	CreatedAt   *time.Time
	CompletedAt *time.Time
	DueAt       *time.Time
	// Projects are semantically identified by a `+`.
	Projects []string
	// Contexts are semantically identified by a `@`.
	Contexts []string
}
