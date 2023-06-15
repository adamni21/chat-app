package goChat

import "context"

// Array of bytes encoded to base64 string.
type SessionId string

// Represents a session.
type Session struct {
	SessionId SessionId
	UserId    Id
}

type AuthService interface {
	// Verifies user and creates new session for specified user.
	// Returns id of created session.
	//
	// Returns ENotFound if user doesn't exist.
	// Returns EUnauthorized if credentials are invalid.
	// Login(ctx context.Context, user User, password string) (SessionId, error)

	// Deletes specified session.
	//
	// 	Returns ENotFound if session doesn't exist.
	// DeleteSession(ctx context.Context, sessionId SessionId) error

	// Retrieves userId for given sessionId.
	//
	// Returns ENotFound error if specified sessionId doesn't exist.
	// UserId(ctx context.Context, sessionId SessionId) (Id, error)

	// Returns true if credentials are valid.
	//
	// Returns ENotFound if user doesn't exist.
	// Can return EInternal.
	VerifyUser(ctx context.Context, user User, password string) (bool, error)
}
