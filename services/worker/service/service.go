package main

import (
	"io"
	"log"
	pb "service-provider/services/worker/proto"
)

type WorkerService struct {
	pb.UnimplementedWorkerServiceServer
}

// func (service *WorkerService) mustEmbedUnimplementedWorkerServiceServer() {}

func NewWorkerService() *WorkerService {
	return &WorkerService{}
}

func (service *WorkerService) ComputeGPA(stream pb.WorkerService_ComputeGPAServer) error {

	log.Println("Starting to compute GPA on worker service")
	for {

		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ComputeGPAResponse{
				Test: "Hello from worker service!",
			})
		}
		if err != nil {
			return err
		}

	}
	return nil
}
