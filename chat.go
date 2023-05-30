package goChat

import "time"

type Chat struct {
	member1 User
	member2 User

	message message

	CreatedAt time.Time
	UpdatedAt time.Time
}

type message struct {
	content string

	CreatedAt time.Time
	UpdatedAt time.Time
}
