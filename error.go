package goChat

import (
	"fmt"
)

type ErrCode = uint16

const (
	EInternal ErrCode = 1
	ENotFound ErrCode = 2
)

type Error struct {
	Code    ErrCode
	Message string

	Op  string
	Err error
}

func (e Error) Error() string {
	return fmt.Sprintf("goChat error: code=%d message=%s", e.Code, e.Message)
}

func NewInternalErr(message, op string, err error) Error {
	return Error{Code: EInternal, Message: message, Op: op, Err: err}
}
