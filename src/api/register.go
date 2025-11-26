package api

import (
	"fastlink/src/db"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		c.JSON(400, resp.Error(400, "Invalid request"))
		return
	}
	// 用户名和密码验证
	if !isValidUsername(req.Username) || !isValidPassword(req.Password) {
		c.JSON(400, resp.Error(400, "Invalid username or password"))
		return
	}

	// 创建用户
	NewUser.Username = req.Username
	NewUser.PasswordHash = hash(req.Password)
	NewUser.AccessTokenID = randStr()

	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// 使用OnConflict避免竞态条件
		err := db.MySQLClient.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&NewUser).Error

		if err != nil {
			return err
		}
		// 写后检查
		if NewUser.ID == 0 {
			c.JSON(409, resp.Error(409, "Username already exists"))
			return gorm.ErrRegistered
		}

		return nil

	})

	if err != nil {
		if err == gorm.ErrRegistered {
			c.JSON(409, resp.Error(409, err.Error()))
		}
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	c.JSON(201, resp.OK(201, RegisterResponse{}))
}
