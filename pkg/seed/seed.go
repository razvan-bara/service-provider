package seed

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const (
	directoryName = "generated"
	gradeLimit    = 7
	minGrade      = 1
	maxGrade      = 10
)

var (
	inputFile          = fmt.Sprintf("%s/%s", directoryName, "input.csv")
	columnNames        = []string{"Name", "PAJ", "DA", "PP", "MDS", "SGSC", "IBD", "BT"}
)

func init() {

	err := os.Mkdir(directoryName, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatalln("Failed to create 'generated' directory, got err: ", err)
	}

}

func GenerateCSV() error {

	file, err := os.Create(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(columnNames)
	writer.Flush()

	for name, grades := range studentsWithGrades {
		err := writer.Write(append([]string{name}, grades...))
		if err != nil {
			log.Printf("Failed to write data for student %s, got err: %v\n", name, err)
		}

		writer.Flush()
	}

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
