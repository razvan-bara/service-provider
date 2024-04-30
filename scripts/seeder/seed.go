package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type gradesCSV struct {
	name string
	rows int64
}

const (
	directoryName = "generated"
	gradeLimit    = 7
	minGrade      = 1
	maxGrade      = 10
)

var (
	files = []gradesCSV{
		{name: fmt.Sprintf("%s/%s", directoryName, "input_100k.csv"), rows: 100_000},
		{name: fmt.Sprintf("%s/%s", directoryName, "input_500k.csv"), rows: 500_000},
		{name: fmt.Sprintf("%s/%s", directoryName, "input_1m.csv"), rows: 1_000_000},
		{name: fmt.Sprintf("%s/%s", directoryName, "input_5m.csv"), rows: 5_000_000},
		{name: fmt.Sprintf("%s/%s", directoryName, "input_10m.csv"), rows: 10_000_000},
	}
	columnNames = []string{"Name", "PAJ", "DA", "PP", "MDS", "SGSC", "IBD", "BT"}
)

func init() {

	err := os.Mkdir(directoryName, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatalln("Failed to create 'generated' directory, got err: ", err)
	}

}

func GenerateCSV() error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, file := range files {
		wg.Add(1)
		go func(name string, rows int64) {
			defer wg.Done()

			inputFile, err := os.Create(name)
			if err != nil {
				errChan <- err
				return
			}
			defer inputFile.Close()

			writer := csv.NewWriter(inputFile)
			writer.Write(columnNames)
			writer.Flush()

			studentLen := int64(len(students))
			for i := int64(0); i < rows; i++ {
				name := students[rand.Int63n(studentLen)]
				err := writer.Write(append([]string{name}, generateGrades()...))
				if err != nil {
					log.Printf("Failed to write data for student %s, got err: %v\n", name, err)
				}

				writer.Flush()
			}

		}(file.name, file.rows)

	}

	wg.Wait()
	return nil
}

func generateGrades() []string {
	grades := make([]string, gradeLimit)
	for i := 0; i < len(grades); i++ {
		grades[i] = generateRandomGrade()
	}

	return grades

}

func generateRandomGrade() string {
	n := minGrade + rand.Intn(maxGrade-minGrade+1)
	return strconv.Itoa(n)
}
