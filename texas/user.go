package texas

import "context"

type userKeyType int

const userKey userKeyType = 0

// User represents an authenticated user.
type User struct {
	// Authenticated indicates whether the user is authenticated.
	Authenticated bool
	// Token is the incoming bearer token from the Authorization header.
	Token string
}

// NewContext returns a new context with the given user.
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext retrieves the user from the context. If no user is found, it returns a User with Authenticated set to false.
func FromContext(ctx context.Context) *User {
	if u, ok := ctx.Value(userKey).(*User); ok {
		return u
	}
	return &User{Authenticated: false}
}
