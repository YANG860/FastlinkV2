package api

import (
	"fastlink/src/auth"
	resp "fastlink/src/response"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func Refresh(c *gin.Context) {
	auth.ParseToken(c)
	auth.AuthRefreshToken(c)

	token, exists := c.Get("token")
	if !exists {
		c.JSON(401, resp.Error(401, "Unauthorized"))
		return
	}

	refreshToken, _ := token.(*auth.Token)
	// 通过验证即可刷新Token
	accessToken, err := auth.GenAccessToken(refreshToken)
	if err != nil {
		c.JSON(500, resp.Error(500, "Internal server error"))
		return
	}

	c.JSON(200, resp.OK(200, RefreshResponse{
		AccessToken: accessToken,
	}))
}
