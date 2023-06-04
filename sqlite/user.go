package sqlite

import (
	"context"

	"github.com/adamni21/goChat"
	"github.com/adamni21/goChat/crypto"
)

const service = "sqlite.userService."

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
	const op = service + "Create"
	encryptedPw, err := s.pwHasher.Generate(password)
	if err != nil {
		return goChat.NewInternalErr("generating password hash", op, err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return goChat.NewInternalErr("beginning transaction", op, err)
	}
	defer tx.Rollback()

	err = createUser(ctx, tx, user, encryptedPw)
	if err != nil {
		return goChat.NewInternalErr("", op, err)
	}

	if err = tx.Commit(); err != nil {
		return goChat.NewInternalErr("committing transaction", op, err)
	}

	return nil
}

func createUser(ctx context.Context, tx *Tx, user *goChat.User, encryptedPw string) error {
	const op = service + "create"

	user.CreatedAt = tx.now
	user.UpdatedAt = tx.now

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
		return goChat.NewInternalErr("inserting into users table", op, err)
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		return goChat.NewInternalErr("getting last inserted id", op, err)
	}

	return nil
}
