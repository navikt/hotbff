package texas

import "context"

type userKeyType int

const userKey userKeyType = 0

// User represents an authenticated user.
type User struct {
	Authenticated bool   // indicates whether the user is authenticated
	Token         string // the incoming bearer token from the Authorization header
}

// NewContext returns a new context with the given [User].
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext retrieves the [User] from the context. If no user is found, it returns a [User] with Authenticated set to false.
func FromContext(ctx context.Context) *User {
	if u, ok := ctx.Value(userKey).(*User); ok {
		return u
	}
	return &User{Authenticated: false}
}
