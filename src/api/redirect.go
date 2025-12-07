package api

import (
	"fastlink/src/auth"
	"fastlink/src/config"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"fastlink/src/utils"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var cacheSyncManager = utils.NewTaskManager[string]()

func Redirect(c *gin.Context) {
	
	var link db.Link
	var err error

	shortCode := c.Param("shortCode")

	link, err = db.FetchLink(shortCode)

	switch err {
	case nil:
		//nil
	case redis.Nil:

		link, err = gorm.G[db.Link](db.MySQLClient).Where("short_code = ?", shortCode).First(db.Ctx)

		switch err {
		case nil:
			db.CacheLink(link, true)
		case gorm.ErrRecordNotFound:
			slog.Warn("Link not found", "shortCode", shortCode)
			c.JSON(404, resp.Error(404, "Link not found"))
			return
		default:
			slog.Error("Failed to retrieve link from database", "error", err)
			c.JSON(500, resp.Error(500, "Internal server error"))
			return
		}

	default:
		slog.Error("Failed to retrieve link from cache", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	db.UpdateLinkTTL(link.ShortCode)

	switch link.Type {
	case db.LinkTypeCustom:
		redirectCustom(c, link)
	case db.LinkTypeGeneral:
		redirectGeneral(c, link)
	case db.LinkTypePrivate:
		redirectPrivate(c, link)
	case db.LinkTypeOneShot:
		redirectOneShot(c, link)
	default:
		slog.Error("Invalid link type", "type", link.Type)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

}

func redirectCustom(c *gin.Context, link db.Link) {
	// currently same as general
	redirectGeneral(c, link)
}

func redirectPrivate(c *gin.Context, link db.Link) {
	//todo
	ok, err := auth.AuthAccessToken(c)
	if err != nil {
		slog.Error("Failed to authenticate access token", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if !ok {
		slog.Warn("Unauthorized access attempt")
		c.JSON(403, resp.Error(403, "Forbidden"))
		return
	}
	redirectGeneral(c, link)
}

func redirectGeneral(c *gin.Context, link db.Link) {

	// 更新访问计数

	err := db.UpdateLinkClicks(link.ShortCode)
	if err != nil {
		slog.Error("Failed to update link clicks", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	// 单例定时任务异步处理，定时批量更新
	cacheSyncManager.NewTask(config.Redis().ClickWriteBackInterval, link.ShortCode, func(shortCode string) {
		// 从缓存读取点击数
		cachedLink, err := db.FetchLink(shortCode)
		if err != nil {
			return
		}
		// 更新数据库
		_, err = gorm.G[db.Link](db.MySQLClient).Where("short_code = ?", shortCode).Update(db.Ctx, "clicks", cachedLink.Clicks)
		if err != nil {
			slog.Error("Failed to update link clicks in database", "error", err)
			return
		}
	})

	// 重定向
	c.Redirect(302, link.SourceURL)
}

func redirectOneShot(c *gin.Context, link db.Link) {
	// 检查是否已被访问
	var err error
	var rowsAffected int64

	err = db.MySQLClient.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&db.Link{}).Where("id = ? AND access_count = 0", link.ID).Update("access_count", gorm.Expr("access_count + ?", 1))
		err = result.Error
		rowsAffected = result.RowsAffected
		return err
	})
	if err != nil {
		slog.Error("Failed to update one-shot link access count", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	if rowsAffected == 0 {
		// 已经被访问过
		slog.Warn("One-shot link already accessed", "shortCode", link.ShortCode)
		c.JSON(404, resp.Error(404, "Link not found"))
		return
	}

	redirectGeneral(c, link)
}
