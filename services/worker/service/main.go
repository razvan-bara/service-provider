package main

import (
	"fmt"
	"log"
	"net"
	pbWorker "service-provider/services/worker/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbWorker.RegisterWorkerServiceServer(grpcServer, NewWorkerService())
	log.Println("Worker service started on port 8081")
	grpcServer.Serve(lis)
	// router.Run("localhost:8081")
}
