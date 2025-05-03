-- name: CreateTask :one
INSERT INTO tasks(id, created_at, updated_at, title, description, priority, category, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = ?;

-- name: GetTaskByTitle :one
SELECT * FROM tasks WHERE UPPER(title) LIKE '%?%';

-- name: GetTaskByDescription :one
SELECT * FROM tasks WHERE UPPER(description) LIKE '%?%';

-- name: GetAllUsersTasks :many
SELECT * FROM tasks WHERE user_id = ?;

-- name: UpdateTaskByID :one
UPDATE tasks
SET title = ?, description = ?, priority = ?, category = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTaskByID :exec
DELETE FROM tasks WHERE id = ?;