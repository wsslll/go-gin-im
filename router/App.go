package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/hello", service.GetIndex)
	r.GET("/send", service.SendMsg)
	
	user := r.Group("/user")
	user.GET("/list", service.GetUserList) // 用户列表接口
	user.POST("/create", service.CreateUser)
	user.PUT("/update", service.UpdateUser)
	user.DELETE("/delete/:id", service.DeleteUser)
	user.POST("/login", service.Login)
	return r
}
