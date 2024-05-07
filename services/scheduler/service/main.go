package main

import (
	"log"
	pbWorker "service-provider/services/worker/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to dial GRPC server, got err: ", err)
	}
	defer conn.Close()

	pbWorkerClinet := pbWorker.NewWorkerServiceClient(conn)

	schedulerHandler := &SchedulerHandler{
		workerService: pbWorkerClinet,
	}

	var router = gin.Default()

	router.GET("/ping", schedulerHandler.getHealth)
	router.POST("/file", schedulerHandler.handleCSVRequest)

	router.Run() // listen and serve on 0.0.0.0:8080
}
