package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Token struct {
	jwt.RegisteredClaims
	Type     string `json:"type"`
	UserID   string `json:"userId"`
	Username string `json:"username"`
}
