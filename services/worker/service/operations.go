package main

// const (
// 	numberOfCourses = 7
// )

// func computeGPAForStudents(students []types.StudentRequestRow) ([]*types.StudentResponseRow, error) {

// 	studentsOrderedByGPA := make([]*types.StudentResponseRow, len(students))

// 	for i := 0; i < len(students); i++ {
// 		gpa, err := computeGPA(students[i])
// 		if err != nil {
// 			return nil, err
// 		}
// 		studentsOrderedByGPA[i] = &types.StudentResponseRow{
// 			Name: students[i].Name,
// 			GPA:  gpa,
// 		}
// 	}

// 	sort.Slice(studentsOrderedByGPA, func(i, j int) bool {
// 		return studentsOrderedByGPA[i].GPA > studentsOrderedByGPA[j].GPA
// 	})

// 	return studentsOrderedByGPA, nil
// }

// func computeGPA(student types.StudentRequestRow) (float64, error) {
// 	var sum int
// 	x, err := convertGradeToInt(student.PAJ)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.BT)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.PP)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.DA)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.MDS)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.SGSC)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	x, err = convertGradeToInt(student.IBD)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sum += x

// 	res := float64(sum) / numberOfCourses
// 	return res, nil
// }

// func convertGradeToInt(grade string) (int, error) {
// 	x, err := strconv.Atoi(grade)
// 	if err != nil {
// 		return 0, errors.New("could not convert grade to int")
// 	}

// 	return x, nil
// }
