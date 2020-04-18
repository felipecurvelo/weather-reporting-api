package authorizer

type AuthMock struct {
	validToken string
}

func (auth *AuthMock) GenerateAccessToken() string {
	return "M0CK3D_T0K3N"
}

func NewAuthMock() *AuthMock {
	return &AuthMock{}
}

func (auth *AuthMock) ValidateToken(token string) bool {
	return token == "M0CK3D_T0K3N"
}
