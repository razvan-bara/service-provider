package main

import (
	"fmt"
	"log"
	pbWorker "service-provider/services/worker/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	scheduler := &Scheduler{}
	for i := 1; i <= 3; i++ {
		conn, err := grpc.Dial(
			fmt.Sprintf("localhost:808%d", i),
			grpc.WithInsecure(),
		)

		if err != nil {
			log.Fatalln("Failed to dial GRPC server, got err: ", err)
		}

		pbWorkerClinet := pbWorker.NewWorkerServiceClient(conn)
		scheduler.Conns = append(scheduler.Conns, conn)
		scheduler.Clients = append(scheduler.Clients, pbWorkerClinet)
	}

	defer scheduler.CloseConnections()

	handler := &SchedulerHandler{
		Scheduler: scheduler,
	}

	var router = gin.Default()

	router.GET("/ping", handler.getHealth)
	router.POST("/file", handler.handleCSVRequest)

	router.Run() // listen and serve on 0.0.0.0:8080
}
