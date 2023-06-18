package goChat

import (
	"context"
	"time"
)

// [16]byte array encoded to base64URL string.
type SessionId string

// Represents a session.
type Session struct {
	Id     SessionId
	UserId Id
	Expiry time.Time
}

type AuthService interface {
	// Verifies user and creates new session for specified user.
	// Returns id of created session.
	//
	// Returns ENotFound if user doesn't exist.
	// Returns EUnauthorized if credentials are invalid.
	Login(ctx context.Context, user User, password string) (Session, error)

	// Deletes specified session.
	//
	// 	Returns ENotFound if session doesn't exist.
	// DeleteSession(ctx context.Context, sessionId SessionId) error

	// Retrieves a single session by sessionId.
	//
	// Returns ENotFound error if specified sessionId doesn't exist.
	FindSession(ctx context.Context, sessionId SessionId) (*Session, error)

	// Returns true if credentials are valid.
	//
	// Returns ENotFound if user doesn't exist.
	// Can return EInternal.
	VerifyUser(ctx context.Context, user User, password string) (bool, error)
}
