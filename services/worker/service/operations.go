package main

import (
	"math"
	pb "service-provider/services/worker/proto"
	"strconv"
)

const (
	numberOfCourses = 7
)

func computeGPA(student *pb.StudentWithGrades) (*pb.StudentWithGPA, error) {
	var sum int

	for _, grade := range student.Grades {
		x, err := strconv.Atoi(grade.Score)
		if err != nil {
			return nil, err
		}
		sum += x

	}

	res := float64(sum) / numberOfCourses
	return &pb.StudentWithGPA{
		StudentName: student.StudentName,
		GPA:         math.Floor(res*100) / 100,
	}, nil
}
