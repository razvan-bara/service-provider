package main

import "github.com/gin-gonic/gin"

var router = gin.Default()

func init() {
	router.GET("/ping", getHealth)
	router.POST("/file", handleCSVRequest)
}
