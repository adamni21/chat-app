package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB

	DSN string
}

func NewDB(dsn string) *DB {
	return &DB{
		DSN: dsn,
	}
}

func (db *DB) Open() error {
	var err error
	if db.db, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}
	if err := db.db.Ping(); err != nil {
		return err
	}
	if _, err := db.db.Exec("PRAGMA journal_mode = wal;"); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}
	if _, err := db.db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
