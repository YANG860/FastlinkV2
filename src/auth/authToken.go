package auth

import (
	"errors"
	"fastlink/src/db"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRefreshToken(c *gin.Context) (bool, error) {
	// 验证AccessToken

	var err error
	var token *Token
	token, err = ParseToken(c)
	if err != nil {
		return false, err
	}

	if token.Type != RefreshTokenType {
		return false, nil
	}

	// 从缓存中获取Token ID进行对比
	cachedID, err := db.FetchRefreshTokenID(token.UserID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	// 缓存未命中，查询数据库
	if errors.Is(err, redis.Nil) {
		valid, err := RefreshTokenValidInDB(token)
		if err != nil {
			return false, err
		}
		return valid, nil
	}

	// 对比Token ID
	if cachedID != token.ID {
		return false, nil
	}

	// 设置或延期缓存
	db.CacheRefreshTokenID(token.UserID, token.ID, true)
	db.UpdateRefreshTokenTTL(token.UserID)
	return true, nil
}

func AuthAccessToken(c *gin.Context) (bool, error) {

	token, err := ParseToken(c)
	if err != nil {
		return false, err
	}

	if token.Type != AccessTokenType {
		return false, err
	}
	return true, nil
}

func AuthAdmin(c *gin.Context) (bool, error) {

	var err error

	token, err := ParseToken(c)
	if err != nil {
		return false, err
	}

	var admin db.Admin
	err = db.MySQLClient.Where("user_id = ?", token.UserID).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
