package api

import (
	"fastlink/src/auth"
	"fastlink/src/db"
	resp "fastlink/src/response"
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

	// TODO: 用户名的布隆过滤器

	var req LoginRequest
	var err error
	var user db.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, resp.Error(400, "Invalid request"))
		return
	}

	err = db.MySQLClient.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		c.JSON(401, resp.Error(401, "Invalid username or password"))
		return
	}
	// 验证密码
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(401, resp.Error(401, "Invalid username or password"))
		return
	}

	if user.Banned {
		c.JSON(403, resp.Error(403, "User is banned"))
		return
	}

	// 更新旧的Tokenid，使用事务
	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {

		newTokenID := randStr()
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
		c.JSON(500, resp.Error(500, err.Error()))
		return
	}

	refreshToken, err := auth.GenRefreshToken(&user)
	if err != nil {
		c.JSON(500, resp.Error(500, err.Error()))
		return
	}

	c.JSON(200, resp.OK(200, LoginResponse{
		RefreshToken: refreshToken,
	}))

}
