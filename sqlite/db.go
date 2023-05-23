package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
	if err := db.migrate(); err != nil {
		return err
	}

	return nil
}

type migration = struct {
	fileName  string
	timestamp int32
}

func (db *DB) migrate() error {
	currentUserVersion, err := db.currentUserVersion()
	if err != nil {
		return fmt.Errorf("getting current user_version: %w", err)
	}

	fileNames, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return err
	}

	files := make([]migration, 0, len(fileNames))
	for _, name := range fileNames {
		rawTimestamp := strings.SplitN(path.Base(name), "_", 2)[0]
		if len(rawTimestamp) != 10 {
			return fmt.Errorf("timestamp '%s' has invalid length. filename must begin with left padded timestamp", rawTimestamp)
		}

		// sqlite user_version is a 32 bit int
		timestamp, err := strconv.ParseInt(rawTimestamp, 10, 32)
		if err != nil {
			return fmt.Errorf("parsing timestamp '%s' gave error: %w", rawTimestamp, err)
		}

		if int32(timestamp) > currentUserVersion {
			files = append(files, migration{fileName: name, timestamp: int32(timestamp)})
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i].timestamp < files[j].timestamp })
	for _, migration := range files {
		if err := db.migrateFile(migration); err != nil {
			return fmt.Errorf("migrating file '%s': %w", migration.fileName, err)
		}
	}

	return nil
}

func (db *DB) migrateFile(migration migration) error {
	buf, err := fs.ReadFile(migrationsFS, migration.fileName)
	if err != nil {
		return err
	}

	if _, err := db.db.Exec(string(buf)); err != nil {
		return err
	}

	return nil
}

func (db *DB) currentUserVersion() (int32, error) {
	query, err := db.db.Query("PRAGMA user_version;")
	if err != nil {
		return 0, err
	}
	defer query.Close()

	var userVersion int32
	query.Next()
	if err := query.Scan(&userVersion); err != nil {
		return 0, err
	}

	return userVersion, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
