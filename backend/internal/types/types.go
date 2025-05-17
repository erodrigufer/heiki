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
	Projects []Project
	// Contexts are semantically identified by a `@`.
	Contexts []Context
}

type Project struct {
	ID        int
	Name      string
	CreatedAt *time.Time
}

type Context struct {
	ID        int
	Name      string
	CreatedAt *time.Time
}
