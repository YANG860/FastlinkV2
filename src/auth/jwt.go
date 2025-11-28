package auth

import (
	"fastlink/src/config"
	"fastlink/src/db"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenRefreshToken(user *db.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Jwt().AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        user.AccessTokenID,
			Issuer:    "fastlink",
		},
		Type:     RefreshTokenType,
		UserID:   strconv.Itoa(int(user.ID)),
		Username: user.Username,
	})

	signedToken, err := token.SignedString([]byte(config.Jwt().JwtKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenAccessToken(token *Token) (string, error) {

	token.Type = AccessTokenType
	token.IssuedAt = jwt.NewNumericDate(time.Now())
	token.ExpiresAt = jwt.NewNumericDate(time.Now().Add(config.Jwt().RefreshTokenTTL))

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)

	signedToken, err := RefreshToken.SignedString([]byte(config.Jwt().JwtKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
