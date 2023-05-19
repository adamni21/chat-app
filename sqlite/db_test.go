package sqlite_test

import (
	"chat-app/sqlite"
	"testing"
)

func TestDB(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

func MustOpenDB(tb testing.TB) *sqlite.DB {
	tb.Helper()

	db := sqlite.NewDB(":memory:")
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}

	return db
}

func MustCloseDB(tb testing.TB, db *sqlite.DB) {
	tb.Helper()

	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}
