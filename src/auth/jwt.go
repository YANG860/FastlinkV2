package auth

import (
	"errors"
	"fastlink/src/config"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ParseToken(c *gin.Context) {
	// 解析Token
	authHeader := c.GetHeader("Authorization")

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.AbortWithStatusJSON(400, resp.Error(400, "Invalid Authorization header"))
		return
	}

	rawToken := authHeader[len(bearerPrefix):]

	// 基本验证，包括过期验证，合法性验证
	claims, err := jwt.ParseWithClaims(rawToken, &Token{}, func(t *jwt.Token) (any, error) { return []byte(config.Jwt().JwtKey), nil },
		jwt.WithExpirationRequired(),
	)

	if errors.Is(err, jwt.ErrTokenExpired) || !claims.Valid {
		c.AbortWithStatusJSON(400, resp.Error(400, "Token is expired or invalid"))
		return
	}

	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	// 设置Token到上下文
	token := claims.Claims.(*Token)
	c.Set("token", token)
	c.Next()
}

func AuthRefreshToken(c *gin.Context) {
	// 验证AccessToken
	var err error
	tokenAny, _ := c.Get("token")

	token := tokenAny.(*Token)

	if token.Type != AccessTokenType {
		c.AbortWithStatusJSON(400, resp.Error(400, "Invalid  token type"))
		return
	}
	// 从缓存中获取Token ID进行对比
	cachedID, err := db.FetchRefreshTokenID(token.UserID)
	if err != nil && !errors.Is(err, redis.Nil) {
		c.AbortWithStatus(500)
		return
	}
	// 缓存未命中，查询数据库
	if errors.Is(err, redis.Nil) {
		valid, err := RefreshTokenValidInDB(token)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if !valid {
			c.AbortWithStatusJSON(401, resp.Error(401, "Invalid token"))
			return
		}
	}
	// 对比Token ID
	if cachedID != token.ID {
		c.AbortWithStatusJSON(401, resp.Error(401, "Invalid token"))
		return
	}

	// 设置或延期缓存
	db.CacheRefreshTokenID(token.UserID, token.ID, true)
	db.UpdateRefreshTokenTTL(token.UserID)
	c.Next()
}

func AuthAccessToken(c *gin.Context) {

	tokenAny, _ := c.Get("token")

	token := tokenAny.(*Token)

	if token.Type != RefreshTokenType {
		c.AbortWithStatusJSON(400, resp.Error(400, "Invalid token type"))
		return
	}

	c.Next()
}

func AuthAdmin(c *gin.Context) {

	tokenAny, _ := c.Get("token")

	token := tokenAny.(*Token)

	//TODO: 数据库查询验证身份

	var admin db.Admin
	err := db.MySQLClient.Where("user_id = ?", token.UserID).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(403, resp.Error(403, "Forbidden"))
		} else {
			c.AbortWithStatusJSON(500, resp.Error(500, "Internal server error"))
		}
		return
	}

	c.Set("token", token)
	c.Next()

}

func RefreshTokenValidInDB(claims *Token) (bool, error) {

	user, err := gorm.G[db.User](db.MySQLClient).Where("id = ?", claims.UserID).First(db.Ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) || user.AccessTokenID != claims.ID {
		return false, nil
	}

	return true, err
}

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
