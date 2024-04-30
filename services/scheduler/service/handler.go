package main

import (
	"encoding/csv"
	"log"
	"mime/multipart"
	"net/http"
	"time"

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

func handleCSVRequest(c *gin.Context) {

	now := time.Now()
	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
	}

	openedFile, _ := form.File.Open()
	defer openedFile.Close()

	reader := csv.NewReader(openedFile)
	records, err := reader.ReadAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file"})
	}

	log.Println(len(records))
	c.JSON(http.StatusOK, gin.H{
		"message":  "Uploaded successfully",
		"duration": time.Since(now).String(),
	})
}
