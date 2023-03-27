package database

import (
	"fmt"
	"log"
	"strings"
	"student-crud/models"
)

func Insert(student models.Student) (*models.Student, error) {
	err := DB.QueryRow("INSERT INTO students (name, roll_number) VALUES ($1, $2) RETURNING id, name, roll_number", student.Name, strings.ToUpper(student.RollNumber)).Scan(&student.Id, &student.Name, &student.RollNumber)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func UpdateByRollNumber(rollNumber string, student models.Student) (*models.Student, error) {
	err := DB.QueryRow("UPDATE students SET name=$1 WHERE roll_number=$2 RETURNING id", student.Name, strings.ToUpper(rollNumber)).Scan(student.Id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &student, nil
}

func DeleteByRollNumber(rollNumber string) (*models.Student, error) {
	student := models.Student{}
	err := DB.QueryRow("DELETE FROM students WHERE roll_number=$1 RETURNING id,name,roll_number", strings.ToUpper(rollNumber)).Scan(student.Id, student.Name, student.RollNumber)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func GetStudentByRollNumber(rollNumber string) (*models.Student, error) {
	stmt, err := DB.Prepare("SELECT id,name, roll_number FROM students WHERE roll_number=$1")
	defer stmt.Close()
	query, err := stmt.Query(strings.ToUpper(rollNumber))
	defer query.Close()
	if err != nil {
		return nil, err
	}
	if query.Next() {
		var student models.Student
		query.Scan(&student.Id, &student.Name, &student.RollNumber)
		return &student, nil
	}
	return nil, err
}

func GetAllStudents() ([]models.Student, error) {
	var students []models.Student
	query, err := DB.Query("SELECT id,name, roll_number FROM students")
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var student models.Student
		query.Scan(&student.Id, &student.Name, &student.RollNumber)
		fmt.Println(student)
		students = append(students, student)
	}
	return students, nil
}
