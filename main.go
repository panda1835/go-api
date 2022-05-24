package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panda1835/go-api/restapi"
)

func main() {
	router := gin.Default()
	router.POST("/upload", restapi.PostImage)
	router.Run("localhost:8080")
}
