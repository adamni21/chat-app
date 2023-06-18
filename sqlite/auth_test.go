package sqlite_test

import (
	"context"
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
	"github.com/adamni21/goChat/sqlite"
)

func TestLogin(t *testing.T) {
	authService, db, closeDB, ctx := InitAuthService(t)
	userService := sqlite.NewUserService(db)
	defer closeDB()

	password := "password"
	user := &goChat.User{
		Username: "user0",
		Email:    "test@mail.io",
	}
	MustCreateUser(t, ctx, userService, user, password)

	// can successfully login
	t.Run("login successfully", func(t *testing.T) {
		session, err := authService.Login(ctx, *user, password)
		if err != nil {
			t.Fatalf("expected no error got %+v", err)
		}

		if session.UserId != user.Id {
			t.Fatalf("userId incorrect, want %d got %d", user.Id, session.UserId)
		}
		if len(session.Id) != 24 {
			t.Fatalf("expected sessionId to be of len 24, sessionId %s", session.Id)
		}
		if session.Expiry.IsZero() {
			t.Fatalf("expected Expiry")
		}

		var persistedSession goChat.Session
		rows, err := db.QueryContext(ctx, "SELECT id, userId, expiry FROM sessions;")
		defer rows.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !rows.Next() {
			t.Fatal("couldn't find created session in DB")
		}
		err = rows.Scan(&persistedSession.Id, &persistedSession.UserId, (*sqlite.NullTime)(&persistedSession.Expiry))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(session, persistedSession) {
			t.Fatalf("session not stored correctly got %+v want %+v", persistedSession, session)
		}
	})

	// return EUnauthorized if password is not correct
	t.Run("wrong password", func(t *testing.T) {
		session, err := authService.Login(ctx, *user, "wrongPassword")
		if err == nil {
			t.Fatal("expected error")
		}
		if val, ok := err.(goChat.Error); val.ErrCode() != goChat.EUnauthorized || !ok {
			t.Fatalf("expected error code %d got %d", goChat.EUnauthorized, err)
		}
		if !reflect.DeepEqual(session, goChat.Session{}) {
			t.Fatalf("expected null session got %+v", session)
		}
	})

	// return ENotFound if user not exists
	t.Run("user doesn't exist", func(t *testing.T) {
		session, err := authService.Login(ctx, goChat.User{Id: -1}, "wrongPassword")
		if err == nil {
			t.Fatal("expected error")
		}
		if val, ok := err.(goChat.Error); val.ErrCode() != goChat.ENotFound || !ok {
			t.Fatalf("expected error code %d got %d", goChat.ENotFound, err)
		}
		if !reflect.DeepEqual(session, goChat.Session{}) {
			t.Fatalf("expected null session got %+v", session)
		}
	})
}

func TestVerifyUser(t *testing.T) {
	authService, db, closeDB, ctx := InitAuthService(t)
	userService := sqlite.NewUserService(db)
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

func MustCreateSession(t testing.TB, ctx context.Context, db *sqlite.DB, session *goChat.Session) {
	t.Helper()
	sessionId, err := crypto.GenerateRandomBytes(16)
	if err != nil {
		t.Fatal(err)
	}

	session.Id = goChat.SessionId(base64.URLEncoding.EncodeToString(sessionId))
	session.Expiry = time.Now().UTC().Add(30 * 24 * time.Hour).Truncate(time.Second)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	query := `
		INSERT INTO sessions (id, userId, expiry)
		VALUES (?, ?, ?)
	`
	_, err = tx.Exec(query, session.Id, session.UserId, (*sqlite.NullTime)(&session.Expiry))
	if err != nil {
		t.Fatal(err)
	}
}

func InitAuthService(t testing.TB) (goChat.AuthService, *sqlite.DB, func(), context.Context) {
	t.Helper()
	db := MustOpenDB(t)
	ctx := context.Background()
	s := sqlite.NewAuthService(db)
	return s, db, func() { MustCloseDB(t, db) }, ctx
}
