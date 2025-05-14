SELECT id, completed, priority, description, created_at, completed_at, due_at
  FROM tasks
  ORDER BY completed, priority, description;
