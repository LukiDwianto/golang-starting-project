package router

import (
	"golang-starting-project/controllers"
	"golang-starting-project/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	authGroup.POST("/login", controllers.LoginAuth)
	authGroup.GET("/refresh-token", middleware.AuthMiddleware(), controllers.RefreshTokenAuth)
}
