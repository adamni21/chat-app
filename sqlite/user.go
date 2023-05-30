package sqlite

import (
	"context"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

type userService struct {
	db       *DB
	pwHasher crypto.PasswordHasher
}

func NewUserService(db *DB) *userService {
	return &userService{
		db: db,
	}
}

func (s *userService) Create(ctx context.Context, username, email, password string) (*goChat.User, error) {
	encryptedPw, err := s.pwHasher.Generate(password)
	if err != nil {
		return nil, err
	}

	user := &goChat.User{
		Username: username,
		Email:    email,
		Verified: false,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = createUser(ctx, tx, user, encryptedPw)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func createUser(ctx context.Context, tx *Tx, user *goChat.User, encryptedPw string) error {

}
