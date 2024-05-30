package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	pbWorker "service-provider/services/worker/proto"
	"time"

	"google.golang.org/grpc"
)

type Scheduler struct {
	ClientConnections map[string]*ClientConnection
	TaskDB            *TaskDB
}

type ClientConnection struct {
	Conn         *grpc.ClientConn
	Client       pbWorker.WorkerServiceClient
	Load         int
	HasHeartbeat bool
}

var (
	index int
)

func (scheduler *Scheduler) ListenForTasks() {
	for {
		log.Println("Listening for tasks...")

		// Iterate through the tasks from the queue
		taskChan := make(chan *Task)
		go func() {
			_, err := scheduler.TaskDB.SelectPendingTasks(taskChan)
			if err != nil {
				log.Panicf("failed when getting tasks, err: %v\n", err)
			}
		}()

		for task := range taskChan {
			go func() {
				var studentsWithGrades []*pbWorker.StudentWithGrades
				err := json.Unmarshal(task.Payload, &studentsWithGrades)
				if err != nil {
					log.Fatal("failed to unmarshal task payload: ", err)
				}
				// Compute the GPA for the students
				studentWithGPA, err := scheduler.distributeRoundRobin(studentsWithGrades)
				if err != nil {
					log.Println("Couldn't compute GPA for students: ", err)
				} else {
					err = scheduler.TaskDB.UpdateRow(task, studentWithGPA)
					if err != nil {
						log.Println("Couldn't update the task status: ", err)
					}
				}
			}()
		}
		time.Sleep(1 * time.Second)
	}
}

func (scheduler *Scheduler) distributeRoundRobin(studentsWithGrades []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {
	log.Println("Distributing workload in round robin fashion")

	workerClient, cc, err := scheduler.decideClient()
	if err != nil {
		return nil, err
	}

	load := len(studentsWithGrades) / 100
	cc.Load += load
	log.Println("Load: ", cc.Load)
	stdudentsWithGPA, err := scheduler.computeGPARoundRobin(workerClient, studentsWithGrades)
	if err != nil {
		cc.Load -= load
		return nil, err
	}
	cc.Load -= load

	return stdudentsWithGPA, nil
}

func (scheduler *Scheduler) decideClient() (pbWorker.WorkerServiceClient, *ClientConnection, error) {
	// Implement the logic to decide which client to us
	minLoad, chosenPort := 9999, ""
	for port, tmpCC := range scheduler.ClientConnections {
		if !tmpCC.HasHeartbeat {
			continue
		}

		if isSmaller(tmpCC.Load, minLoad) {
			minLoad = tmpCC.Load
			chosenPort = port
		}

	}

	if chosenPort == "" {
		return nil, nil, fmt.Errorf("no worker client available")
	}

	log.Println("Chosen port: ", chosenPort)
	cc := scheduler.ClientConnections[chosenPort]
	return cc.Client, cc, nil
}

func isSmaller(x, y int) bool {
	if x < y {
		return true
	}

	return false
}

func (scheduler *Scheduler) computeGPARoundRobin(workerClient pbWorker.WorkerServiceClient, workload []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {

	computeGPAStreamClient, err := workerClient.ComputeGPA(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to open a stream to compute GPA: %v", err)
	}

	waitc := make(chan struct{})
	errChan := make(chan error)
	studentsWithGPA := make([]*pbWorker.StudentWithGPA, 0)
	go func() {
		for {
			gpaResp, err := computeGPAStreamClient.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				errChan <- fmt.Errorf("got error while computing GPA: %v", err)
			}

			for i := 0; i < len(gpaResp.StudentsWithGPA); i++ {
				studentsWithGPA = append(studentsWithGPA, gpaResp.StudentsWithGPA[i])
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

		if err := computeGPAStreamClient.Send(computeGPAReq); err != nil {
			return nil, fmt.Errorf("failed to send a batch of students to compute GPA: %v", err)
		}

	}
	computeGPAStreamClient.CloseSend()

	select {
	case err := <-errChan:
		return nil, err
	case <-waitc:
		return studentsWithGPA, nil
	}
}

func (scheduler *Scheduler) CheckWorkerHeartbeats() {
	for {
		for i := 1; i <= 3; i++ {
			port := fmt.Sprintf("808%d", i)
			conn, err := grpc.Dial(
				fmt.Sprintf("localhost:%s", port),
				grpc.WithInsecure(),
			)

			if err != nil {
				log.Println("Failed to get heartbeat for the worker client on port: ", port)
				continue
			}

			pbWorkerClinet := pbWorker.NewWorkerServiceClient(conn)
			if pbWorkerClinet == nil {
				log.Println("Failed to get heartbeat from worker client on port: ", port)
				continue
			}

			_, err = pbWorkerClinet.GetStatus(context.Background(), &pbWorker.GetStatusRequest{})
			heartbeat := scheduler.ClientConnections[port].HasHeartbeat
			if err != nil && heartbeat {
				log.Println("Failed to get heartbeat from worker client on port: ", port)
				heartbeat = false

			} else if err == nil && !heartbeat {
				log.Println("Got heartbeat for worker client on port: ", port)
				heartbeat = true
			}

			scheduler.ClientConnections[port].HasHeartbeat = heartbeat
		}
		time.Sleep(700 * time.Millisecond)
	}

}

func (scheduler *Scheduler) CloseConnections() {
	for _, cc := range scheduler.ClientConnections {
		cc.Conn.Close()
	}
}

// func (scheduler *Scheduler) distribuiteComputeGPAWork(studentsWithGrades []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {
// 	log.Println("Starting to distribuite work")
// 	log.Println("Number of available clients: ", len(scheduler.ClientConnections))
// 	students, err := scheduler.distribuiteEqualAcrossNodes(studentsWithGrades)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return students, nil
// }

// func (scheduler *Scheduler) distribuiteEqualAcrossNodes(studentsWithGrades []*pbWorker.StudentWithGrades) ([]*pbWorker.StudentWithGPA, error) {
// 	log.Println("Distribuiting workload equally across nodes")
// 	wg := sync.WaitGroup{}
// 	studentChan := make(chan *pbWorker.StudentWithGPA, 10_000)
// 	errChan := make(chan error)
// 	done := make(chan struct{})

// 	students := []*pbWorker.StudentWithGPA{}
// 	partition := len(studentsWithGrades) / len(scheduler.ClientConnections)

// 	for i := 0; i < len(scheduler.ClientConnections); i++ {

// 		wg.Add(1)
// 		go func(workerIndex int) {
// 			defer wg.Done()

// 			beg, end := workerIndex*partition, (workerIndex+1)*partition
// 			if end >= len(studentsWithGrades) {
// 				end = len(studentsWithGrades)
// 			}

// 			workload := studentsWithGrades[beg:end]
// 			err := scheduler.computeGPADistribuiteEqualAcrossNodes(workerIndex, workload, studentChan)
// 			if err != nil {
// 				errChan <- err
// 			}
// 		}(i)
// 	}

// 	go func() {
// 		wg.Wait()
// 		done <- struct{}{}
// 	}()

// 	defer close(studentChan)
// 	for {
// 		select {
// 		case err := <-errChan:
// 			return nil, err
// 		case s := <-studentChan:
// 			students = append(students, s)
// 		case <-done:
// 			return students, nil
// 		}
// 	}
// }

// func (scheduler *Scheduler) computeGPADistribuiteEqualAcrossNodes(workerIndex int, workload []*pbWorker.StudentWithGrades, studentChan chan *pbWorker.StudentWithGPA) error {
// 	computeGPAStreamClient, err := scheduler.ClientConnections[workerIndex].Client.ComputeGPA(context.Background())

// 	if err != nil {
// 		return fmt.Errorf("failed to open a stream to compute GPA: %v", err)
// 	}

// 	waitc := make(chan struct{})
// 	errChan := make(chan error)
// 	go func() {
// 		for {
// 			gpaResp, err := computeGPAStreamClient.Recv()
// 			if err == io.EOF {
// 				close(waitc)
// 				return
// 			}

// 			if err != nil {
// 				errChan <- fmt.Errorf("got error while computing GPA: %v", err)
// 			}

// 			for i := 0; i < len(gpaResp.StudentsWithGPA); i++ {
// 				studentChan <- gpaResp.StudentsWithGPA[i]
// 			}
// 		}
// 	}()

// 	j := 0
// 	batchSize := 1000
// 	for i := 0; i < len(workload); i += batchSize {
// 		j += batchSize
// 		if j > len(workload) {
// 			j = len(workload)
// 		}

// 		batch := workload[i:j]
// 		computeGPAReq := &pbWorker.ComputeGPARequest{
// 			StudentsWithGrades: batch,
// 		}

// 		if err := computeGPAStreamClient.Send(computeGPAReq); err != nil {
// 			return fmt.Errorf("failed to send a batch of students to compute GPA: %v", err)
// 		}

// 	}
// 	computeGPAStreamClient.CloseSend()

// 	select {
// 	case err := <-errChan:
// 		return err
// 	case <-waitc:
// 		break
// 	}

// 	return nil
// }
