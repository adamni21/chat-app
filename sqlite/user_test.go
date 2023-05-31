package sqlite_test

import (
	"context"
	"testing"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/sqlite"
)

func TestCreate(t *testing.T) {
	t.Run("can create user", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s := sqlite.NewUserService(db)
		ctx := context.Background()

		user := &goChat.User{
			Username: "user0",
			Email:    "mail@mail.com",
			Verified: false,
		}

		if err := s.Create(ctx, user, "password"); err != nil {
			t.Fatal(err)
		}
	})
}

func InitUserService(t testing.TB) (goChat.UserService, func(), context.Context) {
	t.Helper()

	db := MustOpenDB(t)

	s := sqlite.NewUserService(db)
	ctx := context.Background()

	return s, func() { MustCloseDB(t, db) }, ctx
}
