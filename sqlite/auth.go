package sqlite

import (
	"context"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

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

func (s *AuthService) VerifyUser(ctx context.Context, user goChat.User, password string) (bool, error) {

}

func (s *AuthService) getPasswordHash(user goChat.User) (string, error) {
	var passwordHash string
	query := `
		SELECT passwordString FROM users
		WHERE id = ?;
	`
	err := s.db.db.QueryRow(query, user.Id).Scan(&passwordHash)
	return "", err
}
