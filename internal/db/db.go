package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)


// DB wraps a SQLite connection.
type DB struct {
	conn *sql.DB
}

// Open opens the SQLite database and runs migrations.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

// Close closes the database.
func (db *DB) Close() error {
	return db.conn.Close()
}

// migrate runs the schema migrations.
func (db *DB) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS guild_settings (
    guild_id TEXT PRIMARY KEY,
    alert_channel_id TEXT,
    language TEXT NOT NULL DEFAULT 'en' CHECK(language IN ('en', 'tc', 'sc', 'bilingual')),
    tide_station TEXT NOT NULL DEFAULT 'QUB',
    bot_status_enabled BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS warning_state (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    subtype TEXT,
    action_code TEXT,
    issue_time TEXT,
    update_time TEXT,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tips_state (
    id INTEGER PRIMARY KEY,
    update_time TEXT NOT NULL,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_warning_state_code ON warning_state(code);
CREATE INDEX IF NOT EXISTS idx_tips_state_update_time ON tips_state(update_time);
`
	if _, err := db.conn.Exec(schema); err != nil {
		return err
	}
	return nil
}

// Exec executes a SQL statement.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

// QueryRow queries a single row.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// Query queries rows.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}
