package main

import (
	"encoding/csv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	pbWorker "service-provider/services/worker/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SchedulerHandler struct {
	Scheduler *Scheduler
	TaskDB    *TaskDB
}

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (service *SchedulerHandler) getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (service *SchedulerHandler) handleCSVRequest(c *gin.Context) {

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
	studentsWithGrades := make([]*pbWorker.StudentWithGrades, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file"})
			return
		}

		studentsWithGrades = append(studentsWithGrades, addStudentWithGrades(record))
		cnt++
	}
	log.Println("Count of records in the file: ", cnt)

	requestId := uuid.New().String()
	err = service.TaskDB.InsertRow(requestId, studentsWithGrades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error trying to insert row"})
		return
	}

	log.Println("Inserted a new task into the queue")
	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"requestId": requestId,
	})

	// scheduler := service.Scheduler
	// students, err := scheduler.roundRobinDistributeComputeGPAWork(studentsWithGrades)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// computeGPAClient, err := scheduler.Clients[0].ComputeGPA(context.Background())
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// waitc := make(chan struct{})
	// students := make([]*pbWorker.StudentWithGPA, 0)
	// go func() {
	// 	for {
	// 		gpaResp, err := computeGPAClient.Recv()
	// 		if err == io.EOF {
	// 			close(waitc)
	// 			return
	// 		}
	// 		if err != nil {
	// 			log.Fatalf("Got error while computing GPA: %v", err)
	// 		}
	// 		students = append(students, gpaResp.StudentsWithGPA...)
	// 	}
	// }()

	// j := 0
	// batchSize := 150
	// for i := 0; i < len(studentsWithGrades); i += batchSize {

	// 	j = i + batchSize
	// 	if j > len(studentsWithGrades) {
	// 		j = len(studentsWithGrades)
	// 	}

	// 	batch := studentsWithGrades[i:j]
	// 	computeGPAReq := &pbWorker.ComputeGPARequest{
	// 		StudentsWithGrades: batch,
	// 	}
	// 	if err := computeGPAClient.Send(computeGPAReq); err != nil {
	// 		log.Fatalf("Failed to send a batch of students to compute GPA: %v", err)
	// 	}

	// }
	// computeGPAClient.CloseSend()
	// <-waitc

	// c.JSON(http.StatusOK, gin.H{
	// 	"studentsWithGPA": students,
	// })
}

func addStudentWithGrades(record []string) *pbWorker.StudentWithGrades {
	return &pbWorker.StudentWithGrades{
		StudentName: record[0],
		Grades: []*pbWorker.Grade{
			{
				CourseId: 1,
				Score:    record[1],
			},
			{
				CourseId: 2,
				Score:    record[2],
			},
			{
				CourseId: 3,
				Score:    record[3],
			},
			{
				CourseId: 4,
				Score:    record[4],
			},
			{
				CourseId: 5,
				Score:    record[5],
			},
			{
				CourseId: 6,
				Score:    record[6],
			},
			{
				CourseId: 7,
				Score:    record[7],
			},
		},
	}
}
