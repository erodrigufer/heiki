package state

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/erodrigufer/serenitynow/internal/types"
)

//go:embed insertProject.sql
var insertProjectQuery string

//go:embed getAllProjects.sql
var getAllProjectsQuery string

func (sm *StateManager) InsertProject(ctx context.Context, name string) error {
	_, err := sm.ExecContext(ctx, insertProjectQuery, name)
	if err != nil {
		return fmt.Errorf("unable to insert project in db: %w", err)
	}
	return nil
}

func (sm *StateManager) GetAllProjects(ctx context.Context) ([]types.Project, error) {
	rows, err := sm.QueryContext(ctx, getAllProjectsQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get all projects from db: %w", err)
	}
	allProjects := make([]types.Project, 0)
	defer rows.Close()

	var createdAtStr *string

	for rows.Next() {
		p := types.Project{}
		err := rows.Scan(&p.ID, &p.Name, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("row.Scan() failed: %w", err)
		}
		createdAt, err := parseSqliteDate(createdAtStr)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		allProjects = append(allProjects, p)
	}
	return allProjects, nil
}
