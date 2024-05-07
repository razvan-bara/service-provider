package main

import (
	"fmt"
	"log"
	"net"
	"os"
	pbWorker "service-provider/services/worker/proto"

	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("PORT")
	fmt.Println("PORT FROM ENV ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pbWorker.RegisterWorkerServiceServer(grpcServer, NewWorkerService())
	log.Println("Worker service started on port ", port)
	grpcServer.Serve(lis)
	// router.Run("localhost:8081")
}
