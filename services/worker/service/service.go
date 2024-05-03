package main

import (
	"context"
	pb "service-provider/services/worker/proto"
)

type WorkerService struct {
	pb.UnimplementedWorkerServiceServer
}

func (service *WorkerService) DoWork(context.Context, *pb.WorkRequest) (*pb.WorkResponse, error) {

	return nil, nil
}

// func (service *WorkerService) mustEmbedUnimplementedWorkerServiceServer() {}

func NewWorkerService() *WorkerService {
	return &WorkerService{}
}
