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

const (
	batchSize = 150
)

func (service *WorkerService) ComputeGPA(stream pb.WorkerService_ComputeGPAServer) error {

	log.Println("Starting to compute GPA on worker service")
	students := make([]*pb.StudentWithGPA, batchSize)
	i := 0
	for {

		computeGPAReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		for _, student := range computeGPAReq.StudentsWithGrades {
			studentWithGPA, err := computeGPA(student)
			if err != nil {
				log.Printf("Error while computing GPA for student %s: %v", student.StudentName, err)
				return err
			}

			students[i] = studentWithGPA
			i++
			if i >= batchSize {
				i = 0
				computeGPAResp := &pb.ComputeGPAResponse{
					StudentsWithGPA: students,
				}
				if err := stream.Send(computeGPAResp); err != nil {
					return err
				}
			}

		}

		if err != nil {
			return err
		}
	}

}
