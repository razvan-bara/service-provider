package main

import "github.com/gin-gonic/gin"

var router = gin.Default()

func init() {
	router.GET("/health", getHealth)
	router.POST("/gpa", handleComputeGPAForStudents)
}
