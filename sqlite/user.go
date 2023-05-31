package sqlite

import (
	"context"
	"fmt"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

type userService struct {
	db       *DB
	pwHasher crypto.PasswordHasher
}

func NewUserService(db *DB) *userService {
	return &userService{
		db:       db,
		pwHasher: crypto.NewArgon2Hasher(),
	}
}

func (s *userService) Create(ctx context.Context, user *goChat.User, password string) error {
	encryptedPw, err := s.pwHasher.Generate(password)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createUser(ctx, tx, user, encryptedPw)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func createUser(ctx context.Context, tx *Tx, user *goChat.User, encryptedPw string) error {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			username,
			email,
    		isVerified,
    		passwordString
		)
		VALUES (?, ?, ?, ?)
	`, user.Username, user.Email, false, encryptedPw)
	if err != nil {
		return fmt.Errorf("insert into users: %w", err)
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting assigned id: %w", err)
	}

	user.CreatedAt = tx.now
	user.UpdatedAt = tx.now

	return nil
}
