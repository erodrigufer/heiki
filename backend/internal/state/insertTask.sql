INSERT INTO tasks 
  (completed, priority, description, created_at, due_at, completed_at) 
VALUES 
  (FALSE, ?, ?, date('now'), ?, ?); 
