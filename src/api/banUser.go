package api

import (
	"fastlink/src/db"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BanUserRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

type BanUserResponse struct {
}

func BanUser(c *gin.Context) {
	//admin only

	var err error
	var body BanUserRequest
	var user db.User
	var links []db.Link

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, resp.Error(400, "Invalid request body"))
		return
	}

	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// 查找用户
		// 更新用户的Banned字段
		// 锁定用户和产生的链接

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, body.UserID).Error; err != nil {
			return err
		}
		user.Banned = true
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("creator_id = ?", body.UserID).Find(&links).Error; err != nil {
			return err
		}

		for _, link := range links {
			link.Type = db.LinkTypePrivate
			if err := tx.Save(&link).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "Failed to ban user"))
		return
	}

	c.JSON(200, BanUserResponse{})
}
