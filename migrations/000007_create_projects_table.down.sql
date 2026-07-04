ALTER TABLE todos DROP FOREIGN KEY fk_todos_project;
ALTER TABLE todos DROP COLUMN project_id;
DROP TABLE IF EXISTS projects;
