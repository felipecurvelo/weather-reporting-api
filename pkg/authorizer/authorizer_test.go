package authorizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken_ReturnToken(t *testing.T) {
	a := NewAuth()
	token := a.GenerateAccessToken()
	assert.NotEmpty(t, token)
}

func TestGenerateAccessToken_CalledTwice_ReturnDifferentTokens(t *testing.T) {
	a := NewAuth()
	token1 := a.GenerateAccessToken()
	token2 := a.GenerateAccessToken()
	assert.NotEqual(t, token1, token2)
}
