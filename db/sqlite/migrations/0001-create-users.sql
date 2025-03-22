-- Up
CREATE TABLE users (
  uuid TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password_hash TEXT NOT NULL
);

ALTER TABLE users ADD INDEX idx_username (username);

-- Down
DROP TABLE users;
