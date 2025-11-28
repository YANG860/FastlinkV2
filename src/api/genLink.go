package api

import (
	"fastlink/src/auth"
	"fastlink/src/config"
	"fastlink/src/db"
	resp "fastlink/src/response"
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
	//TODO: implement me
	var body GenlinkRequest
	var link db.Link
	var accessToken *auth.Token

	ok, err := auth.AuthAccessToken(c)
	if err != nil {
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}
	if !ok {
		c.JSON(403, resp.Error(403, "Forbidden"))
	}

	accessToken, _ = auth.ParseToken(c)


	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, resp.Error(400, "Invalid request body"))
		return
	}
	// 拦截自定义短链接请求
	if body.LinkType == db.LinkTypeCustom {
		genCustomLink(c, body)
		return
	}
	
	// 短链接
	code := genShortCode(config.Server().ShortCodeLength)

	userID, err := strconv.ParseUint(accessToken.UserID, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "Invalid user ID"))
		return
	}
	userID_Uint := uint(userID)

	link = db.Link{
		Type:      body.LinkType,
		CreatorID: userID_Uint,
		SourceURL: body.SourceURL,
		ShortCode: code,
	}
	// 保存到数据库
	if err := gorm.G[db.Link](db.MySQLClient).Create(db.Ctx, &link); err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "Failed to create link"))
		return
	}

	url := config.Server().Domain + "/" + code

	c.JSON(200, resp.OK(200, GenlinkResponse{
		Url: url,
	}))

}

func genShortCode(length int) string {
	//TODO: check if shortcode exists
	// base 64 coded
	return ""
}

func genCustomLink(c *gin.Context, body GenlinkRequest) {
	//TODO: 
}
