package main

import (
	"log"
	"os"
	"service-provider/pkg/seed"
	"time"
)

var (
	args = os.Args[1:]
)

func main() {
	// dockerClient, err := docker.CreateClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(dockerClient)

	if len(args) > 0 && args[0] == "seed" {
		log.Print("Generating CSV file...")

		now := time.Now()
		err := seed.GenerateCSV()
		if err != nil {
			log.Fatalf("Failed to generate CSV data, got err: %v", err)
		}

		log.Printf("CSV file generated successfully, time taken: %v\n", time.Since(now))
	}

}
