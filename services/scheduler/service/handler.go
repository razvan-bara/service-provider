package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"service-provider/pkg/types"
	"time"

	"github.com/gin-gonic/gin"
)

// const (
// 	payloadRowLimit = 10
// )

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
		return
	}

	openedFile, _ := form.File.Open()
	defer openedFile.Close()

	reader := csv.NewReader(openedFile)

	_, err = reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file"})
		return
	}

	cnt := 0
	batch := make([]*types.StudentRequestRow, 0)
	for {
		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file"})
			return
		}
		batch = append(batch, &types.StudentRequestRow{
			Name: record[0],
			PAJ:  record[1],
			DA:   record[2],
			PP:   record[3],
			MDS:  record[4],
			SGSC: record[5],
			IBD:  record[6],
			BT:   record[7],
		})
		cnt++
	}
	log.Println("Count of records in the file: ", cnt)

	payload, err := json.Marshal(batch)
	if err != nil {
		log.Println("Error marshalling the payload")
	}
	req, err := http.NewRequest("POST", "http://localhost:8081/gpa", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error creating the request")

	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	resBody := response.Body

	c.JSON(http.StatusOK, gin.H{
		"message":  "Uploaded successfully",
		"duration": time.Since(now).String(),
		"body":     resBody,
	})
}
