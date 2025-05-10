-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  user_uuid UUID PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password_hash TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
