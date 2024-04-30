package main

import (
	"log"
	"time"
)

func main() {

	log.Print("Generating CSV file...")

	now := time.Now()
	err := GenerateCSV()
	if err != nil {
		log.Fatalf("Failed to generate CSV data, got err: %v", err)
	}

	log.Printf("CSV file generated successfully, time taken: %v\n", time.Since(now))

}
