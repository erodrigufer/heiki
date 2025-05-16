package state

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/erodrigufer/serenitynow/internal/types"
)

//go:embed insertTask.sql
var insertTaskQuery string

//go:embed getAllTasks.sql
var getAllTasksQuery string

//go:embed updateCompletedTask.sql
var updateCompletedTaskQuery string

func (sm *StateManager) InsertTask(ctx context.Context, priority, description, dueDate string) error {
	var dueDatePtr *string
	// If dueDate is an empty string, store a NULL value.
	if dueDate == "" {
		dueDatePtr = nil
	} else {
		dueDatePtr = &dueDate
	}
	_, err := sm.ExecContext(ctx, insertTaskQuery, priority, description, dueDatePtr, nil)
	if err != nil {
		return fmt.Errorf("unable to insert task in db: %w", err)
	}
	return nil
}

func (sm *StateManager) GetAllTasks(ctx context.Context) ([]types.Task, error) {
	rows, err := sm.QueryContext(ctx, getAllTasksQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get all tasks from db: %w", err)
	}
	allTasks := make([]types.Task, 0)
	defer rows.Close()

	var createdAtStr *string
	var completedAtStr *string
	var dueAtStr *string

	for rows.Next() {
		t := &types.Task{}
		err := rows.Scan(&t.ID, &t.Completed, &t.Priority, &t.Description, &createdAtStr,
			&completedAtStr, &dueAtStr)
		if err != nil {
			return nil, fmt.Errorf("row.Scan() failed: %w", err)
		}
		tParsed, err := parseDatesIntoTask(*t, createdAtStr, completedAtStr, dueAtStr)
		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, tParsed)
	}
	return allTasks, nil
}

func (sm *StateManager) UpdateCompletedTask(ctx context.Context, completed bool, id int) error {
	var err error
	if completed {
		completedAt := time.Now().Format("2006-01-02")
		_, err = sm.ExecContext(ctx, updateCompletedTaskQuery, completed, completedAt, id)
	} else {
		_, err = sm.ExecContext(ctx, updateCompletedTaskQuery, completed, nil, id)
	}
	if err != nil {
		return fmt.Errorf("unable to update completed column of task: %w", err)
	}
	return nil
}

func parseSqliteDate(date *string) (*time.Time, error) {
	if date == nil {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", *date)
	if err != nil {
		return nil, fmt.Errorf("unable to parse date: %w", err)
	}
	return &parsedDate, nil
}

func parseDatesIntoTask(task types.Task, createdAt, completedAt, dueAt *string) (types.Task, error) {
	var err error
	task.CreatedAt, err = parseSqliteDate(createdAt)
	if err != nil {
		return types.Task{}, err
	}
	task.CompletedAt, err = parseSqliteDate(completedAt)
	if err != nil {
		return types.Task{}, err
	}
	task.DueAt, err = parseSqliteDate(dueAt)
	if err != nil {
		return types.Task{}, err
	}

	return task, nil
}
