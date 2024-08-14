package controllers

import (
	"golang-starting-project/middleware"
	"golang-starting-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshTokenAuth(c *gin.Context) {

	username := c.GetString("username")
	role := c.GetString("role")
	privileges := c.GetStringSlice("privileges")

	tokenString, refreshToken, err := middleware.GenerateToken(username, role, privileges, c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User refresh token successfully",
		"token":         tokenString,
		"refresh_token": refreshToken,
	})
}

func LoginAuth(c *gin.Context) {

	type LoginDataDto struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginData LoginDataDto

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username and password can't empty"})
		return
	}

	userData, err := models.GetUserByUsername(loginData.Username)

	if userData.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username not found"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	match := middleware.ComparePassword(loginData.Password, userData.Password)

	if !match {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	var privileges []string

	for _, item := range *userData.Role.Privileges {
		privileges = append(privileges, item.Name)
	}

	tokenString, refreshToken, err := middleware.GenerateToken(userData.UserName, userData.Role.Name, privileges, c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully",
		"token":         tokenString,
		"refresh_token": refreshToken,
	})
}
