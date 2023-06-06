package goChat

import (
	"context"
	"time"
)

type Id = int64

type User struct {
	Id       Id
	Username string
	Email    string
	Verified bool

	Chats []Chat

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserService interface {
	Create(ctx context.Context, user *User, password string) error
	FindById(ctx context.Context, id Id) (*User, error)
}
