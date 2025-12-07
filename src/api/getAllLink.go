package api

import (
	"fastlink/src/auth"
	"fastlink/src/db"
	resp "fastlink/src/response"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type GetAllLinkRequest struct {
}

type GetAllLinkResponse struct {
	Links []db.Link `json:"links"`
}


func GetAllLink(c *gin.Context) {
	//admin only

	var err error
	var records []db.Link

	ok, err := auth.AuthAdmin(c)
	if err != nil {
		slog.Error("Failed to authenticate admin", "error", err)
		c.JSON(500, resp.Error(500, "Internal server error"))
	}
	if !ok {
		slog.Warn("Unauthorized admin access attempt")
		c.JSON(403, resp.Error(403, "Forbidden"))
	}

	records, err = gorm.G[db.Link](db.MySQLClient).Order("id ASC").Find(db.Ctx)
	if err != nil {
		slog.Error("Failed to retrieve all links", "error", err)
		c.JSON(500, resp.Error(500, "internal server error"))
		return
	}
	c.JSON(200, resp.OK(200, GetAllLinkResponse{Links: records}))

}
