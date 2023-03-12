package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetIndex
// @Summary 首页
// @Description 这是一个首页的API
// @Accept json
// @Tags 首页
// @Success 200 string json{"code","message"}
// @Router /hello [get]
func GetIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "welcome",
	})
}
