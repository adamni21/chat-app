package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

const authServiceOp = "sqlite.AuthService."

type AuthService struct {
	db       *DB
	pwHasher crypto.PasswordHasher
}

func NewAuthService(db *DB) *AuthService {
	return &AuthService{
		db:       db,
		pwHasher: crypto.NewArgon2Hasher(),
	}
}

// Returns bool whether password is correct.
//
// Returns ENotFound if user doesn't exist.
// Can return EInternal.
func (s *AuthService) VerifyUser(ctx context.Context, user goChat.User, password string) (bool, error) {
	const op = authServiceOp + "VerifyUser"
	passwordHash, err := s.getPasswordHash(user)
	if err != nil {
		return false, goChat.Error{Op: op, Err: err}
	}

	isCorrect, err := s.pwHasher.Verify(password, passwordHash)
	if err != nil {
		return false, goChat.Error{Op: op, Err: err}
	}
	return isCorrect, nil
}

// Retrieves password from DB for specified user.
//
// Returns ENotFound if user doesn't exist.
func (s *AuthService) getPasswordHash(user goChat.User) (string, error) {
	const op = "getPasswordHash"
	var passwordHash string
	query := `
		SELECT passwordString FROM users
		WHERE id = ?;
	`
	err := s.db.db.QueryRow(query, user.Id).Scan(&passwordHash)
	if err == sql.ErrNoRows {
		info := fmt.Sprintf("user not found, this indicates a bug, user: %+v", user)
		return "", goChat.NewNotFoundErr(info, op, "", err)
	}
	return passwordHash, nil
}
