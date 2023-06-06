package goChat

import "context"

type contextKey int

const userIdKey = contextKey(1)

func NewContextWithUserId(ctx context.Context, id Id) context.Context {
	return context.WithValue(ctx, userIdKey, id)
}

func UserIdFromContext(ctx context.Context) Id {
	id, _ := ctx.Value(userIdKey).(Id)
	return id
}
