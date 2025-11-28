package api

import (
	"fastlink/src/auth"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
)

func Redirect(c *gin.Context) {
	//TODO:

}

func redirectCustom(c *gin.Context) {
	// currently same as general
	redirectGeneral(c)
}
func redirectGeneral(c *gin.Context) {
	//todo
}
func redirectPrivate(c *gin.Context) {
	//todo
	ok, err := auth.AuthAccessToken(c)
	if err != nil {
		c.JSON(500, resp.Error(500, "Internal server error"))
	}
	if !ok {
		c.JSON(403, resp.Error(403, "Forbidden"))
	}

}
func redirectOneShot(c *gin.Context) {
	//todo
}
