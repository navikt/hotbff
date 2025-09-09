package texas

import (
	"context"
)

type key int

const userKey = key(0)

type User struct {
	Authenticated bool
	Token         string
}

func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func FromContext(ctx context.Context) *User {
	if u, ok := ctx.Value(userKey).(*User); ok {
		return u
	} else {
		return &User{Authenticated: false}
	}
}
