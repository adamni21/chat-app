package goChat

import (
	"context"
	"time"
)

type User struct {
	Id       int64
	Username string
	Email    string
	Verified bool

	Chats []Chat

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserService interface {
	Create(ctx context.Context, username, email, password string) (*User, error)
	Get(ctx context.Context, id int64) (*User, error)
}
