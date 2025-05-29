-- name: CreateTask :one
INSERT INTO tasks(id, created_at, updated_at, due_until, title, description, priority, category, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = ?;

-- name: GetTasksByTitle :many
SELECT * FROM tasks WHERE title LIKE ?;

-- name: GetTasksByDescription :many
SELECT * FROM tasks WHERE description LIKE ?;

-- name: GetTaskByTitleAndDescription :many
SELECT * FROM tasks WHERE title LIKE ? OR description LIKE ?;

-- name: GetAllUsersTasks :many
SELECT * FROM tasks WHERE user_id = ?;

-- name: UpdateTaskByID :one
UPDATE tasks
SET title = ?, description = ?, priority = ?, category = ?, updated_at = ?, due_until = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTaskByID :exec
DELETE FROM tasks WHERE id = ? AND user_id = ?;