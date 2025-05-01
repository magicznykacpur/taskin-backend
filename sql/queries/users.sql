-- name: CreateUser :exec
INSERT INTO users(id, created_at, updated_at, username, hashed_password)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUsers :many
SELECT * FROM users;