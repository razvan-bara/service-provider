package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pbWorker "service-provider/services/worker/proto"
	"sync"

	"google.golang.org/grpc"
)

type Scheduler struct {
	Conns   []*grpc.ClientConn
	Clients []pbWorker.WorkerServiceClient
}

func (scheduler *Scheduler) distribuiteComputeGPAWork(studentsWithGrades []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {
	log.Println("Starting to distribuite work")
	log.Println("Number of available clients: ", len(scheduler.Clients))
	students, err := scheduler.distribuiteEqualAcrossNodes(studentsWithGrades)
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (scheduler *Scheduler) distribuiteEqualAcrossNodes(studentsWithGrades []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {
	log.Println("Distribuiting workload equally across nodes")
	wg := sync.WaitGroup{}
	studentChan := make(chan *pbWorker.StudentWithGPA, 10_000)
	errChan := make(chan error)
	done := make(chan struct{})

	students := []*pbWorker.StudentWithGPA{}
	partition := len(studentsWithGrades) / len(scheduler.Clients)

	for i := 0; i < len(scheduler.Clients); i++ {

		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()

			beg, end := workerIndex*partition, (workerIndex+1)*partition
			if end >= len(studentsWithGrades) {
				end = len(studentsWithGrades)
			}

			workload := studentsWithGrades[beg:end]
			err := scheduler.computeGPA(workerIndex, workload, studentChan)
			if err != nil {
				errChan <- err
			}
		}(i)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	defer close(studentChan)
	for {
		select {
		case err := <-errChan:
			return nil, err
		case s := <-studentChan:
			students = append(students, s)
		case <-done:
			return students, nil
		}
	}
}

func (scheduler *Scheduler) computeGPA(workerIndex int, workload []*pbWorker.StudentWithGrades, studentChan chan *pbWorker.StudentWithGPA) error {
	computeGPAClient, err := scheduler.Clients[workerIndex].ComputeGPA(context.Background())

	if err != nil {
		return fmt.Errorf("failed to open a stream to compute GPA: %v", err)
	}

	waitc := make(chan struct{})
	errChan := make(chan error)
	go func() {
		for {
			gpaResp, err := computeGPAClient.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				errChan <- fmt.Errorf("got error while computing GPA: %v", err)
			}

			for i := 0; i < len(gpaResp.StudentsWithGPA); i++ {
				studentChan <- gpaResp.StudentsWithGPA[i]
			}
		}
	}()

	j := 0
	batchSize := 1000
	for i := 0; i < len(workload); i += batchSize {
		j += batchSize
		if j > len(workload) {
			j = len(workload)
		}

		batch := workload[i:j]
		computeGPAReq := &pbWorker.ComputeGPARequest{
			StudentsWithGrades: batch,
		}

		if err := computeGPAClient.Send(computeGPAReq); err != nil {
			return fmt.Errorf("failed to send a batch of students to compute GPA: %v", err)
		}

	}
	computeGPAClient.CloseSend()

	select {
	case err := <-errChan:
		return err
	case <-waitc:
		break
	}

	return nil
}

func (scheduler *Scheduler) CloseConnections() {
	for _, conn := range scheduler.Conns {
		conn.Close()
	}
}
