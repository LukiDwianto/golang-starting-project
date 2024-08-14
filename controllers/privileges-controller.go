package controllers

import (
	"golang-starting-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPrivilege(c *gin.Context) {
	privileges, err := models.GetPrivileges()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, privileges)
}
