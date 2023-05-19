package main

import "time"

type User struct {
	Id       int64
	Username string
	Email    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserService interface {
	Create(username, email, password string) (User, error)
	Get(id int64) (User, error)
}
