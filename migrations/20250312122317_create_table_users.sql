-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id        INTEGER PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
