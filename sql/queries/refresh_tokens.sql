-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(user_id, created_at, updated_at, is_revoked, expires_at)
VALUES (?, ?, ?, 0, ?)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET is_revoked = 1 WHERE user_id = ?;