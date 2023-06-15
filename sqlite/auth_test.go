package sqlite_test

import (
	"context"
	"testing"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/sqlite"
)

func TestVerifyUser(t *testing.T) {
	authService, db, closeDB, ctx := InitAuthService(t, nil, nil)
	userService, _, _, _ := InitUserService(t, db)
	defer closeDB()

	password := "password"
	user := &goChat.User{
		Username: "user0",
		Email:    "test@mail.io",
	}
	MustCreateUser(t, ctx, userService, user, password)

	// return true if password is correct
	t.Run("password is correct", func(t *testing.T) {
		isCorrect, err := authService.VerifyUser(ctx, *user, password)
		if err != nil {
			t.Fatal("expected no error")
		}
		if !isCorrect {
			t.Fatal("expected true")
		}
	})

	// return false if password is incorrect
	t.Run("password is incorrect", func(t *testing.T) {
		isCorrect, err := authService.VerifyUser(ctx, *user, "wrong")
		if err != nil {
			t.Fatal("expected no error")
		}
		if isCorrect {
			t.Fatal("expected false")
		}
	})

	// return ENotFound if user doesn't exist
	t.Run("user doesn't exist", func(t *testing.T) {
		user := goChat.User{Id: -1}
		isCorrect, err := authService.VerifyUser(ctx, user, "password")
		if isCorrect {
			t.Fatal("expected false")
		}
		if val, ok := err.(goChat.Error); ok && val.ErrCode() != goChat.ENotFound {
			t.Fatalf("expected error code %d got %d", goChat.ENotFound, val.ErrCode())
		}
	})
}

func InitAuthService(t testing.TB, db *sqlite.DB, ctx context.Context) (goChat.AuthService, *sqlite.DB, func(), context.Context) {
	t.Helper()
	if db == nil {
		db = MustOpenDB(t)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	s := sqlite.NewAuthService(db)
	return s, db, func() { MustCloseDB(t, db) }, ctx
}
