package main

import (
	"fmt"
	"golang-starting-project/models"
	"golang-starting-project/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	fmt.Printf("Waiting for db service up.. \n")
	time.Sleep(5 * time.Second)
	models.InitDB()
	router.AuthRoutes(r)
	router.RoleRoutes(r)

	r.Run()
}
