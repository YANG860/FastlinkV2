package api

import (
	"errors"
	"fastlink/src/config"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"fastlink/src/utils"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
}

func Register(c *gin.Context) {

	var err error
	var req RegisterRequest
	var NewUser db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.JSON(400, resp.Error(400, "Invalid request"))
		return
	}
	// 用户名和密码验证
	if !isValidUsername(req.Username) || !isValidPassword(req.Password) {
		slog.Warn("Invalid username or password", "username", req.Username)
		c.JSON(400, resp.Error(400, "Invalid username or password"))
		return
	}
	exist, err := db.UsernameBloomFilterExists(req.Username)
	if err != nil {
		slog.Error("Failed to check username existence in bloom filter", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if exist {
		slog.Warn("Username already exists", "username", req.Username)
		c.JSON(409, resp.Error(409, "Username already exists"))
		return
	}

	// 创建用户
	NewUser.Username = req.Username
	NewUser.PasswordHash = hash(req.Password)
	NewUser.AccessTokenID = utils.RandStr(config.Jwt().RefreshTokenIDLength)

	err = gorm.G[db.User](db.MySQLClient).Create(db.Ctx, &NewUser)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		slog.Warn("Username already exists", "username", req.Username)
		c.JSON(409, resp.Error(409, "Username already exists"))
		return
	}

	if err != nil {
		slog.Error("Failed to create new user", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	// 添加到布隆过滤器
	err = db.UsernameBloomFilterAdd(req.Username)

	c.JSON(201, resp.OK(201, RegisterResponse{}))
}
