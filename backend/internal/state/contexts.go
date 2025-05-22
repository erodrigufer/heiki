package state

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/erodrigufer/serenitynow/internal/types"
)

//go:embed insertContext.sql
var insertContextQuery string

//go:embed getAllContexts.sql
var getAllContextsQuery string

func (sm *StateManager) InsertContext(ctx context.Context, name string) error {
	_, err := sm.ExecContext(ctx, insertContextQuery, name)
	if err != nil {
		return fmt.Errorf("unable to insert context in db: %w", err)
	}
	return nil
}

func (sm *StateManager) GetAllContexts(ctx context.Context) ([]types.Context, error) {
	rows, err := sm.QueryContext(ctx, getAllContextsQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get all contexts from db: %w", err)
	}
	allContexts := make([]types.Context, 0)
	defer rows.Close()

	var createdAtStr *string

	for rows.Next() {
		c := types.Context{}
		err := rows.Scan(&c.ID, &c.Name, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("row.Scan() failed: %w", err)
		}
		createdAt, err := parseSqliteDate(createdAtStr)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = createdAt
		allContexts = append(allContexts, c)
	}
	return allContexts, nil
}
