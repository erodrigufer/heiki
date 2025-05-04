INSERT INTO tasks 
  (completed, priority, description, createdAt, dueAt, completedAt) 
VALUES 
  (FALSE, ?, ?, date('now'), ?, ?); 
