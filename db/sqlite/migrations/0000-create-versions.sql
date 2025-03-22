-- Up
CREATE TABLE versions (
  uuid TEXT PRIMARY KEY,
  version TEXT NOT NULL UNIQUE,
  migrated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  status TEXT NOT NULL DEFAULT 'pending'
);

-- Down
DROP TABLE versions;
