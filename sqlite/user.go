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

func (s *userService) FindById(ctx context.Context, id int64) (*goChat.User, error) {
	const op = service + "FindById"

	query := `
		SELECT id, username, email, isVerified, createdAt, updatedAt
		FROM users
		WHERE id = ?

	`
	rows, err := s.db.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, goChat.NewInternalErr("querying user", op, err)
	}
	defer rows.Close()

	rows.Next()
	user := &goChat.User{}
	err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Verified, (*NullTime)(&user.CreatedAt), (*NullTime)(&user.UpdatedAt))
	if err != nil {
		return nil, goChat.NewInternalErr("scanning row", op, err)
	}

	return user, nil
}

func createUser(ctx context.Context, tx *Tx, user *goChat.User, encryptedPw string) error {
	const op = service + "create"

	user.CreatedAt = tx.now
	user.UpdatedAt = tx.now

	query := `
		INSERT INTO users (
			username,
			email,
    		isVerified,
    		passwordString,
			createdAt,
			updatedAt
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		false,
		encryptedPw,
		(*NullTime)(&user.CreatedAt),
		(*NullTime)(&user.UpdatedAt),
	)
	if err != nil {
		return goChat.NewInternalErr("inserting into users table", op, err)
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		return goChat.NewInternalErr("getting last inserted id", op, err)
	}

	return nil
}
