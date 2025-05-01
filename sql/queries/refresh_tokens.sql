-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens(user_id, token, created_at, updated_at, is_revoked, expires_at)
VALUES (?, ?, ?, ?, 0, ?);

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET is_revoked = 1 WHERE user_id = ?;

-- name: GetValidRefreshTokenForUserId :one
SELECT * FROM refresh_tokens WHERE is_revoked = 0 AND expires_at > ? AND user_id = ?;