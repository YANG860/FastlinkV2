package api

import (
	"fastlink/src/auth"

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

}
func redirectOneShot(c *gin.Context) {
	//todo
}
