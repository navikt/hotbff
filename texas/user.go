package texas

import "context"

type userKeyType int

const userKey userKeyType = 0

// User represents a user with authentication status.
type User struct {
	Authenticated bool   // Indicates whether the user is authenticated.
	Token         string // The incoming bearer token from the Authorization header, or empty if not authenticated.
}

// NewContext returns a new context with the given user.
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext retrieves the user from the context. If no user is found, it returns a non-nil user with Authenticated set to false.
func FromContext(ctx context.Context) *User {
	if u, ok := ctx.Value(userKey).(*User); ok {
		return u
	}
	return &User{Authenticated: false}
}
