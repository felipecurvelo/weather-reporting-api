package auth

import "context"

type Auth struct {
	Name string
}

type contextKey struct{}

func FromContext(ctx context.Context) Auth {
	auth, _ := ctx.Value(contextKey{}).(Auth)
	return auth
}

func NewContext(parentContext context.Context, auth Auth) context.Context {
	return context.WithValue(parentContext, contextKey{}, auth)
}
