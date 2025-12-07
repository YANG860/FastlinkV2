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

type GenlinkRequest struct {
	LinkType  string `json:"linkType" binding:"required"`
	SourceURL string `json:"sourceUrl" binding:"required,url"`
	ShortCode string `json:"shortCode" binding:"omitempty,alphanum"`
}

type GenlinkResponse struct {
	Url string `json:"shortCode"`
}

func Genlink(c *gin.Context) {

	var body GenlinkRequest
	var link db.Link
	var accessToken *auth.Token
	var code string

	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.JSON(400, resp.Error(400, "Invalid request body"))
		return
	}

	ok, err := auth.AuthAccessToken(c)
	if err != nil {
		slog.Error("Failed to authenticate access token", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if !ok {
		slog.Warn("Unauthorized access attempt")
		c.JSON(403, resp.Error(403, "Forbidden"))
	}

	accessToken, _ = auth.ParseToken(c)

	if body.LinkType == db.LinkTypeCustom {
		// 拦截自定义短链接请求
		exist, err := db.ShortCodeBloomFilterExists(body.ShortCode)
		if err != nil {
			slog.Error("Failed to check short code existence", "error", err)
			c.JSON(500, resp.Error(500, "Internal server error"))
			return
		}

		if exist {
			slog.Warn("Short code already exists", "shortCode", body.ShortCode)
			c.JSON(400, resp.Error(400, "Short code already exists"))
			return
		}

		code = body.ShortCode
	} else {
		// 一般短链接
		code = genShortCode(config.Server().ShortCodeLength)
	}

	userID, err := strconv.ParseUint(accessToken.UserID, 10, 32)
	if err != nil {
		slog.Error("Invalid user ID", "error", err)
		c.JSON(500, resp.Error(500, "Invalid user ID"))
		return
	}
	userID_Uint := uint(userID)
	// 检查短链接是否在黑名单中
	err = db.MySQLClient.Find(&db.UrlBlacklist{}, "url = ?", body.SourceURL).Error
	if err == nil {
		slog.Warn("Source URL is blacklisted", "sourceUrl", body.SourceURL)
		c.JSON(400, resp.Error(400, "The source URL is blacklisted"))
		return
	}

	link = db.Link{
		Type:      body.LinkType,
		CreatorID: userID_Uint,
		SourceURL: body.SourceURL,
		ShortCode: code,
	}
	// 保存到数据库
	if err := gorm.G[db.Link](db.MySQLClient).Create(db.Ctx, &link); err != nil {
		slog.Error("Failed to create link", "error", err)
		c.JSON(500, resp.Error(500, "Failed to create link"))
		return
	}

	url := config.Server().Domain + "/" + code

	
	c.JSON(200, resp.OK(200, GenlinkResponse{
		Url: url,
	}))

}

func genShortCode(length int) string {

	for {
		code := utils.RandStr(length)
		exist, err := db.ShortCodeBloomFilterExists(code)
		if err != nil {
			continue
		}
		if exist {
			continue
		}
		db.ShortCodeBloomFilterAdd(code)
		return code
	}

}
