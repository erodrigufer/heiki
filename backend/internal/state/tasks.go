package state

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/erodrigufer/serenitynow/internal/types"
)

//go:embed insertTask.sql
var insertTaskQuery string

//go:embed getAllTasks.sql
var getAllTasksQuery string

//go:embed getAllOpenTasks.sql
var getAllOpenTasksQuery string

//go:embed updateCompletedTask.sql
var updateCompletedTaskQuery string

//go:embed insertIntoProjectsByTask.sql
var insertIntoProjectsByTask string

//go:embed insertIntoContextsByTask.sql
var insertIntoContextsByTask string

//go:embed getAllTasksByProjectID.sql
var getAllTasksByProjectID string

//go:embed getAllTasksByContextID.sql
var getAllTasksByContextID string

func (sm *StateManager) InsertTask(ctx context.Context, priority, description, dueDate string, projectID, contextID int) error {
	var dueDatePtr *string
	// If dueDate is an empty string, store a NULL value.
	if dueDate == "" {
		dueDatePtr = nil
	} else {
		dueDatePtr = &dueDate
	}

	tx, err := sm.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction to insert task: %w", err)
	}
	defer tx.Rollback()

	results, err := tx.ExecContext(ctx, insertTaskQuery, priority, description, dueDatePtr, nil)
	if err != nil {
		return fmt.Errorf("unable to insert task in db: %w", err)
	}

	taskID, err := results.LastInsertId()
	if err != nil {
		return fmt.Errorf("unable to retrieve ID of inserted task: %w", err)
	}

	if projectID != 0 {
		_, err = tx.ExecContext(ctx, insertIntoProjectsByTask, int(taskID), projectID)
		if err != nil {
			return fmt.Errorf("unable to insert task in db: %w", err)
		}
	}

	if contextID != 0 {
		_, err = tx.ExecContext(ctx, insertIntoContextsByTask, int(taskID), contextID)
		if err != nil {
			return fmt.Errorf("unable to insert task in db: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}
	return nil
}

func (sm *StateManager) GetAllTasks(ctx context.Context, showCompletedTasks bool) ([]types.Task, error) {
	var rows *sql.Rows
	var err error
	if showCompletedTasks {
		rows, err = sm.QueryContext(ctx, getAllTasksQuery)
		if err != nil {
			return nil, fmt.Errorf("unable to get all tasks from db: %w", err)
		}
	} else {
		rows, err = sm.QueryContext(ctx, getAllOpenTasksQuery)
		if err != nil {
			return nil, fmt.Errorf("unable to get all tasks from db: %w", err)
		}

	}
	allTasks, err := parseRowsToTasks(rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tasks from rows: %w", err)
	}

	return allTasks, nil
}

func (sm *StateManager) GetAllTasksByProjectID(ctx context.Context, projectID int) ([]types.Task, error) {
	rows, err := sm.QueryContext(ctx, getAllTasksByProjectID, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to get all tasks from db: %w", err)
	}
	allTasks, err := parseRowsToTasks(rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tasks from rows: %w", err)
	}

	return allTasks, nil
}

func (sm *StateManager) GetAllTasksByContextID(ctx context.Context, contextID int) ([]types.Task, error) {
	rows, err := sm.QueryContext(ctx, getAllTasksByContextID, contextID)
	if err != nil {
		return nil, fmt.Errorf("unable to get all tasks from db: %w", err)
	}
	allTasks, err := parseRowsToTasks(rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tasks from rows: %w", err)
	}

	return allTasks, nil
}

func parseRowsToTasks(rows *sql.Rows) ([]types.Task, error) {
	allTasks := make([]types.Task, 0)
	defer rows.Close()

	var createdAtStr *string
	var completedAtStr *string
	var dueAtStr *string
	var projectNamesStr *string
	var projectIDsStr *string
	var contextNamesStr *string
	var contextIDsStr *string

	for rows.Next() {
		t := &types.Task{}
		err := rows.Scan(&t.ID, &t.Completed, &t.Priority, &t.Description, &createdAtStr,
			&completedAtStr, &dueAtStr, &projectNamesStr, &projectIDsStr,
			&contextNamesStr, &contextIDsStr)
		if err != nil {
			return nil, fmt.Errorf("row.Scan() failed: %w", err)
		}
		tParsed, err := parseDatesIntoTask(*t, createdAtStr, completedAtStr, dueAtStr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse dates into task: %w", err)
		}
		projects, err := parseProjects(projectNamesStr, projectIDsStr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse projects: %w", err)
		}
		contexts, err := parseContexts(contextNamesStr, contextIDsStr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse contexts: %w", err)
		}

		tParsed.Projects = projects
		tParsed.Contexts = contexts
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

func parseProjects(concatenatedNames, concatenatedIDs *string) ([]types.Project, error) {
	if concatenatedNames == nil {
		return []types.Project{}, nil
	}
	names := strings.Split(*concatenatedNames, ",")
	if len(names) == 0 {
		return []types.Project{}, nil
	}
	if concatenatedIDs == nil {
		return []types.Project{}, nil
	}
	ids := strings.Split(*concatenatedIDs, ",")
	if len(ids) == 0 {
		return []types.Project{}, nil
	}
	projects := make([]types.Project, 0)
	for i, name := range names {
		idInt, err := strconv.Atoi(ids[i])
		if err != nil {
			return []types.Project{}, fmt.Errorf("unable to convert id of project into int: %w", err)
		}
		project := types.Project{
			ID:   idInt,
			Name: name,
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func parseContexts(concatenatedNames, concatenatedIDs *string) ([]types.Context, error) {
	if concatenatedNames == nil {
		return []types.Context{}, nil
	}
	names := strings.Split(*concatenatedNames, ",")
	if len(names) == 0 {
		return []types.Context{}, nil
	}
	if concatenatedIDs == nil {
		return []types.Context{}, nil
	}
	ids := strings.Split(*concatenatedIDs, ",")
	if len(ids) == 0 {
		return []types.Context{}, nil
	}
	contexts := make([]types.Context, 0)
	for i, name := range names {
		idInt, err := strconv.Atoi(ids[i])
		if err != nil {
			return []types.Context{}, fmt.Errorf("unable to convert id of context into int: %w", err)
		}
		context := types.Context{
			ID:   idInt,
			Name: name,
		}
		contexts = append(contexts, context)
	}
	return contexts, nil
}
