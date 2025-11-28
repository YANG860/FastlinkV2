package auth

import (
	"errors"
	"fastlink/src/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(c *gin.Context) (*Token, error) {

	var err error
	var token *Token
	// 解析Token
	authHeader := c.GetHeader("Authorization")

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		err = errors.New("invalid Authorization header")
		return token, err
	}

	rawToken := authHeader[len(bearerPrefix):]

	// 基本验证，包括过期验证，合法性验证
	jwtToken, err := jwt.ParseWithClaims(rawToken, &Token{}, func(t *jwt.Token) (any, error) { return []byte(config.Jwt().JwtKey), nil },
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return token, err
	}

	token, _ = jwtToken.Claims.(*Token)
	return token, nil
}
