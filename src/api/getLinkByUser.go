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
	//todo

	var token *auth.Token
	var records []db.Link
	t, exists := c.Get("token")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	token = t.(*auth.Token)
	records, err := gorm.G[db.Link](db.MySQLClient).Where("user_id = ?", token.UserID).Find(db.Ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "internal server error"))
		return
	}

	c.JSON(200, resp.OK(200, GetLinkByUserResponse{
		Links: records,
	}))

}
