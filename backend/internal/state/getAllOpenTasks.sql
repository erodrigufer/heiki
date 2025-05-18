SELECT t.id, t.completed, t.priority, t.description, t.created_at, 
    t.completed_at, t.due_at, group_concat(p.name) AS project_names, 
    group_concat(p.id) AS project_ids 
  FROM tasks AS t 
    LEFT JOIN projects_by_task AS pt ON t.id = pt.task_id 
    LEFT JOIN projects AS p ON pt.project_id = p.id 
  WHERE t.completed = 0
  GROUP BY t.id
  ORDER BY t.completed, t.priority, t.due_at ASC NULLS LAST, t.description;
