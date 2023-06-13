package goChat

import (
	"fmt"
)

type ErrCode = uint16

const (
	EInternal     ErrCode = 1
	ENotFound     ErrCode = 2
	EUnauthorized ErrCode = 3
)

type Error struct {
	// Machine readable error code.
	Code ErrCode
	// Info for the operator. For example state.
	Info string

	// Func in which the error occured.
	Op string
	// Nested error.
	Err error

	// Error message for end user.
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("goChat error: code=%d message=%s", e.Code, e.Message)
}

func (e Error) ErrCode() ErrCode {
	if e.Code != 0 {
		return e.Code
	}
	if err, ok := e.Err.(Error); ok {
		return err.ErrCode()
	}
	return EInternal
}

func (e Error) ErrMessage() string {
	if e.Message != "" {
		return e.Message
	}
	if err, ok := e.Err.(Error); ok {
		err.ErrMessage()
	}
	return "Internal error."
}

func NewInternalErr(info, op, message string, err error) Error {
	return Error{Code: EInternal, Info: info, Op: op, Err: err, Message: message}
}
