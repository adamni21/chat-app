package sqlite

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

const authServiceOp = "sqlite.AuthService."

// Authservice represents a service for authentication
type AuthService struct {
	db       *DB
	pwHasher crypto.PasswordHasher
}

// returns new instance of AuthService
func NewAuthService(db *DB) *AuthService {
	return &AuthService{
		db:       db,
		pwHasher: crypto.NewArgon2Hasher(),
	}
}

// Verifies user and creates new session for specified user.
// Returns id of created session.
//
// Returns ENotFound if user doesn't exist.
// Returns EUnauthorized if credentials are invalid.
func (s *AuthService) Login(ctx context.Context, user goChat.User, password string) (goChat.Session, error) {
	const op = authServiceOp + "Login"
	correct, err := s.VerifyUser(ctx, user, password)
	if err != nil {
		return goChat.Session{}, goChat.Error{Op: op, Err: err}
	}
	if !correct {
		return goChat.Session{}, goChat.NewUnauthorizedErr("", op, "Wrong password.", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return goChat.Session{}, goChat.Error{Op: op, Err: err}
	}
	defer tx.Rollback()

	session, err := createSession(ctx, tx, user.Id)
	if err != nil {
		return goChat.Session{}, goChat.Error{Op: op, Err: err}
	}

	if err = tx.Commit(); err != nil {
		return goChat.Session{}, goChat.NewInternalErr("committing sessionTx", op, "", err)
	}

	return session, nil
}

// Deletes specified session.
func (s *AuthService) DeleteSession(ctx context.Context, sessionId goChat.SessionId) error {
	const op = authServiceOp + "DeleteSession"
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return goChat.NewInternalErr("", op, "", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM sessions WHERE id = ?;", sessionId)
	if err != nil {
		return goChat.NewInternalErr("deleting session from DB", op, "", err)
	}
	err = tx.Commit()
	if err != nil {
		return goChat.NewInternalErr("committing session delete tx", op, "", err)
	}
	return nil
}

// Retrieves a single session by sessionId.
//
// Returns ENotFound error if specified sessionId doesn't exist.
func (s *AuthService) FindSession(ctx context.Context, sessionId goChat.SessionId) (*goChat.Session, error) {
	const op = authServiceOp + "FindSession"
	session := &goChat.Session{}

	query := `
		SELECT id, userId, expiry FROM sessions
		WHERE id = ?;
	`
	row := s.db.db.QueryRowContext(ctx, query, sessionId)
	err := row.Scan(&session.Id, &session.UserId, (*NullTime)(&session.Expiry))
	if err != nil {
		info := fmt.Sprintf("sessionId: %s", sessionId)
		if err == sql.ErrNoRows {
			return nil, goChat.NewNotFoundErr(info, op, "", nil)
		}

		return nil, goChat.NewInternalErr(info, op, "", err)
	}

	return session, nil
}

// Returns bool whether password is correct.
//
// Returns ENotFound if user doesn't exist.
func (s *AuthService) VerifyUser(ctx context.Context, user goChat.User, password string) (bool, error) {
	const op = authServiceOp + "VerifyUser"
	passwordDigest, err := s.getPasswordDigest(user)
	if err != nil {
		return false, goChat.Error{Op: op, Err: err}
	}

	isCorrect, err := s.pwHasher.Verify(password, passwordDigest)
	if err != nil {
		return false, goChat.Error{Op: op, Err: err}
	}
	return isCorrect, nil
}

func createSession(ctx context.Context, tx *Tx, userId goChat.Id) (goChat.Session, error) {
	const op = "createSession"
	sessionId, err := crypto.GenerateRandomBytes(16)
	if err != nil {
		return goChat.Session{}, goChat.NewInternalErr("generate random bytes", op, "", err)
	}

	b64SessionId := base64.URLEncoding.EncodeToString(sessionId)
	expiry := tx.now.Add(30 * 24 * time.Hour).Truncate(time.Second)
	query := `
		INSERT INTO sessions (id, userId, expiry)
		VALUES (?, ?, ?)
	`
	_, err = tx.Exec(query, b64SessionId, userId, (*NullTime)(&expiry))
	if err != nil {
		return goChat.Session{}, goChat.NewInternalErr("inserting into sessions table", op, "", err)
	}

	return goChat.Session{
		Id:     goChat.SessionId(b64SessionId),
		UserId: userId,
		Expiry: expiry,
	}, nil
}

// Retrieves password digest from DB for specified user.
//
// Returns ENotFound if user doesn't exist.
func (s *AuthService) getPasswordDigest(user goChat.User) (string, error) {
	const op = "getPasswordDigest"
	var passwordDigest string
	query := `
		SELECT passwordString FROM users
		WHERE id = ?;
	`
	err := s.db.db.QueryRow(query, user.Id).Scan(&passwordDigest)
	if err == sql.ErrNoRows {
		info := fmt.Sprintf("user not found, this indicates a bug, user: %+v", user)
		return "", goChat.NewNotFoundErr(info, op, "", err)
	}
	return passwordDigest, nil
}
