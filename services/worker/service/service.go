package main

import (
	"context"
	pb "service-provider/services/worker/proto"
)

type WorkerService struct {
	pb.UnimplementedWorkerServiceServer
}

// func (service *WorkerService) mustEmbedUnimplementedWorkerServiceServer() {}

func NewWorkerService() *WorkerService {
	return &WorkerService{}
}

func (service *WorkerService) ComputeGPA(context.Context, *pb.ComputeGPARequest) (*pb.ComputeGPAResponse, error) {

	return &pb.ComputeGPAResponse{
		Test: "worker service",
	}, nil
}
