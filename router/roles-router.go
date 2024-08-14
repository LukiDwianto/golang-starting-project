package router

import (
	"golang-starting-project/controllers"
	"golang-starting-project/middleware"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(router *gin.Engine) {
	roleGroup := router.Group("/role")
	roleGroup.GET("/:id", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"PAGE_ROLE"}), controllers.FindOneRole)
	roleGroup.POST("/get-role", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"PAGE_ROLE"}), controllers.GetRole)
	roleGroup.POST("", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"CRUD_ROLE"}), controllers.CreateRole)
	roleGroup.PUT("/update/:id", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"CRUD_ROLE"}), controllers.UpdateRole)
	roleGroup.PUT("/update-privileges/:id", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"CRUD_ROLE"}), controllers.UpdateRolePrivilages)
	roleGroup.DELETE("/delete/:id", middleware.AuthMiddleware(), middleware.CheckPrivilegesMiddleware([]string{"CRUD_ROLE"}), controllers.DeleteRole)
}
