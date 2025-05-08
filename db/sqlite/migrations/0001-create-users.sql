-- Up
CREATE TABLE users (
  uuid TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password_hash TEXT
);

-- Down
DROP TABLE users;
