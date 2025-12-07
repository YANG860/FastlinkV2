package api

import (
	"fastlink/src/auth"
	"fastlink/src/config"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"fastlink/src/utils"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	RefreshToken string `json:"refreshToken"`
}

func Login(c *gin.Context) {

	var req LoginRequest
	var err error
	var user db.User
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.JSON(400, resp.Error(400, "Invalid request"))
		return
	}
	exist, err := db.UsernameBloomFilterExists(req.Username)
	if err != nil {
		slog.Error("Failed to check username existence", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if !exist {
		slog.Warn("Username does not exist", "username", req.Username)
		c.JSON(401, resp.Error(401, "Invalid username or password"))
		return
	}

	err = db.MySQLClient.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		c.JSON(401, resp.Error(401, "Invalid username or password"))
		return
	}
	// 验证密码
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		slog.Warn("Invalid password attempt", "username", req.Username)
		c.JSON(401, resp.Error(401, "Invalid username or password"))
		return
	}

	if user.Banned {
		slog.Warn("Banned user login attempt", "username", req.Username)
		c.JSON(403, resp.Error(403, "User is banned"))
		return
	}

	// 更新旧的Tokenid，使用事务
	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {

		newTokenID := utils.RandStr(config.Jwt().RefreshTokenIDLength)
		// 数据库更新用户的AccessTokenID
		user.AccessTokenID = newTokenID
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// 更新缓存
		// 使用SET

		err := db.CacheRefreshTokenID(strconv.Itoa(int(user.ID)), newTokenID, false)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to update access token ID", "error", err)
		c.JSON(500, resp.Error(500, err.Error()))
		return
	}

	refreshToken, err := auth.GenRefreshToken(&user)
	if err != nil {
		slog.Error("Failed to generate refresh token", "error", err)
		c.JSON(500, resp.Error(500, err.Error()))
		return
	}

	c.JSON(200, resp.OK(200, LoginResponse{
		RefreshToken: refreshToken,
	}))

}
