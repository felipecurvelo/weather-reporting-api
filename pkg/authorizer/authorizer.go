package authorizer

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"
)

type Authorizer interface {
	GenerateAccessToken() string
}

type MainAuth struct {
	validToken string
}

func (auth MainAuth) GenerateAccessToken() string {
	auth.validToken = auth.createHash(time.Now().String())
	return auth.validToken
}

func (auth *MainAuth) createHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type contextKey struct{}

func NewAuth() MainAuth {
	return MainAuth{}
}

func FromContext(ctx context.Context) Authorizer {
	auth, _ := ctx.Value(contextKey{}).(Authorizer)
	return auth
}

func NewContext(parentContext context.Context, auth Authorizer) context.Context {
	return context.WithValue(parentContext, contextKey{}, auth)
}
