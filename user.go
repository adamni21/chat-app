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
	Create(ctx context.Context, user *User, password string) error
	FindById(ctx context.Context, id int64) (*User, error)
}
