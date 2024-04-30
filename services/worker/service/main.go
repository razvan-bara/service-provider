package main

import (
	"encoding/csv"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/file", func(c *gin.Context) {

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
		c.JSON(http.StatusOK, gin.H{"message": "Uploaded successfully"})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
