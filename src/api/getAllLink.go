package api

import (
	"fastlink/src/db"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllLink(c *gin.Context) {
	//admin only
	var err error
	var records []db.Link

	records, err = gorm.G[db.Link](db.MySQLClient).Order("id ASC").Find(db.Ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, resp.Error(500, "internal server error"))
		return
	}
	c.JSON(200, resp.OK(200, records))

}
