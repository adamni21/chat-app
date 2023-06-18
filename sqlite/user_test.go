package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/sqlite"
)

func TestCreate(t *testing.T) {
	t.Run("can create user", func(t *testing.T) {
		s, _, closeDB, ctx := InitUserService(t)
		defer closeDB()

		username, email, verified := "user0", "mail@mail.com", false

		user := &goChat.User{
			Username: username,
			Email:    email,
			Verified: verified,
		}

		if err := s.Create(ctx, user, "password"); err != nil {
			t.Fatal(err)
		} else if user.Id != 1 {
			t.Fatalf("Username=%d, want %d", user.Id, 1)
		} else if user.Username != username {
			t.Fatalf("Username=%s, want %s", user.Username, username)
		} else if user.Email != email {
			t.Fatalf("Email=%s, want %s", user.Email, email)
		} else if user.Verified != verified {
			t.Fatalf("Verified=%t, want %t", user.Verified, verified)
		} else if user.CreatedAt.IsZero() {
			t.Fatal("expected created at")
		} else if user.UpdatedAt.IsZero() {
			t.Fatal("expected updated at")
		}

		persistedUser, err := s.FindById(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(user, persistedUser) {
			t.Fatalf("persisted user=%#v, want %#v", persistedUser, user)
		}
	})
}

func TestFindById(t *testing.T) {
	t.Run("can find user by id", func(t *testing.T) {
		s, _, closeDB, ctx := InitUserService(t)
		defer closeDB()

		user := &goChat.User{
			Username: "user0",
			Email:    "mail@mail.com",
			Verified: false,
		}
		MustCreateUser(t, ctx, s, user, "password")

		if foundUser, err := s.FindById(ctx, 1); err != nil {
			t.Fatal(err)
		} else if foundUser.Username != user.Username {
			t.Fatalf("Username=%s, want %s", foundUser.Username, user.Email)
		} else if foundUser.Email != user.Email {
			t.Fatalf("Email=%s, want %s", foundUser.Email, user.Email)
		} else if foundUser.Verified != user.Verified {
			t.Fatalf("Verified=%t, want %t", foundUser.Verified, user.Verified)
		}
	})
}

// pass shared db if used by multiple services, otherwise pass nil
func InitUserService(t testing.TB) (goChat.UserService, *sqlite.DB, func(), context.Context) {
	t.Helper()
	db := MustOpenDB(t)
	s := sqlite.NewUserService(db)
	ctx := context.Background()
	return s, db, func() { MustCloseDB(t, db) }, ctx
}

func MustCreateUser(t testing.TB, ctx context.Context, s goChat.UserService, user *goChat.User, password string) *goChat.User {
	t.Helper()
	if err := s.Create(ctx, user, "password"); err != nil {
		t.Fatal(err)
	}
	return user
}
