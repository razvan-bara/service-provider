package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	pbWorker "service-provider/services/worker/proto"
)

type Task struct {
	Id              int
	Name            string
	Payload         []byte
	StudentsWithGPA []byte
	Status          string
	CreatedAt       string
	UpdatedAt       string
}

func (task Task) String() string {
	return fmt.Sprintf("Task: %d, Name: %s, Payload: %s, Status: %s, Created At: %s, Updated At: %s",
		task.Id, task.Name, task.Payload, task.Status, task.CreatedAt, task.UpdatedAt)
}

type TaskDB struct {
	db *sql.DB
}

func (taskDB *TaskDB) UpdateRow(task *Task, studentWithGPA []*pbWorker.StudentWithGPA) error {
	updateStatusSQL := `UPDATE queue SET status = $1, students_with_gpa = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	bytes, err := json.Marshal(studentWithGPA)
	if err != nil {
		return err
	}

	_, err = taskDB.db.Exec(updateStatusSQL, "done", bytes, task.Id)
	if err != nil {
		return err
	}
	log.Printf("Updated task with requestId: %v", task.Name)

	return nil
}

func (taskDB *TaskDB) InsertRow(requestId string, studentsWithGrades []*pbWorker.StudentWithGrades) error {
	bytes, err := json.Marshal(studentsWithGrades)
	if err != nil {
		return err
	}

	insertTaskSQL := `INSERT INTO queue (task_name, payload) VALUES ($1, $2)`
	_, err = taskDB.db.Exec(insertTaskSQL, requestId, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (taskDB *TaskDB) SelectPendingTasks(taskChan chan *Task) ([]*Task, error) {
	queryTasksSQL := `SELECT id, task_name, payload, students_with_gpa, status, created_at, updated_at FROM queue WHERE status = $1`
	rows, err := taskDB.db.Query(queryTasksSQL, "pending")
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}

		err = rows.Scan(&task.Id, &task.Name, &task.Payload, &task.StudentsWithGPA, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			rows.Close()
			log.Fatal(err)
		}
		taskChan <- task
		tasks = append(tasks, task)
	}
	close(taskChan)

	return tasks, nil
}
