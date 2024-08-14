package controllers

import (
	"golang-starting-project/middleware"
	"golang-starting-project/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	var filterUser models.SearchUser

	if err := c.ShouldBindJSON(&filterUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userList, totalSize, err := models.GetUsers(filterUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userList,
		"total_size": totalSize,
	})
}

func UpdateRoleUser(c *gin.Context) {

	var updateUser models.User

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateUser.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username and password can't empty"})
		return
	}

	err := models.UpdateRoleUser(updateUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sukses Ubah role"})
}

func UpdateUser(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	var createdUser models.User

	createdUser.CreatedBy = ""

	if err := c.ShouldBindJSON(&createdUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existUser, err := models.CheckUpdateUSerExist(createdUser.UserName, createdUser.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if existUser.ID != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username already exist"})
		return
	}

	_, err = models.UpdateUser(ID, createdUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}

func CreateUser(c *gin.Context) {
	var createdUser models.User

	createdUser.CreatedBy = ""

	if err := c.ShouldBindJSON(&createdUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if createdUser.UserName == "" || createdUser.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username and password can't empty"})
		return
	}

	existUser, err := models.GetUserByUsername(createdUser.UserName)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if existUser.ID != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username already exist"})
		return
	}

	hashedPassword, err := middleware.HashPassword(createdUser.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	createdUser.Password = hashedPassword

	_, err = models.CreateUser(createdUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func UpdatePasswordUser(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	var createdUser models.User

	createdUser.CreatedBy = ""

	if err := c.ShouldBindJSON(&createdUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if createdUser.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username and password can't empty"})
		return
	}

	hashedPassword, err := middleware.HashPassword(createdUser.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	createdUser.Password = hashedPassword

	_, err = models.UpdatePasswordUser(ID, createdUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}