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
	ok, err := auth.AuthRefreshToken(c)
	if err != nil {
		c.JSON(500, resp.Error(500, "Internal server error"))
	}
	if !ok {
		c.JSON(401, resp.Error(401, "Unauthorized"))
	}

	refreshToken, _ := auth.ParseToken(c)
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
