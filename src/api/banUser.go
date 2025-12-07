package api

import (
	"fastlink/src/auth"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"log/slog"

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

	var err error
	var body BanUserRequest
	var user db.User
	var links []db.Link

	//admin only
	ok, err := auth.AuthAdmin(c)
	if err != nil {
		slog.Error("Failed to authenticate admin", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if !ok {
		slog.Warn("Unauthorized admin access attempt")
		c.JSON(403, resp.Error(403, "Forbidden"))
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.AbortWithStatusJSON(400, resp.Error(400, "Invalid request body"))
		return
	}

	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// 查找用户
		// 更新用户的Banned字段
		// 锁定用户和产生的链接
		// 失效所有链接缓存

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

		// 失效缓存
		for _, link := range links {
			link.Type = db.LinkTypePrivate
			if err := tx.Save(&link).Error; err != nil {
				return err
			}
			if err := db.DeleteLinkCache(link.ShortCode); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to ban user", "error", err)
		c.AbortWithStatusJSON(500, resp.Error(500, "Failed to ban user"))
		return
	}

	c.JSON(200, resp.OK(200, BanUserResponse{}))
}
