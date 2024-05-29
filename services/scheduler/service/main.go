package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	pbWorker "service-provider/services/worker/proto"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	host     = "localhost"
	port     = 5435
	user     = "myuser"
	password = "mypassword"
	dbname   = "da_project"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	waitForAtLeastOneWorkerChan := make(chan struct{})
	var once sync.Once
	onceBody := func() {
		waitForAtLeastOneWorkerChan <- struct{}{}
		close(waitForAtLeastOneWorkerChan)
	}

	log.Println("Connecting to the worker servers... Please wait")

	clientConnections := map[string]*ClientConnection{}
	timeoutDuration := 15 * time.Second
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeoutDuration)
	mutex := &sync.Mutex{}
	for i := 1; i <= 3; i++ {
		go func(i int) {
			waitForWorker := make(chan struct{})
			port := fmt.Sprintf("808%d", i)

			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						log.Println("Cancel func triggered", port)
						return
					case <-time.Tick(500 * time.Millisecond):
						{
							conn, err := grpc.Dial(
								fmt.Sprintf("localhost:%s", port),
								grpc.WithInsecure(),
							)

							if err != nil {
								log.Println("Failed to connect to the worker server on port: ", port)
								time.Sleep(1 * time.Second)
								continue
							}

							pbWorkerClinet := pbWorker.NewWorkerServiceClient(conn)
							if pbWorkerClinet == nil {
								log.Println("Failed to create a new worker client")
								time.Sleep(1 * time.Second)
								continue
							}

							_, err = pbWorkerClinet.GetStatus(context.Background(), &pbWorker.GetStatusRequest{})
							if err != nil {
								mutex.Lock()
								clientConnections[port] = &ClientConnection{
									Conn:         conn,
									Client:       pbWorkerClinet,
									HasHeartbeat: false,
								}
								mutex.Unlock()
								continue
							}

							log.Println("Connected to the worker server on port: ", port)
							mutex.Lock()
							clientConnections[port] = &ClientConnection{
								Conn:         conn,
								Client:       pbWorkerClinet,
								HasHeartbeat: true,
							}
							mutex.Unlock()
							waitForWorker <- struct{}{}
							return
						}
					}
				}
			}(ctx)

			select {
			case <-waitForWorker:
				once.Do(onceBody)
			case <-time.After(timeoutDuration):

			}
			cancelFunc()
		}(i)
	}

	select {
	case <-waitForAtLeastOneWorkerChan:

		break
	case <-time.After(timeoutDuration):
		log.Fatal("Failed to connect to any worker client, exiting...")
	}

	taskDB := &TaskDB{db: db}

	scheduler := &Scheduler{
		ClientConnections: clientConnections,
		TaskDB:            taskDB,
	}
	defer scheduler.CloseConnections()

	handler := &SchedulerHandler{
		Scheduler: scheduler,
		TaskDB:    taskDB,
	}
	go scheduler.ListenForTasks()
	go scheduler.CheckWorkerHeartbeats()

	var router = gin.Default()

	router.GET("/ping", handler.getHealth)
	router.POST("/file", handler.handleCSVRequest)

	router.Run() // listen and serve on 0.0.0.0:8080

}
