-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY,
    admin_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id)
      REFERENCES admins (id)
      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS sessions;
