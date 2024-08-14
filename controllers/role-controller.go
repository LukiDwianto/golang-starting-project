package controllers

import (
	"golang-starting-project/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRole(c *gin.Context) {
	var filterRole models.SearchRole

	if err := c.ShouldBindJSON(&filterRole); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roleList, totalSize, err := models.GetRolesPagination(filterRole)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roleList,
		"total_size": totalSize,
	})
}

func UpdateRole(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	var createdRole models.Role
	if err := c.ShouldBindJSON(&createdRole); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.UpdateRole(ID, createdRole)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sukses Ubah role"})
}

func FindOneRole(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	role, err := models.FindOneRole(ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, role)
}

func DeleteRole(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	err := models.DeleteRole(ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sukses hapus role"})
}

func UpdateRolePrivilages(c *gin.Context) {

	var ID int64

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); id != 0 && err == nil {
		ID = id
	}

	var privilegesIDS []int64

	if err := c.ShouldBindJSON(&privilegesIDS); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateRolePrivilages(ID, privilegesIDS)

	if err != nil {
		log.Println("error :", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role update successfully"})

}

func CreateRole(c *gin.Context) {

	type RequestData struct {
		Data         models.Role `json:"data"`
		PrivilegesID []uint      `json:"privileges_id"`
	}

	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.CreateRole(requestData.Data, requestData.PrivilegesID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role registered successfully"})
}
