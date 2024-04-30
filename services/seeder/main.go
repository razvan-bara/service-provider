package seeder

import (
	"log"
	"os"
	"time"
)

var (
	args = os.Args[1:]
)

func main() {
	if len(args) > 0 && args[0] == "seed" {
		log.Print("Generating CSV file...")

		now := time.Now()
		err := GenerateCSV()
		if err != nil {
			log.Fatalf("Failed to generate CSV data, got err: %v", err)
		}

		log.Printf("CSV file generated successfully, time taken: %v\n", time.Since(now))
	}

}
