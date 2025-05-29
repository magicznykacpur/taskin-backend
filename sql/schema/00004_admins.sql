-- +goose Up
ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN is_admin;