SELECT id, completed, priority, description, createdAt, completedAt, dueAt
  FROM tasks
  ORDER BY completed, priority, description;
