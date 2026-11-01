CREATE TABLE IF NOT EXISTS projects (
  id            TEXT        PRIMARY KEY,
  project_url   TEXT        UNIQUE,
  title         TEXT        NOT NULL CHECK (length(trim(title)) > 0),
  description   TEXT        NOT NULL DEFAULT '',
  languages     TEXT        NOT NULL DEFAULT '[]',
  created_at    TEXT        NOT NULL,
  visible       BOOLEAN     NOT NULL DEFAULT 1,
  custom        BOOLEAN     NOT NULL DEFAULT 0,
  body          TEXT        NOT NULL DEFAULT '',
)
