package api

import (
	"fastlink/src/db"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetAllUserRequest struct {
}

type GetAllUserResponse struct {
	Users []db.User `json:"users"`
}

func GetAllUser(c *gin.Context) {
	// Admin only
	var records []db.User
	var err error
	records, err = gorm.G[db.User](db.MySQLClient).Order("id ASC").Find(db.Ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "internal server error"))
		return
	}

	c.JSON(200, resp.OK(200, GetAllUserResponse{Users: records}))
}
