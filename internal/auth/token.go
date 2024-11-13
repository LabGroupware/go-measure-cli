package auth

import (
	"net/http"
	"time"
)

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Expiry       time.Time
}

func NewAuthToken(accessToken, refreshToken, tokenType string, expiry time.Time) *AuthToken {
	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		Expiry:       expiry,
	}
}

func (t *AuthToken) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.TokenType+" "+t.AccessToken)
}
