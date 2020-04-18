package authorizer

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"
)

type Authorizer interface {
	GenerateAccessToken() string
	ValidateToken(string) bool
}

type MainAuth struct {
	validToken string
}

func (auth *MainAuth) GenerateAccessToken() string {
	auth.validToken = auth.createHash(time.Now().String())
	return auth.validToken
}

func (auth *MainAuth) ValidateToken(token string) bool {
	if token == "" {
		return false
	}

	return token == auth.validToken
}

func (auth *MainAuth) createHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewAuth() *MainAuth {
	return &MainAuth{}
}

type contextKey struct{}

func FromContext(ctx context.Context) Authorizer {
	auth, _ := ctx.Value(contextKey{}).(Authorizer)
	return auth
}

func NewContext(parentContext context.Context, auth Authorizer) context.Context {
	return context.WithValue(parentContext, contextKey{}, auth)
}
