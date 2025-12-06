package api

import (
	"fastlink/src/auth"
	"fastlink/src/db"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetLinkByUserRequest struct {
}

type GetLinkByUserResponse struct {
	Links []db.Link `json:"links"`
}

func GetLinkByUser(c *gin.Context) {
	

	var token *auth.Token
	var records []db.Link

	ok, err := auth.AuthAccessToken(c)
	if err != nil {
		c.JSON(500, resp.Error(500, "internal server error"))
		return
	}
	if !ok {
		c.JSON(401, resp.Error(401, "Unauthorized"))
		return	
	}
	
	token, _ = auth.ParseToken(c)
	records, err = gorm.G[db.Link](db.MySQLClient).Where("user_id = ?", token.UserID).Find(db.Ctx)
	if err != nil {
		c.JSON(500, resp.Error(500, "internal server error"))
		return
	}

	c.JSON(200, resp.OK(200, GetLinkByUserResponse{
		Links: records,
	}))

}
