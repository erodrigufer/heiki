-- migrate:up
ALTER TABLE projects
  RENAME COLUMN createdAT TO created_at;
ALTER TABLE contexts
  RENAME COLUMN createdAT TO created_at;
ALTER TABLE tasks
  RENAME COLUMN createdAT TO created_at;
ALTER TABLE tasks
  RENAME COLUMN dueAt TO due_at;
ALTER TABLE tasks
  RENAME COLUMN completedAt TO completed_at;

-- migrate:down
ALTER TABLE projects
  RENAME COLUMN created_at TO createdAT;
ALTER TABLE contexts
  RENAME COLUMN created_at TO createdAT;
ALTER TABLE tasks
  RENAME COLUMN created_at TO createdAt;
ALTER TABLE tasks
  RENAME COLUMN due_at TO dueAt;
ALTER TABLE tasks
  RENAME COLUMN completed_at TO completedAt;
