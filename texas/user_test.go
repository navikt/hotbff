package texas

import (
	"context"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestUser(t *testing.T) {
	ctx := context.Background()
	u := FromContext(ctx)
	assert.Equal(t, u.Authenticated, false)
	assert.Equal(t, u.Token, "")

	ctx = NewContext(ctx, &User{Authenticated: true, Token: "userToken"})
	u = FromContext(ctx)
	assert.Equal(t, u.Authenticated, true)
	assert.Equal(t, u.Token, "userToken")
}
