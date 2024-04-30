package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"service-provider/pkg/types"

	"github.com/gin-gonic/gin"
)

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handleComputeGPAForStudents(c *gin.Context) {
	var students []types.StudentRequestRow
	err := c.ShouldBind(&students)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported request payload"})
		return
	}

	_, err = computeGPAForStudents(students)
	fmt.Println("func")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while computing GPA"})
		return
	}

	// b, err := json.Marshal(res)
	c.JSON(http.StatusOK, gin.H{
		"s": "studnet",
	})
}
