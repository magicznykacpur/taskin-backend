-- name: CreateUser :exec
INSERT INTO users(id, created_at, updated_at, email, username, hashed_password)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUsers :many
SELECT * FROM users;

-- name: UpdateUserByID :one
UPDATE users
SET email = ?, username = ?, hashed_password = ?, updated_at = ?
WHERE id = ?
RETURNING *;